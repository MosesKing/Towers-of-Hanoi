package v1alpha1

import (
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1" // Use Crossplane's common v1 for Conditions
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TowerChallengeSpec defines the desired state of TowerChallenge
type TowerChallengeSpec struct {
	// Discs is the number of discs in the Tower of Hanoi challenge
	// +kubebuilder:validation:Minimum=1
	Discs int `json:"discs"`
}

// TowerChallengeStatus defines the observed state of TowerChallenge
type TowerChallengeStatus struct {
	// Standard condition fields used by Crossplane to report the observed state of the resource.
	Conditions []xpv1.Condition `json:"conditions,omitempty"`

	// Steps represent the moves to solve the problem, formatted as a series of instructions
	Steps   []string `json:"steps,omitempty"`
	Message string   `json:"message,omitempty"`

	// Phase represents the current phase of the operation (e.g., "Pending", "Completed")
	Phase string `json:"phase,omitempty"`
	// ConfigMapsCreated indicates whether the config maps were successfully created
	ConfigMapsCreated bool `json:"configMapsCreated"`
	// ConfigMapNames lists the names of the created config maps
	ConfigMapNames []string `json:"configMapNames,omitempty"`
	// StartTime is the time when the operation started
	StartTime metav1.Time `json:"startTime,omitempty"`
	// EndTime is the time when the operation completed
	EndTime metav1.Time `json:"endTime,omitempty"`
	// ErrorMessage contains details of any errors that occurred
	ErrorMessage string `json:"errorMessage,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="Discs",type="integer",JSONPath=".spec.discs"
//+kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
//+kubebuilder:printcolumn:name="StartTime",type="date",JSONPath=".status.startTime"
//+kubebuilder:printcolumn:name="EndTime",type="date",JSONPath=".status.endTime"

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
