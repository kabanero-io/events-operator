package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EventConnectionsSpec defines the desired state of EventConnections
type EventConnectionsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
    Connections []EventConnection `json:"connections"`
}

/* Connections are from subscriber to publishers
   from sender to receivers
*/
type EventConnection struct {
    From EventEndpoint `json:"from"`
    To  [] EventEndpoint  `json:"to"`
}

type EventEndpoint struct {
    Group string `json:"group,omitempty"`
    Kind  string `json:"kind"`
    Namespace string `json:"namespace,omitempty"`
    Name  string `json:"name"`
    Id    string `json:"id,omitempty"` // Identifier of the endpoint
}


// EventConnectionsStatus defines the observed state of EventConnections
type EventConnectionsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
    Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventConnections is the Schema for the eventconnections API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eventconnections,scope=Namespaced
type EventConnections struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventConnectionsSpec   `json:"spec,omitempty"`
	Status EventConnectionsStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventConnectionsList contains a list of EventConnections
type EventConnectionsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventConnections `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EventConnections{}, &EventConnectionsList{})
}
