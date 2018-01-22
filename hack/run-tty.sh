#!/bin/bash

set -o errexit
set -o nounset

source ./hack/loadenv.sh

docker run --rm -ti -v `pwd`:$WORKDIR -w $WORKDIR --env-file $DOCKER_ENV neowaylabs/klbdev "$@"
