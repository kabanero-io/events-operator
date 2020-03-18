package eventconnections_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/kabanero-io/events-operator/internal/eventtest"
	"github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
	"github.com/kabanero-io/events-operator/pkg/controller/eventconnections"
	"github.com/kabanero-io/events-operator/pkg/eventenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/onsi/gomega/gexec"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var env *eventtest.Environment

func TestEventConnections(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t,
		"EventConnectionsController Suite",
		[]Reporter{envtest.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	By("Bootstrapping test environment")

	var err error
	env, err = eventtest.NewEnvironment(eventtest.EnvironmentOptions{
		AddFunc: eventconnections.Add,
		MediatorName: "example",
	})
	Expect(err).ToNot(HaveOccurred())

	go env.Start()
	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("Tearing down the test environment")
	gexec.KillAndWait(5 * time.Second)
	err := env.Stop()
	Expect(err).ToNot(HaveOccurred())
})

var _ = Describe("EventConnectionsController", func() {
	timeout := 30 * time.Second
	interval := 1 * time.Second

	key := types.NamespacedName{
		Name:      "example",
		Namespace: "default",
	}

	spec := v1alpha1.EventConnectionsSpec{
		Connections: []v1alpha1.EventConnection{
			{
				From: v1alpha1.EventSourceEndpoint{
					Mediator: &v1alpha1.EventMediatorSourceEndpoint{
						Name:        "webhook",
						Mediation:   "webhook",
						Destination: "dest",
					},
				},
				To: []v1alpha1.EventDestinationEndpoint{
					{
						Https: &[]v1alpha1.HttpsEndpoint{
							{
								Url: "https://mediator1/mediation1",
							},
							{
								Url: "https://mediator2/mediation1",
							},
						},
					},
				},
			},
		},
	}

	// Get a resource and retry 5 times every 100 ms if resource doesn't satisfy predicate.
	getResource := func(predicate func(*v1alpha1.EventConnections) bool) bool {
		retriesLeft := 5
		for retriesLeft > 0 {
			By(fmt.Sprintf("Waiting for EventConnections with name %s in namespace %s that satifies condition", key.Name, key.Namespace))
			f := &v1alpha1.EventConnections{}
			err := env.GetClient().Get(context.Background(), key, f)
			if err == nil && predicate(f) {
				return true
			}
			retriesLeft--
			time.Sleep(100 * time.Millisecond)
		}
		return false
	}

	Context("EventConnections", func() {
		It("should be created, updated, and deleted successfully", func() {
			By("Applying a new EventConnections CR")
			created := &v1alpha1.EventConnections{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},

				Spec: spec,
			}

			numInitialConnections := eventenv.GetEventEnv().ConnectionsMgr.ConnectionCount()
			Expect(env.GetClient().Create(context.Background(), created)).Should(Succeed())

			// Wait for EventConnections to be applied
			Expect(getResource(func(f *v1alpha1.EventConnections) bool {
				return len(f.Spec.Connections) > 0
			})).To(BeTrue())

			By("Checking that it has only 1 EventConnections more than it did before")
			numConnections := eventenv.GetEventEnv().ConnectionsMgr.ConnectionCount()
			Expect(numConnections).Should(Equal(numInitialConnections + 1))

			By("Verifying that the returned endpoint from LookupDestinationEndpoints matches spec")
			for _, conn := range spec.Connections {
				endpoints := eventenv.GetEventEnv().ConnectionsMgr.LookupDestinationEndpoints(&conn.From)
				Expect(reflect.DeepEqual(endpoints, conn.To)).Should(BeTrue())
			}

			By("Updating an existing EventConnections CR")
			numInitialConnections = eventenv.GetEventEnv().ConnectionsMgr.ConnectionCount()

			// Update the spec
			spec.Connections[0].From.Mediator.Mediation = "updated-mediator"

			updated := &v1alpha1.EventConnections{}
			Expect(env.GetClient().Get(context.Background(), key, updated)).Should(Succeed())
			Expect(len(updated.Spec.Connections)).ShouldNot(BeZero())
			updated.Spec.Connections[0].From.Mediator.Mediation = "updated-mediator"
			Expect(env.GetClient().Update(context.Background(), updated)).Should(Succeed())

			// Wait for EventConnections to be applied
			Expect(getResource(func(f *v1alpha1.EventConnections) bool {
				return len(f.Spec.Connections) > 0 && f.Spec.Connections[0].From.Mediator.Mediation == "updated-mediator"
			})).To(BeTrue())

			By("Checking that the number of EventConnections is the same")
			numConnections = eventenv.GetEventEnv().ConnectionsMgr.ConnectionCount()
			Expect(numConnections).Should(Equal(numInitialConnections))

			By("Verifying that the returned endpoint from LookupDestinationEndpoints matches spec")
			for _, conn := range spec.Connections {
				endpoints := eventenv.GetEventEnv().ConnectionsMgr.LookupDestinationEndpoints(&conn.From)
				Expect(reflect.DeepEqual(endpoints, conn.To)).Should(BeTrue())
			}

			By("Deleting an EventConnections CR")
			// numInitialConnections = eventenv.GetEventEnv().ConnectionsMgr.ConnectionCount()

			// Delete the CR
			Eventually(func() error {
				f := &v1alpha1.EventConnections{}
				env.GetClient().Get(context.Background(), key, f)
				return env.GetClient().Delete(context.Background(), f)
			}, timeout, interval).Should(Succeed())

			// Trying to get CR again should fail
			Eventually(func() error {
				f := &v1alpha1.EventConnections{}
				return env.GetClient().Get(context.Background(), key, f)
			}, timeout, interval).ShouldNot(Succeed())

			// TODO: Uncomment this once connections are properly removed from the ConnectionManager.
			// By("Checking that the connection manager has 1 fewer EventConnections")
			// numConnections = eventenv.GetEventEnv().ConnectionsMgr.ConnectionCount()
			// Expect(numConnections).Should(Equal(numInitialConnections - 1))
		})
	})
})
