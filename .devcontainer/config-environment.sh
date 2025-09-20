#!/bin/bash

make download-tools
go get
yarn

cd ui
yarn
yarn build