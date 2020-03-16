package connections

import (
    eventsv1alpha1 "github.com/kabanero-io/events-operator/pkg/apis/events/v1alpha1"
    "k8s.io/klog"
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
func (connectionsMgr *ConnectionsManager) LookupDestinationEndpoints(endpoint *eventsv1alpha1.EventSourceEndpoint) []eventsv1alpha1.EventDestinationEndpoint {
    connectionsMgr.mutex.Lock()
    defer connectionsMgr.mutex.Unlock()

    if endpoint.Mediator != nil {
        klog.Infof("LookupDestnationEndpoins for name: %v, mediation: %v, destination: %v", endpoint.Mediator.Name, endpoint.Mediator.Mediation, endpoint.Mediator.Destination)
    }
    ret := make([]eventsv1alpha1.EventDestinationEndpoint, 0)
    /* iterate through each registered connections */
    for _, conn := range connectionsMgr.connections {
        /* iterate through eacn EventConnection */
        for _, eventConn := range conn.Spec.Connections {
            matched := eventEndpointMatch(endpoint, &eventConn.From)
            if endpoint.Mediator != nil && eventConn.From.Mediator != nil  {
                klog.Infof("eventEndpointMatch:acutal: name: %v, mediation: %v, destination: %v, connections: name: %v, mediations: %v, destination: %v, equals: %v", endpoint.Mediator.Name, endpoint.Mediator.Mediation, endpoint.Mediator.Destination, eventConn.From.Mediator.Name, eventConn.From.Mediator.Mediation, eventConn.From.Mediator.Destination, matched)
             }
            if matched {
                 /* TODO: duplicate elimination */
                for _, to := range eventConn.To {
                    ret = append(ret, to)
                }
            }
       }
    }
    klog.Infof("LookupDestnationEndpoins returned %v endpoints", len(ret))
    return ret
}

/* Return true if two event endpoints match. 
   - actual: actual resource in the Eventmediator
     toMatch: resource defined in EventConnections 
*/
func eventEndpointMatch(actual *eventsv1alpha1.EventSourceEndpoint, resource *eventsv1alpha1.EventSourceEndpoint) bool{
    if actual.Mediator == resource.Mediator {
        return true
    }

    if actual.Mediator == nil {
        return false
    }

    if resource.Mediator == nil {
        return false
    }

    if actual.Mediator.Name != resource.Mediator.Name {
        return false
    }

    if actual.Mediator.Mediation != resource.Mediator.Mediation {
        return false
    }

    if actual.Mediator.Destination != resource.Mediator.Destination {
        return false
    }

    return true
}
