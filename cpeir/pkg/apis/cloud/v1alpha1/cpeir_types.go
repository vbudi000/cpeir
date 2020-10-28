package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CPeirSpec defines the desired state of CPeir
type CPeirSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
        CPType string `json:"cptype"`
        CPVersion string `json:"cpversion"`
				CPSizeType string `json:"cpsizetype",omitempty`
				StorageClass string `json:"storageClass",omitempty`
				// +kubebuilder:validation:Enum=Check;Install;HealthCheck;Upgrade
				Action string `json:"action"`
        CPFeatures []CPeirFeature `json:"cpfeatures"`
}

type CPeirFeature struct {
	// Feature name
	Name string `json:"name"`
	// StorageClass name if needed
	StorageClass string `json:"storageClass",omitempty`
	// Determine ha requirements
	//HaEnabled bool `json:"haEnabled"`
}

type CPeirCPReq struct {
	// CPU installation requirement
	CPReqCPU resource.Quantity `json:"cpreqcpu"`
	// Memory installation requirement
	CPReqMemory resource.Quantity `json:"cpreqmemory"`
	// Storage (PVC) installation requirement
	CPReqStorage resource.Quantity `json:"cpreqstorage",omitempty`
}

type CPeirCluster struct {
	// Allocatable CPU in worker nodes
	ClusterCPU resource.Quantity `json:"clusterCpu"`
	// Allocatable memory in worker nodes
	ClusterMemory resource.Quantity `json:"clusterMemory"`
	// Available Storage - TBD - now showing 0
	ClusterStorage resource.Quantity `json:"clusterStorage",omitempty`
	// Worker node architecture
	ClusterArch string `json:"clusterArch,omitempty"`
	// Number of Worker nodes
	ClusterWorkerNum int `json:"clusterWorkerNum,omitempty"`
	// Worker node kubelet version
	ClusterKubelet string `json:"clusterVersion,omitempty"`
	// OpenShift version
	OCPVersion string `json:"ocpVersion",omitempty`
}

// CPeirStatus defines the observed state of CPeir
type CPeirStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

  			// Status of the CloudPak in the cluster
				// +kubebuilder:validation:Enum=Initial;NotInstallable;ReadyToInstall;Installed;ValidationFailed;UpgradeAvailable
        CPStatus string `json:"CPStatus"`
				// Calculated CloudPak requirements
				CPRequirement CPeirCPReq `json:"cpreq"`
				// Retrieved Cluster characteristics
				ClusterStatus CPeirCluster `json:"cluster"`
				// access to cp.icr.io
				OnlineInstall bool `json:"onlineInstall,omitempty"`
				// available image registry space -
				OfflineInstall bool `json:"offlineInstall,omitempty"`
				// Miscellaneous messages from the operator run
				StatusMessages string `json:"statusMessages",omitempty`
				// List of installed features as investigated - TBD - now empty
        InstalledFeatures []string `json:"installedFeatures,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CPeir is the Schema for the cpeirs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=cpeirs,scope=Namespaced
type CPeir struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CPeirSpec   `json:"spec"`
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
