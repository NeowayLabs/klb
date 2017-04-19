FROM ubuntu:17.04

RUN apt-get update

RUN apt-get install -y curl

RUN curl -sL https://deb.nodesource.com/setup_6.x | bash -

RUN apt-get install -y nodejs python3 python3-pip libffi-dev golang-go

RUN ln -s /usr/bin/python3 /usr/bin/python

RUN apt-get install -y libffi-dev libssl-dev wget jq

RUN pip3 install -U pip

RUN pip3 install azure-cli && \
    pip3 install awscli && \
    npm install --no-optional -g azure-cli

ENV NASH_VERSION=0.3
RUN wget https://github.com/NeowayLabs/nash/releases/download/v${NASH_VERSION}/nash -O /usr/bin/nash && \
    chmod +x /usr/bin/nash

ENV NASHPATH=/root/.nash

COPY ./aws ${NASHPATH}/lib/klb/aws
COPY ./azure ${NASHPATH}/lib/klb/azure
