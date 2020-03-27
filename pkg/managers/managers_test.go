package managers

import (
	"testing"

	"github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEventManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Event Manager Suite")
}

var _ = Describe("TestEvbentManager", func() {
	var mgr *EventManager

	BeforeEach(func() {
		mgr = NewEventManager()
	})

	assignStatement := "sendEvent(dest, body, header)"
	mediator := &v1alpha1.EventMediator{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "events.kabanero.io/v1alpha1",
			Kind:       "EventMediator",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "event-mediator-1",
		},
		Spec: v1alpha1.EventMediatorSpec{
			ListenerPort: 9443,
			CreateRoute:  true,
			Mediations: &[]v1alpha1.EventMediationImpl{
				{
					Name:   "mediation-test",
					SendTo: []string{"dest"},
					Body: []v1alpha1.EventStatement{
						{
							Assign: &assignStatement,
						},
					},
				},
			},
		},
	}

	Context("EventManager", func() {
		It("should add an EventMediator successfully", func() {
			numInitialManagers := len(mgr.GetMediatorManagers())
			key := v1alpha1.MediatorHashKey(mediator)
			mgr.AddEventMediator(mediator)
			Expect(mgr.GetMediator(key)).ToNot(BeNil())
			Expect(len(mgr.GetMediatorManagers())).Should(Equal(numInitialManagers + 1))
		})
	})
})
