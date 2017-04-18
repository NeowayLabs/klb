FROM golang:1.8

RUN apt-get update
RUN apt-get install -y unzip python-pip

RUN mkdir -p /tmp/nodejsinstall && \
    cd /tmp/nodejsinstall && \
    wget https://nodejs.org/dist/v6.9.4/node-v6.9.4.tar.gz && \
    tar xfv node-v6.9.4.tar.gz && \
    cd node-v6.9.4 && \
    ./configure && \
    make install && \
    rm -rf /tmp/nodejsinstall

ENV NASH_VERSION=0.3
RUN wget https://github.com/NeowayLabs/nash/releases/download/v${NASH_VERSION}/nash -O /usr/bin/nash && \
    chmod +x /usr/bin/nash

ENV NASHPATH=/root/.nash

COPY ./aws ${NASHPATH}/lib/klb/aws
COPY ./azure ${NASHPATH}/lib/klb/azure
