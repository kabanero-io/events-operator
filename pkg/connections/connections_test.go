// +build unit_test

package connections_test

import (
	"github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
	"github.com/kabanero-io/events-operator/pkg/connections"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

func TestConnectionsManager(t *testing.T) {
	cm := connections.NewConnectionsManager()

	ghUrl := "https://github-service"
	dockerUrl := "https://docker-service"
	svcUrl1 := "https://service1"
	svcUrl2 := "https://service2"

	eventConnections := v1alpha1.EventConnectionsList{
		Items: []v1alpha1.EventConnections{
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "events.kabanero.io/v1alpha1",
					Kind:       "EventConnections",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name: "event-connections-1",
				},
				Spec: v1alpha1.EventConnectionsSpec{
					Connections: []v1alpha1.EventConnection{
						{
							From: v1alpha1.EventEndpoint{
								Kind: "EventMediator",
								Name: "switchboard-mediator-1",
								Id: "webhook-switchboard/githubDest",
							},
							To: []v1alpha1.EventDestinationEndpoint{
								{
                                    Https: &[]v1alpha1.HttpsEndpoint {
                                        {
									       Url: ghUrl,
                                        },
                                    },
								},
							},
						},
						{
							From: v1alpha1.EventEndpoint{
								Kind: "EventMediator",
								Name: "switchboard-mediator-1",
								Id: "webhook-switchboard/dockerDest",
							},
							To: []v1alpha1.EventDestinationEndpoint{
								{
                                    Https: &[]v1alpha1.HttpsEndpoint {
                                        {
									       Url: dockerUrl,
                                        },
                                    },
								},
							},
						},
					},
				},
			},
			{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "events.kabanero.io/v1alpha1",
					Kind:       "EventConnections",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name: "event-connections-2",
				},
				Spec: v1alpha1.EventConnectionsSpec{
					Connections: []v1alpha1.EventConnection{
						{
							From: v1alpha1.EventEndpoint{
								Kind: "EventMediator",
								Name: "switchboard-mediator-2",
								Id: "webhook-switchboard/dest",
							},
							To: []v1alpha1.EventDestinationEndpoint{
								{
                                    Https: &[]v1alpha1.HttpsEndpoint {
                                        {
									       Url: svcUrl1,
                                        },
                                    },
								},
								{
                                    Https: &[]v1alpha1.HttpsEndpoint {
                                        {
									       Url: svcUrl2,
                                           Insecure: true,
                                        },
                                    },
								},
							},
						},
					},
				},
			},
		},
	}

	// Add connections to the connection manager
	for i := range eventConnections.Items {
		cm.AddConnections(&eventConnections.Items[i])
	}

	// Test the connections
	for _, ec := range eventConnections.Items {
		for _, conn := range ec.Spec.Connections {
			endpoints := cm.LookupDestinationEndpoints(&conn.From)
			if !reflect.DeepEqual(endpoints, conn.To) {
				t.Errorf("expected endpoints %v\nbut got: %v", conn.To, endpoints)
			}
		}
	}

	// Remove the first set of connections and check to make sure it was removed from the connection manager
	cm.RemoveConnections(&eventConnections.Items[0])

	// Make sure the first set of connections was removed
	for _, conn := range eventConnections.Items[0].Spec.Connections {
		if endpoints := cm.LookupDestinationEndpoints(&conn.From); len(endpoints) != 0 {
			t.Errorf("expected the first connection set to have been removed from the connection manager, but got: %v", endpoints)
		}
	}

	// Make sure the rest of the connections are still in the manager
	for i := 1; i < len(eventConnections.Items); i++ {
		for _, conn := range eventConnections.Items[i].Spec.Connections {
			endpoints := cm.LookupDestinationEndpoints(&conn.From)
			if !reflect.DeepEqual(endpoints, conn.To) {
				t.Errorf("expected endpoints %v\nbut got: %v", conn.To, endpoints)
			}
		}
	}
}


