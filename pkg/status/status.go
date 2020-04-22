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
)

const (
    MAX_RETAINED_MESSAGES = 100 // maximum number of messages to retain
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

func (sm *StatusManager) AddEventSummary(summary *eventsv1alpha1.EventStatusSummary) {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()

    var elem *list.Element
    for elem = sm.summaryList.Front(); elem != nil; elem = elem.Next() {
         summaryElem, ok := elem.Value.(*summaryElement)
         if !ok {
              // should not happen
              return
         }
         if summaryElem.summary.Equals(summary) {
             /* no change. Just update timestampe  */
             summaryElem.timestamp = time.Now()
             return
         }
    }

    /* summary did not repeat. */
    if sm.summaryList.Len() >= MAX_RETAINED_MESSAGES {
        /* Delete earliest one */
        earliestTime := time.Now()
        var toDelete *list.Element = nil
        for elem = sm.summaryList.Front(); elem != nil; elem = elem.Next() {
            summaryElem, ok := elem.Value.(*summaryElement)
            if !ok {
                // should not happen
                return
            }
            if summaryElem.timestamp.Before(earliestTime) || summaryElem.timestamp.Equal(earliestTime) {
                 earliestTime = summaryElem.timestamp
                 toDelete = elem
            }
        }
        if toDelete != nil {
            sm.summaryList.Remove(toDelete)
        }
    }

    /* Add a new one */
    newElem := &summaryElement {
        summary : summary,
        timestamp: time.Now(),
    }
    sm.summaryList.PushBack(newElem)
    sm.needsUpdate = true
}

func (sm *StatusManager) getStatusSummary() []eventsv1alpha1.EventStatusSummary {
    sm.mutex.Lock()
    defer sm.mutex.Unlock()


    ret := make([]eventsv1alpha1.EventStatusSummary,0)
    var elem *list.Element
    for elem = sm.summaryList.Front(); elem != nil; elem = elem.Next() {
         summaryElem, ok := elem.Value.(*summaryElement)
         if !ok {
              break
         }
         ret = append(ret, *summaryElem.summary)
    }
    return ret
}
