#!/usr/bin/env bash

swagger generate spec --scan-models -o docs/spec.json
(cd docs && packr && git add .)
if [[ `git status --porcelain | grep docs` ]]; then
    echo Swagger or the Packr file is not up-to-date
    exit 1
fi