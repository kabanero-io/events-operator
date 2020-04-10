package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
    DEFAULT_HTTPS_PORT = 9443
)

func MediatorHashKey(mediator *EventMediator) string {
    return mediator.TypeMeta.APIVersion + "/" + mediator.TypeMeta.Kind + "/" + mediator.ObjectMeta.Namespace + "/" + mediator.ObjectMeta.Name
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EventMediatorSpec defines the desired state of EventMediator
type EventMediatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
    CreateListener  bool `json:"createListener,omitempty"`
    ListenerPort int32    `json:"listenerPort,omitempty"`
    CreateRoute    bool `json:"createRoute,omitempty"`
    Repositories *[]EventRepository `json:"repositories,omitempty"`

    // ImportMediations  *[]string `json:"importMediations,omitempty"` // default is to import everything unless code is specified
    Mediations *[]EventMediationImpl `json:"mediations,omitempty"`
    // Functions *[]EventFunctionImpl `json:"functions,omitempty"`
}

type EventRepository struct {
    Github *EventGithubRepository `json:"github,omitempty"`
}

type EventGithubRepository struct {
    Secret string `json:"secret,omitempty"`
}


// type MediationsImpl struct {
//     Mediation *EventMediationImpl `json:"mediation,omitempty"`
//     Function *EventFunctionImpl `json:"function,omitempty"`
// }

/* Valid combinations are:
  1) assignment
  2) if and assignment
  3) if and body
  4) switch
  5) if and switch
  TBD: switch and default
*/
type EventStatement struct {
    If    *string `json:"if,omitempty"`
    Assign  *string `json:"=,omitempty"`
    Switch  *[]EventStatement `json:"switch,omitempty"`
    Body *[]EventStatement `json:"body,omitempty"`
    Default *[]EventStatement `json:"default,omitempty"`
}

type EventFunctionImpl struct {
    Name string `json:"name"`
    Input string `json:"input"`
    Output string `json:"output"`
    Body []EventStatement `json:"body"`
}

type EventMediationImpl  struct {
    Name string `json:"name"`
    // Input string `json:"input,omitempty"`
    SendTo []string `json:"sendTo,omitempty"`
    Selector *EventMediationSelector `json:"selector,omitempty"`
    Variables *[]EventMediationVariable `json:"variables,omitempty"`
    Body []EventStatement `json:"body,omitempty"`
}

type EventMediationVariable struct {
    Name string `json:"name"`
    Value *string `json:"value,omitempty"` // value treated as tring
    ValueExpression *string `json:"valueExpression,omitempty"`     // value intrepreted as CEL expression
}

type EventMediationSelector struct {
    UrlPattern string `json:"urlPattern,omitempty"`
    RepositoryType *EventMediationRepositoryType `json:"repositoryType,omitempty"`
}

type EventMediationRepositoryType struct {
    File string `json:"file"`
    NewVariable string `json:"newVariable"`
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
