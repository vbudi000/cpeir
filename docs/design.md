# Cloud Pak LifeCycle Operator - cpeir project - Design

The cpeir operator is written in Go language to collect information from the cluster and check them against a list of requriements for installing and running an IBM Cloud Pak.

## Iteration 1

The first iteration will include pre-requisite checker for IBM Cloud Pak. The checking is performed based on a pre-coded ConfigMap for various IBM Cloud Pak software and versions. The checking is performed for:

- OCP version compatibility (Iteration 1) - ClusterVersion
- OCP Cluster health (Iteration 1) - ClusterVersion
- Worker node architecture compatibility (Iteration 1) - Node info
- Worker node quantity (Iteration 1) - Node info
- Cluster size - CPU (Iteration 1) - Node info
- Cluster size - memory (Iteration 1) - Node info
- Installation requirement (Iteration 1) - TBD
  - Internet connection to IBM Entitled Registry (online install)
  - Registry space to load images (offline install)
- Storage (Iteration 2)
  - Type
  - Size
  - StorageClass
  - Software module requirements, version compatibility (Iteration 2)
- Common Services  (Iteration 2)
  - ELK (exists, and version)
  - Prometheus (exists, and version)
  - Version compatibility

### Node information

Information will be retrieved from the Node information as listed below:

```bash
$ oc get node ip-10-0-142-27.us-east-2.compute.internal -o yaml
```

```yaml
apiVersion: v1
kind: Node
metadata:
  annotations:
    machine.openshift.io/machine: openshift-machine-api/cp4mcm13-gse01-worker-us-east-2a-vdlmx
    machineconfiguration.openshift.io/currentConfig: rendered-worker-a27993eb6eaf1d884166948fda90b041
    machineconfiguration.openshift.io/desiredConfig: rendered-worker-a27993eb6eaf1d884166948fda90b041
    machineconfiguration.openshift.io/reason: ""
    machineconfiguration.openshift.io/state: Done
    volumes.kubernetes.io/controller-managed-attach-detach: "true"
  creationTimestamp: "2020-03-27T21:10:53Z"
  labels:
    beta.kubernetes.io/arch: amd64
    beta.kubernetes.io/instance-type: m4.4xlarge
    beta.kubernetes.io/os: linux
    failure-domain.beta.kubernetes.io/region: us-east-2
    failure-domain.beta.kubernetes.io/zone: us-east-2a
    kubernetes.io/arch: amd64
    kubernetes.io/hostname: ip-10-0-142-27
    kubernetes.io/os: linux
    management: "true"
    master: "true"
    node-role.kubernetes.io/cp-management: "true"
    node-role.kubernetes.io/cp-master: "true"
    node-role.kubernetes.io/cp-proxy: "true"
    node-role.kubernetes.io/worker: ""
    node.openshift.io/os_id: rhcos
    proxy: "true"
  name: ip-10-0-142-27.us-east-2.compute.internal
  resourceVersion: "6684507"
  selfLink: /api/v1/nodes/ip-10-0-142-27.us-east-2.compute.internal
  uid: 9085c6c4-fd01-433f-8d72-8beb78311f00
spec:
  providerID: aws:///us-east-2a/i-0f66979af747093b1
status:
  addresses:
  - address: 10.0.142.27
    type: InternalIP
  - address: ip-10-0-142-27.us-east-2.compute.internal
    type: Hostname
  - address: ip-10-0-142-27.us-east-2.compute.internal
    type: InternalDNS
  allocatable:
    attachable-volumes-aws-ebs: "39"
    cpu: 15500m
    ephemeral-storage: "192764845352"
    hugepages-1Gi: "0"
    hugepages-2Mi: "0"
    memory: 65328868Ki
    pods: "250"
  capacity:
    attachable-volumes-aws-ebs: "39"
    cpu: "16"
    ephemeral-storage: 209163244Ki
    hugepages-1Gi: "0"
    hugepages-2Mi: "0"
    memory: 65943268Ki
    pods: "250"
  conditions:
  - lastHeartbeatTime: "2020-04-09T04:05:30Z"
    lastTransitionTime: "2020-03-27T21:10:53Z"
    message: kubelet has sufficient memory available
    reason: KubeletHasSufficientMemory
    status: "False"
    type: MemoryPressure
  - lastHeartbeatTime: "2020-04-09T04:05:30Z"
    lastTransitionTime: "2020-03-27T21:10:53Z"
    message: kubelet has no disk pressure
    reason: KubeletHasNoDiskPressure
    status: "False"
    type: DiskPressure
  - lastHeartbeatTime: "2020-04-09T04:05:30Z"
    lastTransitionTime: "2020-03-27T21:10:53Z"
    message: kubelet has sufficient PID available
    reason: KubeletHasSufficientPID
    status: "False"
    type: PIDPressure
  - lastHeartbeatTime: "2020-04-09T04:05:30Z"
    lastTransitionTime: "2020-03-27T21:11:33Z"
    message: kubelet is posting ready status
    reason: KubeletReady
    status: "True"
    type: Ready
  daemonEndpoints:
    kubeletEndpoint:
      Port: 10250
  images:
  - names:
    - image-registry.openshift-image-registry.svc:5000/ibmcom/icp-platform-auth@sha256:329643b7042b7279071639f29e7e710bdf8f10ff895f18e7b49157c0d4686c4d
    - image-registry.openshift-image-registry.svc:5000/ibmcom/icp-platform-auth:3.3.1
    sizeBytes: 1752055944
    . . .
  nodeInfo:
    architecture: amd64
    bootID: a64d1421-2f33-46d0-917f-dc2c9649dec4
    containerRuntimeVersion: cri-o://1.16.3-28.dev.rhaos4.3.git9aad8e4.el8
    kernelVersion: 4.18.0-147.5.1.el8_1.x86_64
    kubeProxyVersion: v1.16.2
    kubeletVersion: v1.16.2
    machineID: ce5e2b57c2f14717a3efc2d309e4d739
    operatingSystem: linux
    osImage: Red Hat Enterprise Linux CoreOS 43.81.202003191953.0 (Ootpa)
    systemUUID: ec2553f8-aa90-3a67-fbfc-bbba72ca6d30
  volumesAttached:
  - devicePath: /dev/xvdbv
    name: kubernetes.io/aws-ebs/aws://us-east-2a/vol-07e74f76dd8f138cd
  - devicePath: /dev/xvdck
    name: kubernetes.io/aws-ebs/aws://us-east-2a/vol-012ba82e1500fbadb
  volumesInUse:
  - kubernetes.io/aws-ebs/aws://us-east-2a/vol-012ba82e1500fbadb
  - kubernetes.io/aws-ebs/aws://us-east-2a/vol-07e74f76dd8f138cd
```

### ClusterVersion

Information (OCP version and cluster health) will be retrieved from the Cluster Version information as listed below:

```bash
$ oc get clusterversion version -o yaml
```


```yaml
apiVersion: config.openshift.io/v1
kind: ClusterVersion
metadata:
  creationTimestamp: "2020-03-27T20:33:20Z"
  generation: 1
  name: version
  resourceVersion: "5840656"
  selfLink: /apis/config.openshift.io/v1/clusterversions/version
  uid: 0850b523-6768-46ae-9869-5a6d8214270d
spec:
  channel: stable-4.3
  clusterID: 57d86247-58ef-4695-8c99-747e3a0ca5ce
  upstream: https://api.openshift.com/api/upgrades_info/v1/graph
status:
  availableUpdates:
  - force: false
    image: quay.io/openshift-release-dev/ocp-release@sha256:f0fada3c8216dc17affdd3375ff845b838ef9f3d67787d3d42a88dcd0f328eea
    version: 4.3.9
  conditions:
  - lastTransitionTime: "2020-03-27T21:58:05Z"
    message: Done applying 4.3.8
    status: "True"
    type: Available
  - lastTransitionTime: "2020-04-03T07:06:14Z"
    status: "True"
    type: RetrievedUpdates
  desired:
    force: false
    image: quay.io/openshift-release-dev/ocp-release@sha256:a414f6308db72f88e9d2e95018f0cc4db71c6b12b2ec0f44587488f0a16efc42
    version: 4.3.8
  history:
  - completionTime: "2020-03-27T21:58:05Z"
    image: quay.io/openshift-release-dev/ocp-release@sha256:a414f6308db72f88e9d2e95018f0cc4db71c6b12b2ec0f44587488f0a16efc42
    startedTime: "2020-03-27T20:33:34Z"
    state: Completed
    verified: false
    version: 4.3.8
  observedGeneration: 1
  versionHash: lnZzahlL8hk=
```


### Network access and registry check

Network access to IBM Entitled registry

```bash
$ curl -k https://cp.icr.io -i
HTTP/1.1 301 Moved Permanently
Date: Thu, 09 Apr 2020 15:29:53 GMT
Content-Length: 0
Connection: keep-alive
Location: https://www.ibm.com/cloud/container-registry
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-XSS-Protection: 1; mode=block
```

Registry check:

```bash
$ oc get configs cluster -o yaml
```

```yaml
apiVersion: imageregistry.operator.openshift.io/v1
kind: Config
metadata:
  finalizers:
  - imageregistry.operator.openshift.io/finalizer
  name: cluster
spec:
  defaultRoute: true
  disableRedirect: false
  logging: 2
  managementState: Managed
  proxy:
    http: ""
    https: ""
    noProxy: ""
  readOnly: false
  replicas: 1
  requests:
    read:
      maxInQueue: 0
      maxRunning: 0
      maxWaitInQueue: 0s
    write:
      maxInQueue: 0
      maxRunning: 0
      maxWaitInQueue: 0s
  storage:
    s3:
      bucket: cp4mcm13-gse01-image-registry-us-east-1-yurxanjfbfmkyagcojgfvl
      encrypt: true
      keyID: ""
      region: us-east-1
      regionEndpoint: ""
status:
  conditions:
  - lastTransitionTime: "2020-03-27T20:45:22Z"
    reason: S3 Bucket Exists
    status: "True"
    type: StorageExists
  observedGeneration: 3
  readyReplicas: 0
  storage:
    s3:
      bucket: cp4mcm13-gse01-image-registry-us-east-1-yurxanjfbfmkyagcojgfvl
      encrypt: true
      keyID: ""
      region: us-east-1
      regionEndpoint: ""
  storageManaged: true
```
or an example with pvc:

```yaml
apiVersion: imageregistry.operator.openshift.io/v1
kind: Config
metadata:
  creationTimestamp: "2020-04-09T22:13:53Z"
  finalizers:
  - imageregistry.operator.openshift.io/finalizer
  generation: 4
  name: cluster
  resourceVersion: "20502"
  selfLink: /apis/imageregistry.operator.openshift.io/v1/configs/cluster
  uid: 11c5ff9c-43ee-406f-b055-02b5a1baf94a
spec:
  defaultRoute: false
  disableRedirect: false
  httpSecret: 6ea95ae10eba9f08f92563857f92f155bb027b5b52d76f815d7aa0eb478015ade1b889b1d2481008ec00441a855cf9e5d60a8fb8ae4e235887cd1ce6a79eb7e5
  logging: 2
  managementState: Managed
  proxy:
    http: ""
    https: ""
    noProxy: ""
  readOnly: false
  replicas: 1
  requests:
    read:
      maxInQueue: 0
      maxRunning: 0
      maxWaitInQueue: 0s
    write:
      maxInQueue: 0
      maxRunning: 0
      maxWaitInQueue: 0s
  storage:
    pvc:
      claim: image-registry-storage
```

```
root@utility:/var/www/html# oc exec -n openshift-image-registry image-registry-6d75cc7975-jw7zn -- df -k
Filesystem                           1K-blocks     Used Available Use% Mounted on
overlay                              209163244  8716024 200447220   5% /
tmpfs                                    65536        0     65536   0% /dev
tmpfs                                 16468080        0  16468080   0% /sys/fs/cgroup
shm                                      65536        0     65536   0% /dev/shm
tmpfs                                 16468080     5300  16462780   1% /etc/passwd
172.16.53.250:/data/registry         123329280 15993088 101028352  14% /registry
tmpfs                                 16468080        8  16468072   1% /etc/secrets
/dev/mapper/coreos-luks-root-nocrypt 209163244  8716024 200447220   5% /etc/hosts
tmpfs                                 16468080       24  16468056   1% /run/secrets/kubernetes.io/serviceaccount
tmpfs                                 16468080        0  16468080   0% /proc/acpi
tmpfs                                 16468080        0  16468080   0% /proc/scsi
tmpfs                                 16468080        0  16468080   0% /sys/firmware
```

## Iteration 2

In the iteration 2, more detailed pre-requisite will be checked, this include: storage requirements and software incompatibility checks.


### Storage check

### Software incompatibility check
