# Cloud Pak pre-req checker design considerations

For MVP description See [Minimum Viable Product](mvp01.md)

For some ideas see [Kubernetes Healthcheck](https://github.com/emrekenci/k8s-healthcheck)


# Cloud Pack Life Cycle Manager (cpmgr) Operator

TODO: (Flesh this out.  What is here is preliminary; based on Slack posting by Budi Darmawan; based on conversation with Gang Chen)

Cloud Pak manager operator would:
- manage a Cloud Pak deployment
- check for requisite resource requirements
- verify the installation
- do stability check
- check for updates

The Operator manages a set of resources (cpmgr).  The `cpmgr` would have properties:
- type, e.g., CP4I, CP4A, CP4D
- version
- features (list of features)

Type and version would determine the type of container image the `cpmgr` would run.  (TBD - Need more detail here. Container for what?  The Cloud Pak pods?) The particular container image then performs life-cycle management.

Interaction with `cpmgr`:
- Status - state of the given Cloud Pak: Installable, Running, UpgradeAvailable, etc.
- Verbose Status - provides a status for all pods associated with the given Cloud Pak
- Describe - The describe output for the Cloud Pak.  (Is there such a thing?) (Or if a pod name is provided then the describe output for that cloud pak pod?  Maybe that's stepping too much into kubectl.)
- Events - Event stream associated with the Cloud Pak. (Is there such a thing?)
- Log(s) - Log associated with the Cloud Pak.

Concerns:
- Has the scope broadened from "pre-req checker" to life cycle manager?

## Cloud Pak Manager (cpmgr)

- The `cpmgr` pod has an endpoint that supports interaction with external clients.
- The `cpmgr` pod has an endpoint that supports interaction with helper pods (`cpmgr-helper`) that need to register with the `cpmgr`.  The helper pods run on each worker node and provide information to the `cpmgr` regarding the node, e.g., configured resources and resource availability.

TBD - `cpmgr` API.

# Some details on pre-req checking

- Deploy a `cpmgr-helper` pod to each worker node.  For purposes of pre-req checking, the `cpmgr-helper` emits a report on the resources configured on that node. First cut is to report on the total resources deployed, e.g., `lscpu`, `/proc/meminfo`, `df -h`. Second cut would be more sophisticated and report on resource availability, e.g., `iostat`; `vmstat -s`; `free -m`  The resource availability measure is dynamic and more accurately reflects the resources available for a given Cloud Pak deployment.

See [8 commands to check cpu information on Linux](https://www.binarytides.com/linux-cpu-information/)
See [5 commands to check memory usage on Linux](https://www.binarytides.com/linux-command-check-memory-usage/)

The `cpmgr-helper` pod reports its endpoint to a `cpmgr` pod.  The `cpmgr` pod thus can poll each `cpmgr-helper` for worker node info as part of the pre-req check before installation.

The `cpmgr-helper` endpoint is configurable (using an env var).
The `cpmgr-helper` needs to know the endpoint of the `cpmgr` (configured using an env var).

# Implementation considerations

- We should likely use the Operator SDK from Red Hat.
- Given the language support of the Operator SDK, we should implement the operator in Go.

- Command line interaction with the Operator?  (I would vote for Python as it has a rich set of libraries.)
