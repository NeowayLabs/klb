#!/bin/bash

set -o nounset
set -o errexit
echo "install python2-devel, gcc and redhat-rpm-config packages"
sudo dnf -y install python2-devel gcc redhat-rpm-config
echo "installing python pip and nodejs npm packages"
sudo dnf -y install python2-pip npm
echo "installing jq packages"
sudo dnf -y install jq
echo "installing azure cli 1.0"
sudo npm install --no-optional -g azure-cli
echo "installing azure cli 2.0"
sudo pip install azure-cli
echo "installing aws cli"
sudo pip install awscli
echo "installing nash and nashfmt"
go get -u github.com/NeowayLabs/nash/cmd/nash
go get -u github.com/NeowayLabs/nash/cmd/nashfmt
