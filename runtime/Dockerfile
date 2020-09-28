FROM node
# registry.redhat.io/rhscl/nodejs-10-rhel7:latest

RUN  mkdir /app
WORKDIR /app

EXPOSE 8080

RUN curl https://mirror.openshift.com/pub/openshift-v4/clients/ocp/4.3.9/openshift-client-linux.tar.gz --output /tmp/client.tar.gz && \
    tar -xzf /tmp/client.tar.gz && \
    mv oc /usr/local/bin/oc && \
    mv kubectl /usr/local/bin/kubectl  && \
    rm -f /tmp/client.tar.gz README.md

COPY . /app

RUN npm install

RUN mkdir /check && mkdir /install

RUN curl https://get.helm.sh/helm-v2.16.8-linux-amd64.tar.gz --output /tmp/helm.tar.gz && \
    tar -xzf /tmp/helm.tar.gz && \
    mv linux-amd64/helm /usr/local/bin/helm && \
    rm -f /tmp/helm.tar.gz && rm -rf linux-amd64 && \
    apt-get install git

#CMD ["node", "app.js"]
CMD ["bash","start.sh"]
