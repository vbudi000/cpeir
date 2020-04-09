# Cloud Pak LifeCycle Operator - cpeir project - Usage

## Pre-requisites

To use this operator, you need the following:

- The `oc` command
- Access to the cluster on the cpeir project

## Using the operator

1. Build a YAML file for the CPeir object. The YAML must include the CloudPak type, version, installation type and list of features that are expected, the following is a sample YAML for Cloud Pak for MultiCloud Management 1.2:

  ```yaml
  apiVersion: cloud.ibm.com/v1alpha1
  kind: CPeir
  metadata:
    name: cp4multicloud
  spec:
    cptype: "cp4multicloud"
    cpversion: "1.2"
    cpsizetype: "development"
    cpfeatures:
      - icam
  ```

2. Load the yaml to OpenShift:

  ```bash
  oc create -f cp4mcm.yaml
  ```

2. The operator will perform its evaluation to your cluster and compare that to the configMap that it has and deliver its verdict in a couple of seconds. Check the newly created object:

	```bash
	oc describe cpeir cp4multicloud
	```

3. The sample output is as follows:

	```bash
  $ oc describe cpeir cp4multicloud
  Name:         cp4multicloud
  Namespace:    cpeir
  Labels:       <none>
  Annotations:  <none>
  API Version:  cloud.ibm.com/v1alpha1
  Kind:         CPeir
  Metadata:
    Creation Timestamp:  2020-04-08T20:48:59Z
    Generation:          2
    Resource Version:    6519472
    Self Link:           /apis/cloud.ibm.com/v1alpha1/namespaces/cpeir/cpeirs/cp4multicloud
    UID:                 e730de4c-0c1f-493f-a0e0-ff8c82ca418b
  Spec:
    Cpfeatures:
      icam
    Cpsizetype:  development
    Cptype:      cp4multicloud
    Cpversion:   1.2.0
  Status:
    Cluster Arch:        amd64
    Cluster Cpu:         46500m
    Cluster Memory:      195986604Ki
    Cluster Status:      ReadyToInstall
    Cluster Storage:     0
    Cluster Version:     v1.16.2
    Cluster Worker Num:  3
    Cpreqcpu:            18
    Cpreqmemory:         48Gi
    Cpreqstorage:        95Gi
    Status Messages:     Allocatable worker nodes capacity is CPU=46500m and memory=200690282496
  Requirement is CPU=18000m and memory=51539607552
  Events:  <none>
	```

  As indicated by Cluster Status field, the CloudPak can then be installed. 
