// +build unit_test

package managers

import (
	"github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestEventManager(t *testing.T) {
	mgr := NewEventManager()
	assignStatement := "sendEvent(dest1, message.body, message.header)"
	filterStatement := `outHeader = filter(inHeader, "key.startsWith(\"X-Github\") || key == \"X-Hub-Signature\"")`

	mediator := &v1alpha1.EventMediator{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "events.kabanero.io/v1alpha1",
			Kind: "EventMediator",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name: "event-mediator-1",
		},
		Spec: v1alpha1.EventMediatorSpec{
			ListenerPort: 9443,
			CreateRoute: true,
			Mediations: &[]v1alpha1.MediationsImpl{
				{
					Mediation: &v1alpha1.EventMediationImpl{
						Name:   "mediation-test",
						Input:  "message",
						SendTo: []string{"dest1", "dest2"},
						Body: []v1alpha1.EventStatement{
							{
								Assign: &assignStatement,
							},
						},
					},
				},
				{
					Function: &v1alpha1.EventFunctionImpl{
						Name: "filterGitHubHeader",
						Input: "inHeader",
						Output: "outHeader",
						Body: []v1alpha1.EventStatement{
							{
								Assign: &filterStatement,
							},
						},
					},
				},
			},
		},
	}
	key := v1alpha1.MediatorHashKey(mediator)
	mgr.AddEventMediator(mediator)

	// Try to retrieve the event mediator
	if em := mgr.GetMediator(key); em == nil {
		t.Fatalf("expected to find an event mediator with key '%s'", key)
	}

	// Verify that there is only one mediator manager
	if mgrs := mgr.GetMediatorManagers(); len(mgrs) != 1 {
		t.Fatalf("expected to only find 1 mediator manager, but found %v: %v", len(mgrs), mgrs)
	}
}
