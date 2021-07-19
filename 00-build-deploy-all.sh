#!/bin/sh


OLDVERSION="v0.0.4"
VERSION="v0.0.5"
DOCKER_REGISTRY="gcr.io/gci-ptfd-host-dev"

APP_NAME="dashgo"
./a-builddocker.sh $APP_NAME $VERSION $DOCKER_REGISTRY
./a-deploy.sh $APP_NAME $OLDVERSION $VERSION $DOCKER_REGISTRY


git add .
git commit -am $VERSION
git push origin main
git tag $VERSION
git push origin --tags