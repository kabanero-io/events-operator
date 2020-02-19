package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EventExternalListenerSpec defines the desired state of EventExternalListener
type EventExternalListenerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
    Url string `json:"url"`
}

// EventExternalListenerStatus defines the observed state of EventExternalListener
type EventExternalListenerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventExternalListener is the Schema for the eventexternallisteners API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eventexternallisteners,scope=Namespaced
type EventExternalListener struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventExternalListenerSpec   `json:"spec,omitempty"`
	Status EventExternalListenerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventExternalListenerList contains a list of EventExternalListener
type EventExternalListenerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventExternalListener `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EventExternalListener{}, &EventExternalListenerList{})
}
