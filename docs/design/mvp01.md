# Minimum Viable Product (MVP)

The team held a webex on 11 Feb 2020 and made the decisions documented here regarding the Cloud Pak pre-req checker MVP.

## Targeted Cloud Paks

The MVP pre-req checker will focus on the following CloudPaks:
- [Cloud Pak for Applications](https://www.ibm.com/cloud/cloud-pak-for-applications) (Iteration 1)
- [Cloud Pak for Integration](https://www.ibm.com/cloud/cloud-pak-for-integration) (Iteration 2)
- [Cloud Pak for Automation](https://www.ibm.com/cloud/cloud-pak-for-automation) (Iteration 2)
- [Cloud Pak for Multicloud Management](https://www.ibm.com/cloud/cloud-pak-for-management) (Iteration 2)

All of the above Cloud Paks are available on Openshift Container Platform version 4.2 and higher.  Cloud Paks that are not available on at least OCP 4.2 will not be supported by the Cloud Pak pre-req checker.

## Prerequisite focus areas

For a given Openshift Container Platform (OCP) cluster, the prerequisite requirements will be focused on the following items:
- OCP version compatibility (Iteration 1)
- OCP Cluster health (Iteration 1)
- Worker node architecture compatibility (Iteration 1)
- Worker node quantity (Iteration 1)
- Worker node compatibility (Iteration 1)
  - CPU
  - Memory
  - Nominal network speed
  - Operating System distribution and version compatibility
  - Software module requirements, version compatibility
- Networking (Iteration 1)
  - Internet connection
  - IBM Entitled Registry (access)
- Storage (Iteration 1)
  - Type
  - Size
  - StorageClass
- Common Services  (Iteration 2)
  - ELK (exists, and version)
  - Prometheus (exists, and version)
  - Version compatibility


# User experience

- Install as an Operator
  - This defines a CRD
- Command line queries
  - oc describe cpmgr --cpname CP4A
  - oc status cpmgr  --cpname CP4A --command prereq-check
  - oc get cmpgr --list
