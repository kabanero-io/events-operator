package eventenv
import (
    "github.com/kabanero-io/events-operator/pkg/managers"
    "github.com/kabanero-io/events-operator/pkg/connections"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "net/url"
    //"os"
)


const (
     MEDIATOR_NAME_KEY = "MEDIATOR-NAME" // environment variable. If not set, we're running as operator.
)

type  ListenerHandler func(env *EventEnv, header map[string][]string, body map[string]interface{}, key string, url *url.URL) error

type ListenerManager interface {
    /* Create a new TLS listener with TLS. Call the handler on every message received */
    NewListenerTLS(env *EventEnv, port int32, key string, tlsCertPath, tlsKeyPath string, handler ListenerHandler) error

    /* Create a new listener. */
    NewListener(env *EventEnv, port int32, key string, handler ListenerHandler ) error
}

type EventEnv  struct {
    Client client.Client
    EventMgr *managers.EventManager
    ConnectionsMgr *connections.ConnectionsManager
    ListenerMgr ListenerManager
    MediatorName string // Kubernetes name of this mediator worker if not ""
    IsOperator bool  // true if this instance is an operator, not a worker
    Namespace string // namespace we're running under
    KabaneroIntegration bool // true to integrate with Kabanero
}


var eventEnv *EventEnv

func InitEventEnv(env *EventEnv) {
    eventEnv = env
}

func GetEventEnv() *EventEnv {
    return eventEnv
}
