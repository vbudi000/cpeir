# Cloud Pak EIR

The cloud pak EIR is controlled by a ConfigMap that has the requirements written in yaml; The following is a sample config for CP4application 4.0.1:

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: cpeir-config
data:
  cp4application-4.0.1.yaml: |
    system:
      cpu: 1000m
      memory: 4Gi
```

Once the EIR is deployed, the cpeir object can be shown as below:

```
$ oc describe cpeir cp4application
Name:         cp4application
Namespace:    default
Labels:       <none>
Annotations:  <none>
API Version:  cloud.ibm.com/v1alpha1
Kind:         CPeir
Metadata:
  Creation Timestamp:  2020-03-09T21:25:51Z
  Generation:          1
  Resource Version:    33908621
  Self Link:           /apis/cloud.ibm.com/v1alpha1/namespaces/default/cpeirs/cp4application
  UID:                 8d569132-624c-11ea-99fc-005056a583a0
Spec:
  Cptype:     cp4application
  Cpversion:  4.0.1
Status:
  Cluster Status:   ReadyToInstall
  Status Messages:  Allocatable worker nodes capacity is CPU=38500m and memory=178048319488
Requirement is CPU=1000m and memory=4294967296
Events:  <none>
```
