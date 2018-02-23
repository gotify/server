#!/usr/bin/env bash

if [ $# -eq 0 ]; then
    echo You need to pass the version as first parameter.
    exit 1
fi

cp ./build/gotify-linux-amd64 ./docker/gotify-app
(cd docker && docker build -t gotify/server:latest -t gotify/server:$1 .)
rm ./docker/gotify-app
cp ./build/gotify-linux-arm-7 ./docker/gotify-app
(cd docker && docker build -f Dockerfile.arm7 -t gotify/server-arm7:latest -t gotify/server-arm7:$1 .)
rm ./docker/gotify-app