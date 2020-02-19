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
    Code []EventImplementation `json:"code",omitempty`
}

type EventImplementation struct {
    Mediation *EventMediation `json:"mediation",omitempty`
    Function *EventFunction  `json:"function",omitempty`
}

type EventFunction struct {
    Name string `json:"name"`
    Input string `json:"input"`
    Output string `json:"output"`
    Body []EventStatement `json:"body"`
}


type EventMediation struct {
    Name string `json:"name",omitempty`
    SubscribeFrom [] string `json:"subscribeFrom",omitempty`
    SendTo [] string `json:"sendTo",omitempty`
    Body [] EventStatement `json:"body"`
}

/* Valid combinations are:
  1) assignment
  2) if and assignment
  3) if and body
  4) switch
  5) if and switch
  TBD: switch and default
*/
type EventStatement struct {
    If    *string `json:"if"`
    Assign  *string `json:"="`
    Switch  *[]EventStatement `json:"switch",omitempty`
    Body *[]EventStatement `json:"body",omitempty`
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
