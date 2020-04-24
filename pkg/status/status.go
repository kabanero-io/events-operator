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
    "container/list"
    "time"
    "sync"
    "k8s.io/klog"
//    "fmt"
)

const (
    MAX_RETAINED_MESSAGES = 100 // maximum number of messages to retain

   /* Operations names */
   OPERATION_VALIDATE_WEBHOOK_SECRET = "validate webhook secret"
   OPERATION_RESOLVE_REPOSITORY_TYPE = "resolve repository type"
   OPERATION_FIND_MEDIATION = "find mediation"
   OPERATION_INITIALIZE_VARIABLES = "initialize mediation variables"
   OPERATION_EVALUATE_MEDIATION = "evaluate mediation"
   OPERATION_SEND_EVENT = "send event"

   /* Parameter names */
   PARAM_FROM = "from"
   PARAM_MEDIATION = "mediation"
   PARAM_FILE = "file"
   PARAM_DESTINATION = "destination"
   PARAM_URL = "url"
   PARAM_URLEXPRESSION = "urlExpression"

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
