package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CPeirSpec defines the desired state of CPeir
type CPeirSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:validation:Enum=cp4application;cp4integration;cp4automation;cp4multicloud
        CPType string `json:"cptype"`
        CPVersion string `json:"cpversion"`
        CPFeatures []string `json:"cpfeatures"`
}

// CPeirStatus defines the observed state of CPeir
type CPeirStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	// +kubebuilder:validation:Enum=Initial;NotInstallable;ReadyToInstall;Installed;ValidationFailed;UpgradeAvailable
        ClusterStatus string `json:"clusterStatus"`
	StatusMessages string `json:"statusMessages"`
        InstalledFeatures []string `json:"installedFeatures"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CPeir is the Schema for the cpeirs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=cpeirs,scope=Namespaced
type CPeir struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CPeirSpec   `json:"spec,omitempty"`
	Status CPeirStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CPeirList contains a list of CPeir
type CPeirList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CPeir `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CPeir{}, &CPeirList{})
}
