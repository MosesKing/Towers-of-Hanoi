package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TowerChallengeSpec defines the desired state of TowerChallenge
type TowerChallengeSpec struct {
	// Discs is the number of discs in the Tower of Hanoi challenge
	// +kubebuilder:validation:Minimum=1
	Discs int `json:"discs"`
}

// TowerChallengeStatus defines the observed state of TowerChallenge
type TowerChallengeStatus struct {
	// Steps represent the moves to solve the problem, formatted as a series of instructions
	Steps   []string `json:"steps,omitempty"`
	Message string   `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TowerChallenge is the Schema for the towerchallenges API
type TowerChallenge struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TowerChallengeSpec   `json:"spec,omitempty"`
	Status TowerChallengeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TowerChallengeList contains a list of TowerChallenge
type TowerChallengeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TowerChallenge `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TowerChallenge{}, &TowerChallengeList{})
}
