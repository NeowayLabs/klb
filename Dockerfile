FROM ubuntu:17.04

RUN apt-get update

RUN apt-get install -y curl

RUN curl -sL https://deb.nodesource.com/setup_6.x | bash -

RUN apt-get install -y nodejs python3 python3-pip libffi-dev golang-go openssh-server

RUN ln -s /usr/bin/python3 /usr/bin/python

RUN apt-get install -y libffi-dev libssl-dev wget jq

RUN pip3 install -U pip

RUN pip3 install azure-cli && \
    pip3 install awscli && \
    npm install --no-optional -g azure-cli

ENV NASH_VERSION=v0.5
RUN cd /tmp && \
    wget https://github.com/NeowayLabs/nash/releases/download/${NASH_VERSION}/nash-${NASH_VERSION}-linux-amd64.tar.gz && \
    tar xfz nash-${NASH_VERSION}-linux-amd64.tar.gz && \
    cp /tmp/cmd/nash/nash /usr/bin/nash && \
    rm nash-${NASH_VERSION}-linux-amd64.tar.gz && \
    rm -rf /tmp/cmd

ENV NASHPATH=/root/.nash

COPY ./aws ${NASHPATH}/lib/klb/aws
COPY ./azure ${NASHPATH}/lib/klb/azure
