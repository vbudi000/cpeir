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

CMD ["node", "app.js"]
