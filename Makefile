LICENSE_DIR=./licenses/
BUILD_DIR=./build/
DOCKER_DIR=./docker/
SHELL := /bin/bash

test: test-coverage test-race
check: additional-checks check-swagger
build-all: build-zip build-docker

require-version:
	if [ -n ${VERSION} ] && [[ $$VERSION == "v"* ]]; then echo "The version may not start with v" && exit 1; fi
	if [ -z ${VERSION} ]; then echo "Need to set VERSION" && exit 1; fi;

test-race:
	go test -v -race ./...

test-coverage:
	echo "" > coverage.txt
	for d in $(shell go list ./... | grep -v vendor); do \
		go test -v -coverprofile=profile.out -covermode=atomic $$d ; \
		if [ -f profile.out ]; then  \
			cat profile.out >> coverage.txt ; \
			rm profile.out ; \
		fi \
	done

additional-checks:
	go vet ./...
	megacheck ./...
	gocyclo -over 10 $(shell find . -iname '*.go' -type f | grep -v /vendor/)
	golint -set_exit_status $(shell go list ./... | grep -v mock)

download-tools:
	go get github.com/golang/lint/golint
	go get honnef.co/go/tools/cmd/megacheck
	go get github.com/fzipp/gocyclo
	go get -u github.com/gobuffalo/packr/...
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get github.com/karalabe/xgo

update-swagger:
	swagger generate spec --scan-models -o docs/spec.json

check-swagger: update-swagger
## add the docs to git, this changes line endings in git, otherwise this does not work on windows
	git add docs
	if [ -n "$(shell git status --porcelain | grep docs)" ]; then \
        git status --porcelain | grep docs; \
        echo Swagger or the Packr file is not up-to-date; \
        exit 1; \
    fi

extract-licenses:
	mkdir ${LICENSE_DIR} || true
	for LICENSE in $(shell find vendor/* -name LICENSE | grep -v monkey); do \
		DIR=`echo $$LICENSE | tr "/" _ | sed -e 's/vendor_//; s/_LICENSE//'` ; \
        cp $$LICENSE ${LICENSE_DIR}$$DIR ; \
    done

build-binary: require-version
	mkdir build || true
	docker pull karalabe/xgo-latest;
	xgo -ldflags "-X main.Version=${VERSION} \
				-X main.BuildDate=$(shell date "+%F-%T") \
				-X main.Commit=$(shell git rev-parse --verify HEAD)" \
				-targets linux/arm64,linux/amd64,linux/arm-7,windows-10/amd64 \
				-dest ${BUILD_DIR} \
				-out gotify \
				github.com/gotify/server

build-zip: build-binary extract-licenses
	for BUILD in $(shell find ${BUILD_DIR}*); do \
       zip -j $$BUILD.zip $$BUILD ./LICENSE; \
       zip -ur $$BUILD.zip ${LICENSE_DIR}; \
    done

build-docker: require-version build-binary
	cp ${BUILD_DIR}gotify-linux-amd64 ./docker/gotify-app
	(cd ${DOCKER_DIR} && docker build -t gotify/server:latest -t gotify/server:${VERSION} .)
	rm ${DOCKER_DIR}gotify-app
	cp ${BUILD_DIR}gotify-linux-arm-7 ./docker/gotify-app
	(cd ${DOCKER_DIR} && docker build -f Dockerfile.arm7 -t gotify/server-arm7:latest -t gotify/server-arm7:${VERSION} .)
	rm ${DOCKER_DIR}gotify-app

.PHONY: test-race test-coverage test additional-checks verify-swagger check download-tools update-swagger build-binary build-zip build-docker build-all