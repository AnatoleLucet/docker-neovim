#!/bin/bash

# Abort if the image already exists on Docker Hub
export IS_EXISTENT=$([ $(curl --silent -f -lSL https://index.docker.io/v1/repositories/anatolelucet/neovim/tags/${TAG:-$TARGET} 2> /dev/null) ] && echo true || echo false)
[ "$ALLOW_OVERRIDE" = "false" ] && [ "$IS_EXISTENT" = "true" ] && echo "Cannot override the existing image. Exiting." && exit 0 || true

export IMAGE_NAME=anatolelucet/neovim:${TAG:-$TARGET}

docker build -f Dockerfile --build-arg TARGET=${TARGET:-TAG} -t $IMAGE_NAME .

docker login -u anatolelucet -p $PASSWORD

docker push $IMAGE_NAME
