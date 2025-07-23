#!/bin/bash

set -e

# Abort if the image already exists on Docker Hub
export EXISTS=$([ $(curl --silent -f -lSL https://index.docker.io/v1/repositories/anatolelucet/neovim/tags/${TAG:-$VERSION} 2> /dev/null) ] && echo true || echo false)
[ "$ALLOW_OVERRIDE" = "false" ] && [ "$EXISTS" = "true" ] && echo "Cannot override an existing image. Exiting." && exit 0 || true

export ALPINE_IMAGE_NAME=anatolelucet/neovim:${TAG:-$VERSION}
export UBUNTU_IMAGE_NAME=anatolelucet/neovim:${TAG:-$VERSION}-ubuntu

docker login -u anatolelucet -p $PASSWORD

docker buildx build --push --platform linux/amd64,linux/arm64 -f alpine/Dockerfile --build-arg VERSION=${VERSION} -t $ALPINE_IMAGE_NAME .
docker buildx build --push --platform linux/amd64,linux/arm64 -f ubuntu/Dockerfile --build-arg VERSION=${VERSION} -t $UBUNTU_IMAGE_NAME .

