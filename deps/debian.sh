#!/bin/bash

set -o nounset
set -o errexit
echo "installing python pip and nodejs npm"
curl -sL https://deb.nodesource.com/setup_6.x | sudo -E bash -
sudo apt-get -y install python-pip nodejs jq
echo "installing azure cli 1.0"
sudo npm install --no-optional -g azure-cli
echo "installing azure cli 2.0"
sudo pip install azure-cli
echo "installing aws cli"
sudo pip install awscli
echo "installing nash"
go get -u github.com/NeowayLabs/nash
