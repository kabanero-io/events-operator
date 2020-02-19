package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EventExternalTopicSpec defines the desired state of EventExternalTopic
type EventExternalTopicSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
    ProviderName string `json:"providerName"` // name of provider
    Topic *string `json:"topic,omitempty"` // if null, use the Kubernetes resource name as topic
}

// EventExternalTopicStatus defines the observed state of EventExternalTopic
type EventExternalTopicStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
    Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventExternalTopic is the Schema for the eventexternaltopics API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eventexternaltopics,scope=Namespaced
type EventExternalTopic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventExternalTopicSpec   `json:"spec,omitempty"`
	Status EventExternalTopicStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventExternalTopicList contains a list of EventExternalTopic
type EventExternalTopicList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventExternalTopic `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EventExternalTopic{}, &EventExternalTopicList{})
}
