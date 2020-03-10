# Technical design

Operator is the framework to be used
  - Common pattern for OCP deployments
  - provides for evolution of capability


Operator type:
  - Go based
  - Ansible is too restrictive, overhead

Use Operator SDK

Define the Operator states:
- For now, read-only is sufficient
- Beyond "iteration 1" we would move to state modification

What is the object that the Operator is managing?
Abstract representation of a given Cloud Pak (CP).

Pod backs the CRD.

Operator is deployed as a pod that has a controller running in the container in the pod.
Operator runs a controller "loop"

Pre-req check - one reconcile function
Use oc commands to get information about the worker nodes

Prior to installation? (CP Status is uninstalled/"not installed")
What state of the CP gives information concerning pre-reqs?

For network connected OCP clusters get CP requirements from git repo.

Pull down yaml/json definition based on CP name, version.
Stick yaml/json in a ConfigMap (Operator can get from the ConfigMap or git repo)

For iteration 1 hard card the CP requirements into the operator. (First 2 week sprint)

# How do we get this info
Supporting only OCP 4+

- OCP version compatibility - From OCP foundation
- OCP Cluster health - Use "cluster health" pod Noel posted (Evaluate if this is in the first sprint)
- Worker node architecture compatibility (intel, power, os390) (oc get node)
- Worker node quantity (oc get nodes)
- Worker node compatibility
  - CPU
  - Memory
  - Nominal network speed
  - Operating System distribution and version compatibility
  - Software module requirements, version compatibility
- Networking
  - Internet connection
  - IBM Entitled Registry (access)
- Storage (oc get storageclass)
  - Type
  - Size
  - StorageClass (oc get sc)


We should be able to get all the info we need from the `oc get nodes` REST API.  TODO - Track down example of doing this in Go.
See [](https://docs.openshift.com/container-platform/4.1/nodes/nodes/nodes-nodes-viewing.html)

CRD status field:
Status will be Ready or NotReady

Detailed status:
"installation-readiness" - can be long string, json/yaml, or whatever.
Or events?

Create CP Objects manually?
Creating a CR - `oc apply` for each CP and version
