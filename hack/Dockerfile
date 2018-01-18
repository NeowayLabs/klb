FROM neowaylabs/klb

ENV GO_VERSION="1.9"
ENV GOROOT="/go"
ENV PATH=${PATH}:${GOROOT}/bin

RUN cd /tmp && \
    wget https://storage.googleapis.com/golang/go$GO_VERSION.linux-amd64.tar.gz && \
    tar -xf go$GO_VERSION.linux-amd64.tar.gz && \
    mkdir -p $GOROOT && \
    mv ./go/* $GOROOT
