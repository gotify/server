#!/usr/bin/env bash

if [ $# -eq 0 ]; then
    echo You need to pass the version as first parameter.
    exit 1
fi

VERSION=$1
DATE=`date "+%F-%T"`
COMMIT=`git rev-parse --verify HEAD`
BRANCH=`git rev-parse --abbrev-ref HEAD`
URL=github.com/gotify/server
PREFIX=gotify-$VERSION
DEST=./build
TARGETS=linux/arm64,linux/amd64,linux/arm-7,windows-10/amd64
LICENSES=./licenses/

xgo -ldflags "-X main.Version=$VERSION -X main.BuildDate=$DATE -X main.Commit=$COMMIT -X main.Branch=$BRANCH" -targets $TARGETS -dest $DEST -out $PREFIX $URL

mkdir $LICENSES
for LICENSE in $(/bin/find vendor/ -name LICENSE | grep -v monkey); do
    DIR=$(echo $LICENSE | tr "/" _ | sed -e 's/vendor_//; s/_LICENSE//')
    mkdir $LICENSES$DIR
    cp $LICENSE $LICENSES$DIR
done

for BIN in build/*; do
   zip -j $BIN.zip $BIN LICENSE
   zip -ur $BIN.zip $LICENSES
done

rm -rf $LICENSES
