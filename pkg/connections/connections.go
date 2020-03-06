package connections

import (
    eventsv1alpha1 "github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
    "sync"
)

type ConnectionsManager struct {
    connections map[string]*eventsv1alpha1.EventConnections // list of connections
    mutex sync.Mutex
}

func NewConnectionsManager() *ConnectionsManager {
    return &ConnectionsManager {
         connections: make(map[string]*eventsv1alpha1.EventConnections),
    }
}

func getKey( connections *eventsv1alpha1.EventConnections) string {
    return connections.APIVersion + "/" + connections.Kind + "/" + connections.Namespace + "/" + connections.Name
}

func (connectionsMgr *ConnectionsManager) AddConnections(connections *eventsv1alpha1.EventConnections) {
    connectionsMgr.mutex.Lock()
    defer connectionsMgr.mutex.Unlock()
    key := getKey(connections)
    connectionsMgr.connections[ key] = connections

}

func (connectionsMgr *ConnectionsManager) RemoveConnections(connections *eventsv1alpha1.EventConnections) {
    connectionsMgr.mutex.Lock()
    defer connectionsMgr.mutex.Unlock()

    key := getKey(connections)
    if _, found := connectionsMgr.connections[key]; found {
        delete(connectionsMgr.connections, key)
    }
}

/* Lookup destination endpoints for an actual endpoint */
func (connectionsMgr *ConnectionsManager) LookupDestinationEndpoints(endpoint *eventsv1alpha1.EventEndpoint) []eventsv1alpha1.EventDestinationEndpoint {
    connectionsMgr.mutex.Lock()
    defer connectionsMgr.mutex.Unlock()

    ret := make([]eventsv1alpha1.EventDestinationEndpoint, 0)
    /* iterate through each registered connections */
    for _, conn := range connectionsMgr.connections {
        /* iterate through eacn EventConnection */
        for _, eventConn := range conn.Spec.Connections {
            if eventEndpointMatch(endpoint, &eventConn.From) {
                 /* TODO: duplicate elimination */
                for _, to := range eventConn.To {
                    ret = append(ret, to)
                }
            }
       }
    }
    return ret
}

/* Return true if two event endpoints match. 
   - actual: actual resource in the Eventmediator
     toMatch: resource defined in EventConnections 
*/
func eventEndpointMatch(actual *eventsv1alpha1.EventEndpoint, resource *eventsv1alpha1.EventEndpoint) bool{
    if actual.Group != resource.Group {
         if resource.Group  != "" {
              return false
         }
    }
    if actual.Kind != resource.Kind {
         return false
    }


    if actual.Name != resource.Name {
          return false
    }

    if actual.Id != resource.Id {
        return false
    }
    return true
}
