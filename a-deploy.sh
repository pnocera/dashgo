#!/bin/sh

APP_NAME=$1
OLDVERSION=$2
NEWVERSION=$3
DOCKER_REGISTRY=$4
ESCAPED_OLDVERSION=$(printf '%s\n' "$OLDVERSION" | sed -e 's/[]\/$*.^[]/\\&/g');
ESCAPED_NEWVERSION=$(printf '%s\n' "$NEWVERSION" | sed -e 's/[]\/$*.^[]/\\&/g');

cd deploy
sed -i "s/$ESCAPED_OLDVERSION/$ESCAPED_NEWVERSION/g" *
cd ..
kc apply -f deploy/