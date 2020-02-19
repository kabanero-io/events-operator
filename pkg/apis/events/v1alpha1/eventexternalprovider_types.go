package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EventExternalProviderSpec defines the desired state of EventExternalProvider
type EventExternalProviderSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
    Type string `json:"type"`
    Url string `json:"url"`
    Timeout string `json:"timeout,omitempty"` 
}

// EventExternalProviderStatus defines the observed state of EventExternalProvider
type EventExternalProviderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
    Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventExternalProvider is the Schema for the eventexternalproviders API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eventexternalproviders,scope=Namespaced
type EventExternalProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventExternalProviderSpec   `json:"spec,omitempty"`
	Status EventExternalProviderStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventExternalProviderList contains a list of EventExternalProvider
type EventExternalProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventExternalProvider `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EventExternalProvider{}, &EventExternalProviderList{})
}
