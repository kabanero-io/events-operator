package status

import (
    eventsv1alpha1 "github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
    "testing"
    "strconv"
    "fmt"
)


func compareList(list1 []eventsv1alpha1.EventStatusSummary, list2[]eventsv1alpha1.EventStatusSummary) bool {
    if len(list1) != len(list1) {
        return false
    }

    for index, elem := range list1{
        if !elem.Equals(&list2[index]) {
            return false
        }
    }
    return true
}

func printList(list []eventsv1alpha1.EventStatusSummary) {
   fmt.Printf("List: ")
   for _, entry := range list {
      fmt.Printf("%v ", entry.Operation)
   }
   fmt.Printf("\n")
}

func TestAppendStatus(t *testing.T) {
    statusMgr := NewStatusManager()

    arraySummary := make([]eventsv1alpha1.EventStatusSummary , 0)
    for i := 0; i < MAX_RETAINED_MESSAGES; i++ {
         summary := &eventsv1alpha1.EventStatusSummary {
             Operation: strconv.Itoa(i),
             Input: make([]eventsv1alpha1.EventStatusParameter,0),
             Result: "",
             Message:"",
         }
         arraySummary = append(arraySummary, *summary)
         statusMgr.AddEventSummary(summary)
    }

    resultSummary := statusMgr.getStatusSummary()
    resultLen := len(resultSummary)
    if resultLen != MAX_RETAINED_MESSAGES {
         t.Fatalf("Cached Status length %v does not matched expected length %v", resultLen, MAX_RETAINED_MESSAGES)
    }

    if !compareList(arraySummary, resultSummary) {
         t.Fatalf("Cached Status did not match expected")
    }
}

func TestDuplicateStatus(t *testing.T) {
    statusMgr := NewStatusManager()

    arraySummary := make([]eventsv1alpha1.EventStatusSummary , 0)

    /* append firs time */
    for i := 0; i < MAX_RETAINED_MESSAGES; i++ {
         summary := &eventsv1alpha1.EventStatusSummary {
             Operation: strconv.Itoa(i),
             Input: make([]eventsv1alpha1.EventStatusParameter,0),
             Result: "",
             Message:"",
         }
         arraySummary = append(arraySummary, *summary)
         statusMgr.AddEventSummary(summary)
    }

    /* append second time */
    for i := MAX_RETAINED_MESSAGES-1; i>=0 ; i-- {
         statusMgr.AddEventSummary(&arraySummary[i])
    }

    resultSummary := statusMgr.getStatusSummary()
    resultLen := len(resultSummary)
    if resultLen != MAX_RETAINED_MESSAGES {
         t.Fatalf("Cached Status length %v does not matched expected length %v", resultLen, MAX_RETAINED_MESSAGES)
    }

    if !compareList(arraySummary, resultSummary) {
         t.Fatalf("Cached Status did not match expected")
    }
}

func TestOverflow(t *testing.T) {
    statusMgr := NewStatusManager()

    arraySummary := make([]eventsv1alpha1.EventStatusSummary , 0)

    /* append firs time */
    for i := 0; i < MAX_RETAINED_MESSAGES; i++ {
         summary := &eventsv1alpha1.EventStatusSummary {
             Operation: strconv.Itoa(i),
             Input: make([]eventsv1alpha1.EventStatusParameter,0),
             Result: "",
             Message:"",
         }
         arraySummary = append(arraySummary, *summary)
         statusMgr.AddEventSummary(summary)
    }

    /* append second time, which should wipe out those added previously */
    arraySummary = make([]eventsv1alpha1.EventStatusSummary , 0)
    for i := MAX_RETAINED_MESSAGES ; i < MAX_RETAINED_MESSAGES*2; i++ {
         summary := &eventsv1alpha1.EventStatusSummary {
             Operation: strconv.Itoa(i),
             Input: make([]eventsv1alpha1.EventStatusParameter,0),
             Result: "",
             Message:"",
         }
         arraySummary = append(arraySummary, *summary)
         statusMgr.AddEventSummary(summary)
    }

    resultSummary := statusMgr.getStatusSummary()
    resultLen := len(resultSummary)
    if resultLen != MAX_RETAINED_MESSAGES {
         t.Fatalf("Cached Status length %v does not matched expected length %v", resultLen, MAX_RETAINED_MESSAGES)
    }

    if !compareList(arraySummary, resultSummary) {
         t.Fatalf("Cached Status did not match expected")
    }
}

func TestOverflowDuplicate(t *testing.T) {
    statusMgr := NewStatusManager()

    arraySummary := make([]eventsv1alpha1.EventStatusSummary , 0)

    /* append first time */
    for i := 0; i < MAX_RETAINED_MESSAGES; i++ {
         summary := &eventsv1alpha1.EventStatusSummary {
             Operation: strconv.Itoa(i),
             Input: make([]eventsv1alpha1.EventStatusParameter,0),
             Result: "",
             Message:"",
         }
         arraySummary = append(arraySummary, *summary)
         statusMgr.AddEventSummary(summary)
    }

    /* Update odd entries. This will refresh their timestamp */
    for i := 0; i < MAX_RETAINED_MESSAGES; i += 2 {
         statusMgr.AddEventSummary(&arraySummary[i])
    }

    /* Update 1/2 of the entries. This will update the even entries. Odd entries will be replaced earlier. */
    for i := 0; i < MAX_RETAINED_MESSAGES/2; i++ {
         /* The expected first half will be the even entries */
         summary := &eventsv1alpha1.EventStatusSummary {
             Operation: strconv.Itoa(i*2),
             Input: make([]eventsv1alpha1.EventStatusParameter,0),
             Result: "",
             Message:"",
         }
         arraySummary[i] = *summary // update expected result


         /* The expected 2nd half will be the new entries being added */
         num := i + MAX_RETAINED_MESSAGES
         summary = &eventsv1alpha1.EventStatusSummary {
             Operation: strconv.Itoa(num),
             Input: make([]eventsv1alpha1.EventStatusParameter,0),
             Result: "",
             Message:"",
         }
         arraySummary[i+MAX_RETAINED_MESSAGES/2] = *summary // update expected result
         statusMgr.AddEventSummary(summary)
    }

    resultSummary := statusMgr.getStatusSummary()
    resultLen := len(resultSummary)
    if resultLen != MAX_RETAINED_MESSAGES {
         t.Fatalf("Cached Status length %v does not matched expected length %v", resultLen, MAX_RETAINED_MESSAGES)
    }

    if !compareList(arraySummary, resultSummary) {
         printList(arraySummary)
         printList(resultSummary)
         t.Fatalf("Cached Status did not match expected")
    }
}
