package connections_test

import (
	"github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
	"github.com/kabanero-io/events-operator/pkg/connections"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestConnections(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Connections Suite")
}

var _ = Describe("TestConnectionsManager", func() {

	var (
		cm					*connections.ConnectionsManager
		eventConnections	v1alpha1.EventConnectionsList
		ghUrl				string
		dockerUrl			string
		svcUrl1				string
		svcUrl2				string
		githubDest			string
		dockerDest			string
		dest				string
	)

	ghUrl = "https://github-service"
	dockerUrl = "https://docker-service"
	svcUrl1 = "https://service1"
	svcUrl2 = "https://service2"

	githubDest = "githubDest"
	dockerDest = "dockerDest"
	dest = "dest"

	eventConnections = v1alpha1.EventConnectionsList{
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
							From: v1alpha1.EventSourceEndpoint{
								Mediator: &v1alpha1.EventMediatorSourceEndpoint {
									Name: "switchboard-mediator-1",
									Mediation: "webhook-switchboard",
									Destination: githubDest,
								},
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
							From: v1alpha1.EventSourceEndpoint{
								Mediator: &v1alpha1.EventMediatorSourceEndpoint{
									Name: "switchboard-mediator-1",
									Mediation: "webhook-switchboard",
									Destination: dockerDest,
								},
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
							From: v1alpha1.EventSourceEndpoint{
								Mediator: &v1alpha1.EventMediatorSourceEndpoint {
									Name: "switchboard-mediator-2",
									Mediation : "webhook-switchboard",
									Destination: dest,
								},
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

	BeforeEach(func() {
		cm = connections.NewConnectionsManager()
	})

	Context("ConnectionManager", func() {

		It("should be tested after creation", func() {

			By("Checking if connection manager has zero connections")
			Expect(cm.ConnectionCount()).Should(BeZero())

		})

		It("should have adding connections functionality", func() {

			By("first adding the connections and", func() {
				for items := range eventConnections.Items {
					cm.AddConnections(&eventConnections.Items[items])
				}
			})

			By("testing the connections endpoints", func() {
				for _, ec := range eventConnections.Items {
					for _, conn := range ec.Spec.Connections {
						endpoints := cm.LookupDestinationEndpoints(&conn.From)
						Expect(endpoints).Should(Equal(conn.To))
					}
				}
			})

		})

		It("should have remove connections functionality", func() {

			for items := range eventConnections.Items {
				cm.AddConnections(&eventConnections.Items[items])
			}

			Expect(cm.ConnectionCount()).ShouldNot(BeZero())

			By("trying to remove a set of connections and", func() {
				cm.RemoveConnections(&eventConnections.Items[0])
			})

			By("making sure the first set of connections was removed", func() {
				for _, conn := range eventConnections.Items[0].Spec.Connections {
					endpoints := cm.LookupDestinationEndpoints(&conn.From)
					Expect(len(endpoints)).Should(BeZero())
				}
			})

			By("making sure the second set of connections is still there", func() {
				for _, conn := range eventConnections.Items[1].Spec.Connections {
					endpoints := cm.LookupDestinationEndpoints(&conn.From)
					Expect(len(endpoints)).Should(Not(BeZero()))
					Expect(endpoints).Should(Equal(conn.To))
				}
			})

		})
	})
})


