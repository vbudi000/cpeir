#!/bin/bash

entitlement=$1
# Installing
if [ -z $entitlement ]; then
  exit 9
fi

# Create entitelemt key -> imagePullSecrets
oc create secret docker-registry icr-io \
    --docker-username=cp \
    --docker-server="cp.icr.io" \
    --docker-password=${entitlement} \
    --docker-email="vbudi@us.ibm.com"

# collect cluster info
sc=$(oc get sc --no-headers -o custom-columns=name:metadata.name | head -n 1)
workers=$(oc get node --no-headers -o  custom-columns=name:metadata.name -l node-role.kubernetes.io/worker | paste -s -d, -)
console=$(oc get route console -n openshift-console -o jsonpath='{.spec.host}')
roks_url=${console#"console-openshift-console."}
if [[ $roks_url == apps* ]]; then
  roks="false"
else
  roks="true"
fi

# Create configuration files
cat cp4multicloud-1.3.0-configMap.yaml | \
    sed "s/<entitlement>/${entitlement}/g" | \
    sed "s/<workers>/${workers}/g" | \
    sed "s/<sc>/${sc}/g" | \
    sed "s/<roks>/${roks}/g" | \
    sed "s/<roks_url>/${roks_url}/g"| oc create -f -

# Create installation Job
oc apply -f cp4multicloud-1.3.0-install.yaml

# loop to check installation
# Install features

echo "{\"rc\":0}"
## cleanup

#oc delete configmap cp4multicluster-1.3.0-configmap
#oc delete secret icr-io
exit
