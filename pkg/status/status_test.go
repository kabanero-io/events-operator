package status_test

import (
    eventsv1alpha1 "github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
    "github.com/kabanero-io/events-operator/pkg/status"
    "github.com/kabanero-io/events-operator/pkg/status/helpers"
    "strconv"
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestEvent(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Event Suite")
}

var _ = Describe("StatusTest", func() {

    var (
        statusMgr       *status.StatusManager
        arraySummary    []eventsv1alpha1.EventStatusSummary
    )

    BeforeEach(func() {
        statusMgr = status.NewStatusManager()

        arraySummary = make([]eventsv1alpha1.EventStatusSummary, 0)

        for i := 0; i < status.MAX_RETAINED_MESSAGES; i++ {
            summary := &eventsv1alpha1.EventStatusSummary {
                Operation: strconv.Itoa(i),
                Input: make([]eventsv1alpha1.EventStatusParameter,0),
                Result: "",
                Message:"",
            }
            arraySummary = append(arraySummary, *summary)
            statusMgr.AddEventSummary(summary)
        }
    })

    It("should append correctly the expected amount of status summary events to the status manager", func() {
        resultSummary := statusMgr.GetStatusSummary()
        resultLen := len(resultSummary)
        Expect(resultLen).Should(Equal(status.MAX_RETAINED_MESSAGES))
        Expect(helpers.CompareList(arraySummary,resultSummary)).Should(BeTrue())
    })

    It("should ignore duplicates and the cache status length should match the constant expected length", func() {
        /* append second time, in reverse order */
        arraySummary1 := make([]eventsv1alpha1.EventStatusSummary , 0)
        for i := status.MAX_RETAINED_MESSAGES-1; i>=0 ; i-- {
            summary := arraySummary[i]
            statusMgr.AddEventSummary(&summary)
            arraySummary1 = append(arraySummary1, arraySummary[i])
        }
        resultSummary := statusMgr.GetStatusSummary()
        resultLen := len(resultSummary)
        Expect(resultLen).Should(Equal(status.MAX_RETAINED_MESSAGES))
        Expect(helpers.CompareList(arraySummary1,resultSummary)).Should(BeTrue())
    })

    It("should wipe out existing data when adding new one when the cache is at full capacity", func() {
        /* append second time, which should wipe out those added previously */
        arraySummary = make([]eventsv1alpha1.EventStatusSummary , 0)
        for i := status.MAX_RETAINED_MESSAGES ; i < status.MAX_RETAINED_MESSAGES*2; i++ {
            summary := &eventsv1alpha1.EventStatusSummary {
                Operation: strconv.Itoa(i),
                Input: make([]eventsv1alpha1.EventStatusParameter,0),
                Result: "",
                Message:"",
            }
            arraySummary = append(arraySummary, *summary)
            statusMgr.AddEventSummary(summary)
        }
        resultSummary := statusMgr.GetStatusSummary()
        resultLen := len(resultSummary)
        Expect(resultLen).Should(Equal(status.MAX_RETAINED_MESSAGES))
        Expect(helpers.CompareList(arraySummary,resultSummary)).Should(BeTrue())
    })

    It("should wipe out existing data and ingnore duplicates when adding new one when the cache is at full capacity", func() {
        /* Update even entries. This will refresh their timestamp, and push them down the list. */
        for i := 0; i < status.MAX_RETAINED_MESSAGES; i += 2 {
            summary := arraySummary[i]
            statusMgr.AddEventSummary(&summary)
        }

        /* Update 1/2 of the entries. The top bottom half will be the new entries, and the top half the old even entries */
        for i := 0; i < status.MAX_RETAINED_MESSAGES/2; i++ {
            /* The expected first half will be the new entries */
            summary := &eventsv1alpha1.EventStatusSummary {
                Operation: strconv.Itoa(i+ status.MAX_RETAINED_MESSAGES),
                Input: make([]eventsv1alpha1.EventStatusParameter,0),
                Result: "",
                Message:"",
            }
            arraySummary[i+status.MAX_RETAINED_MESSAGES/2] = *summary
            statusMgr.AddEventSummary(summary)

            /* The expected top half will be the prevous even entries */
            summary = &eventsv1alpha1.EventStatusSummary {
                Operation: strconv.Itoa(i*2),
                Input: make([]eventsv1alpha1.EventStatusParameter,0),
                Result: "",
                Message:"",
            }
            arraySummary[i] = *summary // update expected result
        }
        resultSummary := statusMgr.GetStatusSummary()
        resultLen := len(resultSummary)
        Expect(resultLen).Should(Equal(status.MAX_RETAINED_MESSAGES))
        Expect(helpers.CompareList(arraySummary,resultSummary)).Should(BeTrue())
    })
})

