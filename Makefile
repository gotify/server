LICENSE_DIR=./licenses/
BUILD_DIR=./build/
DOCKER_DIR=./docker/
SHELL := /bin/bash

test: test-coverage test-race
check: additional-checks check-swagger

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
	gocyclo -over 10 $(shell find . -iname '*.go' -type f | grep -v /vendor/)
	golint -set_exit_status $(shell go list ./... | grep -v mock)

download-tools:
	go get github.com/golang/lint/golint
	go get github.com/fzipp/gocyclo
	go get -u github.com/gobuffalo/packr/...
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

update-swagger-spec:
	swagger generate spec --scan-models -o docs/spec.json

update-swagger: update-swagger-spec
	(cd docs && packr)

check-swagger: update-swagger-spec
## add the docs to git, this changes line endings in git, otherwise this does not work on windows
	git add docs
	if [ -n "$(shell git status --porcelain | grep docs)" ]; then \
        git status --porcelain | grep docs; \
        echo Swagger Spec is not up-to-date; \
        exit 1; \
    fi

extract-licenses:
	mkdir ${LICENSE_DIR} || true
	for LICENSE in $(shell find vendor/* -name LICENSE | grep -v monkey); do \
		DIR=`echo $$LICENSE | tr "/" _ | sed -e 's/vendor_//; s/_LICENSE//'` ; \
        cp $$LICENSE ${LICENSE_DIR}$$DIR ; \
    done

package-zip: extract-licenses
	for BUILD in $(shell find ${BUILD_DIR}*); do \
       zip -j $$BUILD.zip $$BUILD ./LICENSE; \
       zip -ur $$BUILD.zip ${LICENSE_DIR}; \
    done

build-docker: require-version
	cp ${BUILD_DIR}gotify-linux-amd64 ./docker/gotify-app
	(cd ${DOCKER_DIR} && docker build -t gotify/server:latest -t gotify/server:${VERSION} .)
	rm ${DOCKER_DIR}gotify-app
	cp ${BUILD_DIR}gotify-linux-arm-7 ./docker/gotify-app
	(cd ${DOCKER_DIR} && docker build -f Dockerfile.arm7 -t gotify/server-arm7:latest -t gotify/server-arm7:${VERSION} .)
	rm ${DOCKER_DIR}gotify-app

.PHONY: test-race test-coverage test additional-checks verify-swagger check download-tools update-swagger package-zip build-docker