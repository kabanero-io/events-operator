package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EventMediatorSpec defines the desired state of EventMediator
type EventMediatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
    Listeners *[]EventListenerConfig `json:"listeners",omitempty` // default is no listener
    ImportMediations  *[]string `json:"importMediations",omitempty` // default is to import everything unless code is specified
    Mediations *[]MediationsImpl `json:"mediations"`
}


type EventListenerConfig struct {
    Name string  `json:"name",omitempty`  // name of the listener configuration. Default is the name of the MediatorSpec
    Mediations []string `json:"mediations"` // if not specified, applies to all mediations
    HttpPort    int         `json:"httpPort",omitempty`
    HttpsPort   int         `json:"httpsPort",omitempty`
    CreateService bool `json:"createService",omitempty`
    CreateRoute bool      `json:"createRoute",omitempty`
}

// EventMediatorStatus defines the observed state of EventMediator
type EventMediatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
    Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventMediator is the Schema for the eventmediators API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eventmediators,scope=Namespaced
type EventMediator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventMediatorSpec   `json:"spec,omitempty"`
	Status EventMediatorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventMediatorList contains a list of EventMediator
type EventMediatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventMediator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EventMediator{}, &EventMediatorList{})
}
