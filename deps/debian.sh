#!/bin/bash

set -o nounset
set -o errexit
echo "installing python pip and nodejs npm"
apt-get install python-pip jq
echo "installing nodejs"
nvm install 7.7.4
echo "installing azure cli 1.0"
npm install --no-optional -g azure-cli
echo "installing azure cli 2.0"
pip install azure-cli
echo "installing aws cli"
pip install awscli
echo "installing nash"
go install github.com/NeowayLabs/nash
