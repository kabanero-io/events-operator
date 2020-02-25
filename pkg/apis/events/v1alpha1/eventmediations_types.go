package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EventMediationsSpec defines the desired state of EventMediations
type EventMediationsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
    Mediations []MediationsImpl `json:"mediations"`
}

type MediationsImpl struct {
    Mediation *EventMediationImpl `json:"mediation"`
    Function *EventFunctionImpl `json:"function"`
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
    Switch  *[]EventStatement `json:"switch,omitempty"`
    Body *[]EventStatement `json:"body,omitempty"`
}


type  EventFunctionImpl struct {
    Name string `json:"name"`
    Input string `json:"input"`
    Output string `json:"output"`
    Body []EventStatement `json:"body"`
}

type EventMediationImpl  struct {
    Name string `json:"name"`
    SendTo []string `json:"sendTo,omitempty"`
    Body []EventStatement `json:"body,omitempty"`
}

// EventMediationsStatus defines the observed state of EventMediations
type EventMediationsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventMediations is the Schema for the eventmediations API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=eventmediations,scope=Namespaced
type EventMediations struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventMediationsSpec   `json:"spec,omitempty"`
	Status EventMediationsStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventMediationsList contains a list of EventMediations
type EventMediationsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventMediations `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EventMediations{}, &EventMediationsList{})
}
