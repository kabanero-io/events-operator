package eventtest

import (
	"os"
	"path/filepath"

	"github.com/kabanero-io/events-operator/pkg/apis"
	"github.com/kabanero-io/events-operator/pkg/connections"
	"github.com/kabanero-io/events-operator/pkg/eventenv"
	"github.com/kabanero-io/events-operator/pkg/listeners"
	"github.com/kabanero-io/events-operator/pkg/managers"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	controller "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type AddFunc func(manager.Manager) error

type Environment struct {
	cfg             *rest.Config
	manager         controller.Manager
	testEnvironment *envtest.Environment
}

type EnvironmentOptions struct {
	AddFunc AddFunc
	MediatorName string
}

func NewEnvironment(opts EnvironmentOptions) (*Environment, error) {
	t := true
	var testEnv *envtest.Environment
	if os.Getenv("TEST_USE_EXISTING_CLUSTER") == "true" {
		testEnv = &envtest.Environment{
			UseExistingCluster: &t,
		}
	} else {
		testEnv = &envtest.Environment{
			CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "deploy", "crds")},
		}
	}

	cfg, err := testEnv.Start()
	if err != nil {
		return nil, err
	}

	err = scheme.AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, err
	}

	// +kubebuilder:scaffold:scheme

	manager, err := controller.NewManager(cfg, controller.Options{
		Scheme: scheme.Scheme,
	})
	if err != nil {
		return nil, err
	}

	// Register components
	err = apis.AddToScheme(manager.GetScheme())
	if err != nil {
		return nil, err
	}

	// set up env
	eventEnv := &eventenv.EventEnv{
		Client:         manager.GetClient(),
		EventMgr:       managers.NewEventManager(),
		ConnectionsMgr: connections.NewConnectionsManager(),
		ListenerMgr:    listeners.NewDefaultListenerManager(),
		IsOperator:     opts.MediatorName == "",
		MediatorName:   opts.MediatorName,
	}
	eventenv.InitEventEnv(eventEnv)

	if opts.AddFunc != nil {
		err = opts.AddFunc(manager)
		if err != nil {
			return nil, err
		}
	}

	env := &Environment{
		cfg:             cfg,
		manager:         manager,
		testEnvironment: testEnv,
	}

	return env, nil
}

func (env *Environment) Start() error {
	return env.manager.Start(controller.SetupSignalHandler())
}

func (env *Environment) Stop() error {
	return env.testEnvironment.Stop()
}

func (env *Environment) GetClient() client.Client {
	return env.manager.GetClient()
}
