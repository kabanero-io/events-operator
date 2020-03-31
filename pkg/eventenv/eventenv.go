package eventenv

import (
	"net/url"

	"github.com/kabanero-io/events-operator/pkg/connections"
	"github.com/kabanero-io/events-operator/pkg/managers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	MEDIATOR_NAME_KEY = "MEDIATOR-NAME" // environment variable. If not set, we're running as operator.
)

type ListenerOptions struct {
	Port        int32
	TLSCertPath string
	TLSKeyPath  string
}

type ListenerHandler func(env *EventEnv, message map[string]interface{}, key string, url *url.URL) error

type ListenerManager interface {
	/* Create a new TLS listener with TLS. Call the handler on every message received */
	NewListenerTLS(env *EventEnv, key string, handler ListenerHandler, options ListenerOptions) error

	/* Create a new listener. */
	NewListener(env *EventEnv, key string, handler ListenerHandler, options ListenerOptions) error
}

type EventEnv struct {
	Client         client.Client
	EventMgr       *managers.EventManager
	ConnectionsMgr *connections.ConnectionsManager
	ListenerMgr    ListenerManager
	MediatorName   string // Kubernetes name of this mediator worker if not ""
	IsOperator     bool   // true if this instance is an operator, not a worker
}

var eventEnv *EventEnv

func InitEventEnv(env *EventEnv) {
	eventEnv = env
}

func GetEventEnv() *EventEnv {
	return eventEnv
}
