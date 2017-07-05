#!/bin/bash

set -o nounset
set -o errexit
echo "installing python pip and nodejs npm"
sudo apt-get -y install python-pip npm jq
echo "installing azure cli 1.0"
sudo npm install --no-optional -g azure-cli
echo "installing azure cli 2.0"
sudo pip install azure-cli
echo "installing aws cli"
sudo pip install awscli
echo "installing nash"
go get -u github.com/NeowayLabs/nash
