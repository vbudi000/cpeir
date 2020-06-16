FROM node

RUN mkdir /opt/app
WORKDIR /opt/app

EXPOSE 8080

RUN curl https://mirror.openshift.com/pub/openshift-v4/clients/ocp/4.3.9/openshift-client-linux.tar.gz --output /tmp/client.tar.gz && \
    tar -xzf /tmp/client.tar.gz && \
    ls -la && pwd && \
    mv oc /usr/local/bin/oc && \
    mv kubectl /usr/local/bin/kubectl  && \
    rm -f /tmp/client.tar.gz README.md

COPY . /opt/app

RUN npm install

RUN mkdir /check && mkdir /install

RUN curl https://get.helm.sh/helm-v2.16.8-linux-amd64.tar.gz --output /tmp/helm.tar.gz && \
    tar -xzf /tmp/helm.tar.gz && \
    mv linux-amd64/helm /usr/local/bin/helm && \
    rm -f /tmp/helm.tar.gz && rm -rf linux-amd64 && \
    apt-get install git

CMD ["node", "app.js"]
