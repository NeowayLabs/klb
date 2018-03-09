FROM neowaylabs/klb

ENV GO_VERSION="1.9"
ENV GOROOT="/goroot"
ENV PATH=${PATH}:${GOROOT}/bin

RUN cd /tmp && \
    wget https://storage.googleapis.com/golang/go$GO_VERSION.linux-amd64.tar.gz && \
    tar -xf go$GO_VERSION.linux-amd64.tar.gz && \
    mkdir -p $GOROOT && \
    mv ./go/* $GOROOT

COPY ./tools/azure/createsp.sh ${NASHPATH}/bin/azure-createsp.sh
COPY ./tools/azure/getcredentials.sh ${NASHPATH}/bin/azure-getcredentials.sh

COPY ./aws ${NASHPATH}/lib/klb/aws
COPY ./azure ${NASHPATH}/lib/klb/azure

ENV PATH $PATH:${NASHPATH}/bin
