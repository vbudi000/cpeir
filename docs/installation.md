# Cloud Pak LifeCycle Operator - cpeir project - Installation

## Pre-requisites

To deploy this operator, you need the following:

- Git command line interface
- The `oc` command
- The `operator-sdk` command
- Docker repository where you will store the operator image
- Cluster admin access to the cluster that you want to install the tool

## Deploying the operator

1. Clone the git repository

	```bash
	git clone https://github.ibm.com/CASE/cpeir
	```

2. Build your operator image:

	```bash
	cd cpeir/cpeir
	operator-sdk build <namespace>/cpeir:v0.0.1
	cd ../runtime
	docker build -t <namespace>/runtime:v0.0.1 .
	```

3. Push up the image to a docker repository

	```bash
	docker push <namespace>/cpeir:v0.0.1
	docker push <namespace>/runtime:v0.0.1
	```

4. Create OpenShift resources:

	```bash
	oc login api.<cluster>.<domain>:6443
	oc new-project cpeir
	oc create -f deploy/crds/cloud.ibm.com_cpeirs_crd.yaml
	oc create -f deploy/service_account.yaml
	oc create -f deploy/clusterrole_binding.yaml
	```

4. Modify the operator.yaml with the image names you push to docker repo:

	```bash
	oc create -f deploy/operator.yaml
	```

5. Create your entitlement key secret:

	```bash
	oc create secret generic entitlement --from-literal entitlement="entitlementkey"
	```

Check that the operator pod is running, use the command `oc get pod -n cpeir` and make sure that the cpeir pod status is `Running`.

For information on using this operator, see [usage](usage.md).
