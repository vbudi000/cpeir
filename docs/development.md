# Cloud Pak LifeCycle Operator - cpeir project

## Overview

The IBM Cloud Pak lifecycle management operator will run as native extension to Kuberentes (OpenShift) to provide insight and manage IBM Cloud Pak components.
This operator will:

  - manage a Cloud Pak deployment
  - check for requisite resource requirements
  - verify the installation
  - do stability check
  - check for updates

For more details see [Technical Design](docs/design/technical-design.md)

[Minimum Viable Product (MVP) Iteration 1](docs/design/mvp01.md)

## Requirements

- [git](https://git-scm.com/downloads)
- [go](https://golang.org/dl/) version v1.13+.
- [mercurial](https://www.mercurial-scm.org/downloads) version 3.9+
- docker version 17.03+ or podman v1.2.0+ or buildah v1.7+
- kubectl version v1.12.0+ or oc version 4.2+
- operator-sdk `brew install operator-sdk`
- Access to a Kubernetes v1.12.0+ cluster.

## Creation

This repository is created this way:

1. Create a git repo and clone it; change to the git path

	```
	git clone git@github.ibm.com:CASE/cpeir.git
	cd cpeir
	```

2.  Create a new operator:

	```
	operator-sdk new cpeir --repo github.ibm.com/CASE/cpeir
	```

3. Create the API and controller:

	```
	cd cpeir
	operator-sdk add api --api-version=cloud.ibm.com/v1alpha1 --kind=CPeir
	operator-sdk add controller --api-version=cloud.ibm.com/v1alpha1 --kind=CPeir
	```

	![sdk-run](images/sdk-run.png)

4. Modify the API part (the Custom Resource object definition)

	```
	vi pkg/apis/cloud/v1alpha1/cpeir_types.go
	```

	**Note**:
	- The path cloud/v1alpha1 is from the api-version definition
	- The objecct structure that must be modified are `CPeirSpec` and `CPeirStatus`
	- Spec represents things that you define (yaml input), Status represents things that the system observes (result from the operator)

	```go
	// CPeirSpec defines the desired state of CPeir
	type CPeirSpec struct {
	        // +kubebuilder:validation:Enum=cp4application;cp4integration;cp4automation;cp4multicloud
          CPType string `json:"cptype"`
          CPVersion string `json:"cpversion"`
  				CPSizeType string `json:"cpsizetype",omitempty`
          CPFeatures []string `json:"cpfeatures"`
	}

	// CPeirStatus defines the observed state of CPeir
	type CPeirStatus struct {
	        // +kubebuilder:validation:Enum=Initial;NotInstallable;ReadyToInstall;Installed;ValidationFailed;UpgradeAvailable
          ClusterStatus string `json:"clusterStatus"`
  				CPReqCPU resource.Quantity `json:"cpreqcpu"`
  				CPReqMemory resource.Quantity `json:"cpreqmemory"`
  				CPReqStorage resource.Quantity `json:"cpreqstorage",omitempty`
  				ClusterCPU resource.Quantity `json:"clustercpu"`
  				ClusterMemory resource.Quantity `json:"clustermemory"`
  				ClusterStorage resource.Quantity `json:"clusterstorage",omitempty`
          ClusterArch string `json:"clusterarch,omitempty"`
  				ClusterKubelet string `json:"clusterversion,omitempty"`
  				StatusMessages string `json:"statusMessages",omitempty`
          InstalledFeatures []string `json:"installedFeatures,omitempty"`
	}
	```
5. Generate kubernetes definition

	```
	operator-sdk generate k8s
	```

6. Modify CRD definition in-accordance with the API definition:

	```
	operator-sdk generate crds
  vi deploy/crds/cloud.ibm.com_cpeir_crd.yaml
  ```

	```yaml
	spec:
	  description: CPeirSpec defines the desired state of CPeir
	  properties:
	    cpfeatures:
	      items:
	        type: string
	      type: array
	    cpsizetype:
	      type: string
	    cptype:
	      type: string
	    cpversion:
	      type: string
	  required:
	  - cptype
	  - cpversion
	  type: object
	status:
	  description: CPeirStatus defines the observed state of CPeir
	  properties:
	    clusterStatus:
	      enum:
	      - Initial
	      - NotInstallable
	      - ReadyToInstall
	      - Installed
	      - ValidationFailed
	      - UpgradeAvailable
	      type: string
	    clustercpu:
	      type: string
	    clustermemory:
	      type: string
	    clusterstorage:
	      type: string
	    cpreqcpu:
	      type: string
	    cpreqmemory:
	      type: string
	    cpreqstorage:
	      type: string
	    installedFeatures:
	      items:
	        type: string
	      type: array
	    statusMessages:
	      type: string
	  required:
	  - clusterStatus
	  - statusMessages
	  type: object
	```

7. Modify the CR sample object:

	```
	vi deploy/crds/cloud.ibm.com_v1alpha1_cpeir_cr.yaml
	```

	```yaml
  apiVersion: cloud.ibm.com/v1alpha1
  kind: CPeir
  metadata:
    name: cp4application
  spec:
    cptype: "cp4application"
    cpversion: "4.0"
    cpfeatures:
      - transadv
      - kabanero
	```

8. Modify the controller program:

	```
	vi pkg/controller/cpeir/cpeir_controller.go
	```

9. Create clusterrole.yaml and clusterrole_binding.yaml to add cluster-wide roles for the operator service account.

  ```yaml
  apiVersion: rbac.authorization.k8s.io/v1
  kind: ClusterRole
  metadata:
    name: cpeir
  rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "watch", "list"]
  ```
  and
  ```yaml
  kind: ClusterRoleBinding
  apiVersion: rbac.authorization.k8s.io/v1
  metadata:
    name: cpeir
  subjects:
  - kind: ServiceAccount
    name: cpeir
    namespace: cpeir
  roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: cpeir
  ```

10. Create configMap to host the requirement sizes.

  ```yaml
  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: cpeir-config
  data:
    cp4application-4.0.yaml: |
      requirements:
        default:
          cpu: 1
          memory: 1Gi
          pv: 3Gi
    transadv-4.0.yaml: |
      requirements:
        default:
          cpu: 2
          memory: 3584Mi
          pv: 8Gi
    kabanero-4.0.yaml: |
      requirements:
        default:
          cpu: 8
          memory: 2048Mi
          pv: 25Gi
    cp4multicloud-1.2.0.yaml: |
      requirements:
        default:
          cpu: 8000m
          memory: 16Gi
        development:
          cpu: 8000m
          memory: 16Gi
          pv: 20Gi
          disk: 100Gi
        minimal:
          cpu: 16
          memory: 32Gi
          pv: 20Gi
          disk: 100Gi
        standard:
          cpu: 32
          memory: 60Gi
          pv: 20Gi
          disk: 300Gi
        enterprise:
          cpu: 54
          memory: 97Gi
          pv: 60Gi
          disk: 700Gi
    icam-1.2.0.yaml: |
      requirements:
        development:
          cpu: 10
          memory: 32Gi
          pv: 75Gi
        minimal:
          cpu: 35
          memory: 55Gi
          pv: 1200Gi
        standard:
          cpu: 90
          memory: 180Gi
          pv: 7000Gi
        enterprise:
          cpu: 105
          memory: 230Gi
          pv: 10500Gi
    cam-1.2.0.yaml: |
      requirements:
        development:
          cpu: 12
          memory: 20Gi
          pv: 65Gi
        minimal:
          cpu: 12
          memory: 30Gi
          pv: 65Gi
        standard:
          cpu: 15
          memory: 48Gi
          pv: 65Gi
        enterprise:
          cpu: 18
          memory: 60Gi
          pv: 65Gi
    endpoint-1.2.0.yaml: |
      requirements:
        default:
          cpu: 1325m
          memory: 2390Mi
    cp4multicloud-1.3.0.yaml: |
      requirements:
        default:
          cpu: 3000m
          memory: 10Gi
  ```

For information on implementing this operator, see [installation](installation.md).
