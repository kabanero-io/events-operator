/*
Copyright 2020 IBM Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package status

import (
    eventsv1alpha1 "github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
    "github.com/kabanero-io/events-operator/pkg/utils"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "container/list"
    "time"
    "sync"
    "k8s.io/klog"
//    "fmt"
)

const (
    MAX_RETAINED_MESSAGES = 100 // maximum number of messages to retain

   /* Operations names */
   OPERATION_VALIDATE_WEBHOOK_SECRET = "validate-webhook-secret"
   OPERATION_RESOLVE_REPOSITORY_TYPE = "resolve-repository-type"
   OPERATION_FIND_MEDIATION = "find-mediation"
   OPERATION_INITIALIZE_VARIABLES = "initialize-mediation-variables"
   OPERATION_EVALUATE_MEDIATION = "evaluate-mediation"
   OPERATION_SEND_EVENT = "send-event"

   /* Parameter names */
   PARAM_FROM = "from"
   PARAM_MEDIATION = "mediation"
   PARAM_FILE = "file"
   PARAM_DESTINATION = "destination"
   PARAM_URL = "url"
   PARAM_URLEXPRESSION = "urlExpression"
   PARAM_REPOSITORY = "repository"
   PARAM_BRANCH = "branch"
   PARAM_GITHUB_EVENT = "github-event"
   PARAM_STACK = "stack"

   /* Results */
   RESULT_FAILED = "failed"
   RESULT_COMPLETED = "completed"

)



type  StatusManager struct {
    summaryList *list.List 
    needsUpdate bool
    mutex sync.Mutex
}

type summaryElement struct {
   summary *eventsv1alpha1.EventStatusSummary
   timestamp time.Time
}

func NewStatusManager() *StatusManager {
    sm := &StatusManager {
        summaryList : list.New(),
        needsUpdate: true,
    }
    return sm
}

/* Add an EventSummary. 
input:
   summary: the summary to add. The caller no longer owns the summary after calling.
*/
func (sm *StatusManager) AddEventSummary(summary *eventsv1alpha1.EventStatusSummary) {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    klog.Infof("AddEventSummary: %v", *summary)


    // fmt.Printf("AddEventSummary: inserting %v\n", summary.Operation)
    /* set the time of update */
    summary.Time.Time = time.Now()

    // fmt.Printf("AddEventSummary: checking: ")
    for elem := sm.summaryList.Front(); elem != nil; elem = elem.Next() {
         summaryElem, ok := elem.Value.(*eventsv1alpha1.EventStatusSummary)
         if !ok {
              // should not happen
              return
         }
         // fmt.Printf("%v ", summaryElem.Operation)
         if summaryElem.Equals(summary) {
             /* Duplicate. Move the element to the back */
             // fmt.Printf("AddEventSummary: duplicate found %v %v\n", summaryElem.Operation, summary.Operation)
             sm.summaryList.Remove(elem)
             sm.summaryList.PushBack(summary)
             sm.needsUpdate = true
             return
         }
    }

    /* summary did not repeat. */
    if sm.summaryList.Len() >= MAX_RETAINED_MESSAGES {
        /* Delete earliest one */
        front := sm.summaryList.Front()
        if front != nil {
            sm.summaryList.Remove(front)
        }
    }

    /* Append the new one */
    sm.summaryList.PushBack(summary)
    sm.needsUpdate = true
}

func (sm *StatusManager) getStatusSummary() []eventsv1alpha1.EventStatusSummary {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    return sm.getStatusSummaryHelper()
}

/*
Unsynchronized version to fetch status summary
*/
func (sm *StatusManager) getStatusSummaryHelper() []eventsv1alpha1.EventStatusSummary {

    ret := make([]eventsv1alpha1.EventStatusSummary,0)
    var elem *list.Element
    for elem = sm.summaryList.Front(); elem != nil; elem = elem.Next() {
         summaryElem, ok := elem.Value.(*eventsv1alpha1.EventStatusSummary)
         if !ok {
              break
         }
         ret = append(ret, *summaryElem)
    }
    return ret
}

/* Send status to Status Updater if there is a status change
*/
func (sm *StatusManager) SendStatus(updater *Updater) {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    if sm.needsUpdate {
        summaryArray := sm.getStatusSummaryHelper()
        updater.SendUpdate(summaryArray)
        sm.needsUpdate = false
    }
}

/* The Status Updater is used to update status in a resource efficient manner.
   Status changes are not updated immediately. Instead, they are accumulated for up to a configurable time.
   In addition, there is no polling thread to check for status updates.
  */
type Updater struct {
    duration time.Duration // how long to wait for changes before updating to Kubernetes
    client client.Client // controller client
    timerStarted bool // true if timer started
    timerChan chan struct{} // channel for timer
    statusChan chan *[]eventsv1alpha1.EventStatusSummary // channel to send status update

    summary *[]eventsv1alpha1.EventStatusSummary  // latest summary

    mutex sync.Mutex
}

func NewSatusUpdater(client client.Client, namespace string, name string, duration time.Duration) *Updater {
    updater := &Updater {
        client: client,
        duration: duration,
        timerStarted: false,
        timerChan: make(chan struct{}, 1),
        statusChan: make(chan *[]eventsv1alpha1.EventStatusSummary),
        summary: nil,
    }

    // Thread to update status
    go func() {
        for  {
             status := updater.getStatus()
             err := utils.UpdateStatus(updater.client, namespace, name, *status)
             if err != nil {
                updater.putBack(status)
             }
        }
    }()
    return updater
}

/* Put back what was processed. */
func (updater *Updater) putBack(summary *[]eventsv1alpha1.EventStatusSummary) {
    updater.mutex.Lock()
    defer updater.mutex.Unlock()

    /* Put back only if there is nothing newer */
    if updater.summary == nil {
        updater.summary = summary
        updater.startTimer()
    }
}

/* Get available status. Block if needed. */
func (updater *Updater) getStatus() *[]eventsv1alpha1.EventStatusSummary {
    for {
         select {
              case summary, _:= <- updater.statusChan:
                  klog.Infof("Updater getStatus: Received status")
                  updater.mutex.Lock()
                  updater.summary = summary
                  updater.startTimer()
                  updater.mutex.Unlock()
              case <- updater.timerChan:
                  updater.mutex.Lock()
                  updater.timerStarted = false
                  klog.Infof("Updater getStatus: Timer fired, has status: %v", updater.summary != nil)
                  ret := updater.summary
                  updater.summary = nil
                  updater.mutex.Unlock()
                  if ret != nil {
                      return ret
                  }
         }
    }

}

/* Send Update */
func (updater *Updater) SendUpdate(summary []eventsv1alpha1.EventStatusSummary) {
    klog.Infof("Updater SendUpdate called")
    updater.statusChan <- &summary
}

func (updater *Updater) startTimer() {
    if !updater.timerStarted {
        updater.timerStarted = true
        timerChan := updater.timerChan
        duration := updater.duration
        go func() {
            time.Sleep(duration)
            timerChan <- struct{}{}
        }()
    }
}


type StatusParameters struct {
   params [] eventsv1alpha1.EventStatusParameter
}

func NewStatusParameters() *StatusParameters {
   return  &StatusParameters {
       params: make([]eventsv1alpha1.EventStatusParameter, 0),
   }
}

func (sp *StatusParameters) AddParameter(name string, value string) {
    sp.params = append(sp.params, eventsv1alpha1.EventStatusParameter { Name: name, Value: value, })
}

func (sp *StatusParameters) GetStatusParameters() []eventsv1alpha1.EventStatusParameter {
    return sp.params
}
