#!/bin/bash

make download-tools
go get

cd ui
yarn
