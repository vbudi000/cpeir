#!/bin/bash

entitlement=$1
# Installing

# Create entitelemt key -> imagePullSecrets
oc create secret docker-registry icr-io \
    --docker-username=cp \
    --docker-server="cp.icr.io" \
    --docker-password=${entitlement} \
    --docker-email="vbudi@us.ibm.com" \
    -n kube-system

oc patch serviceaccount default -p '{"imagePullSecrets": [{"name": "icr-io"}]}' -n kube-system

oc get secret tiller-secret -n kube-system -o yaml | grep -A5 '^data:' |awk -F: '{system("echo "$2" | base64 -d >"$1)}'
helm init --client-only
cp ca.crt ~/.helm/ca.pem
cp tls.key ~/.helm/key.pem
cp tls.crt ~/.helm/cert.pem

# collect cluster info
icpconsole=$(oc get configmap ibmcloud-cluster-info -n kube-public -o=jsonpath='{.data.cluster_address}')
icpproxy=$(oc get configmap ibmcloud-cluster-info -n kube-public -o=jsonpath='{.data.proxy_address}')

# Create configuration files
helm repo add entitled https://raw.githubusercontent.com/IBM/charts/master/repo/entitled/index.yaml

echo helm install --tls ibm-cloud-appmgmt-prod -n ibmcloudappmgmt --namespace kube-system \
      --set global.ingress.domain=${icpconsole}, \
            global.ingress.port=443, \
            global.icammcm.domain=${icpproxy}, \
            global.masterIP=${icpconsole}, \
            global.masterPort=443

# loop to check installation
# Install features

echo "{\"rc\":0}"
## cleanup

#oc delete configmap cp4multicluster-1.3.0-configmap
#oc delete secret icr-io
exit
