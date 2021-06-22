#!/bin/bash

set -e

# Abort if the image already exists on Docker Hub
export IS_EXISTENT=$([ $(curl --silent -f -lSL https://index.docker.io/v1/repositories/anatolelucet/neovim/tags/${TAG:-$TARGET} 2> /dev/null) ] && echo true || echo false)
[ "$ALLOW_OVERRIDE" = "false" ] && [ "$IS_EXISTENT" = "true" ] && echo "Cannot override the existing image. Exiting." && exit 0 || true

export ALPINE_IMAGE_NAME=anatolelucet/neovim:${TAG:-$TARGET}
export UBUNTU_IMAGE_NAME=anatolelucet/neovim:${TAG:-$TARGET}-ubuntu

docker build -f alpine/Dockerfile --build-arg TARGET=${TARGET:-TAG} -t $ALPINE_IMAGE_NAME .
docker build -f ubuntu/Dockerfile --build-arg TARGET=${TARGET:-TAG} -t $UBUNTU_IMAGE_NAME .

docker login -u anatolelucet -p $PASSWORD

docker push $ALPINE_IMAGE_NAME
docker push $UBUNTU_IMAGE_NAME
