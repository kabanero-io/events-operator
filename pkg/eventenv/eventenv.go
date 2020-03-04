package eventenv
import (
    "github.com/kabanero-io/events-operator/pkg/managers"
    "github.com/kabanero-io/events-operator/pkg/connections"
    "sigs.k8s.io/controller-runtime/pkg/client"
)

type  ListenerHandler func(env *EventEnv, message map[string]interface{}, key string) error

type ListenerManager interface {
    /* Create a new TLS listener with TLS. Call the hndler on every emssage recevied*/
    NewListenerTLS(env *EventEnv, port int, key string, tlsCertPath, tlsKeyPath string, handler ListenerHandler) error

    /* reate a new listener. */
    NewListener(env *EventEnv, port int, key string, handler ListenerHandler ) error
}

type EventEnv  struct {
    Client client.Client
    EventMgr *managers.EventManager
    ConnectionsMgr *connections.ConnectionsManager
    ListenerMgr ListenerManager
}


var eventEnv *EventEnv

func InitEventEnv(env *EventEnv) {
    eventEnv = env
}

func GetEventEnv() *EventEnv {
    return eventEnv
}
