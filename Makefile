LICENSE_DIR=./licenses/
BUILD_DIR=./build
DOCKER_DIR=./docker/
SHELL := /bin/bash
GO_VERSION=`cat GO_VERSION`
DOCKER_BUILD_IMAGE=gotify/build
DOCKER_WORKDIR=/proj
DOCKER_RUN=docker run --rm -e LD_FLAGS="$$LD_FLAGS" -v "$$PWD/.:${DOCKER_WORKDIR}" -v "`go env GOPATH`/pkg/mod/.:/go/pkg/mod:ro" -w ${DOCKER_WORKDIR}
DOCKER_GO_BUILD=go build -mod=readonly -a -installsuffix cgo -ldflags "$$LD_FLAGS"
DOCKER_TEST_LEVEL ?= 0 # Optionally run a test during docker build
NODE_OPTIONS=$(shell if node --help | grep -q -- "--openssl-legacy-provider"; then echo --openssl-legacy-provider; fi)

test: test-coverage test-js
check: check-go check-swagger check-js
check-ci: check-swagger check-js

require-version:
	if [ -n ${VERSION} ] && [[ $$VERSION == "v"* ]]; then echo "The version may not start with v" && exit 1; fi
	if [ -z ${VERSION} ]; then echo "Need to set VERSION" && exit 1; fi;

test-coverage:
	go test --race -coverprofile=coverage.txt -covermode=atomic ./...

format:
	goimports -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

test-js:
	go build -ldflags="-s -w -X main.Mode=prod" -o removeme/gotify app.go
	(cd ui && CI=true GOTIFY_EXE=../removeme/gotify yarn test)
	rm -rf removeme

check-go:
	golangci-lint run

check-js:
	(cd ui && yarn lint)
	(cd ui && yarn testformat)

download-tools:
	go install github.com/go-swagger/go-swagger/cmd/swagger@v0.31.0

update-swagger:
	swagger generate spec --scan-models -o docs/spec.json
	sed -i 's/"uint64"/"int64"/g' docs/spec.json

check-swagger: update-swagger
## add the docs to git, this changes line endings in git, otherwise this does not work on windows
	git add docs
	if [ -n "$(shell git status --porcelain | grep docs)" ]; then \
        echo Swagger Spec is not up-to-date; \
        exit 1; \
    fi

extract-licenses:
	mkdir ${LICENSE_DIR} || true
	for LICENSE in $(shell find vendor/* -name LICENSE); do \
		DIR=`echo $$LICENSE | tr "/" _ | sed -e 's/vendor_//; s/_LICENSE//'` ; \
        cp $$LICENSE ${LICENSE_DIR}$$DIR ; \
    done

package-zip: extract-licenses
	for BUILD in $(shell find ${BUILD_DIR}/*); do \
       zip -j $$BUILD.zip $$BUILD ./LICENSE; \
       zip -ur $$BUILD.zip ${LICENSE_DIR}; \
    done

build-docker-multiarch: require-version
	docker buildx build --sbom=true --provenance=true \
		$(if $(DOCKER_BUILD_PUSH),--push) \
		-t gotify/server:latest \
		-t gotify/server:${VERSION} \
		-t gotify/server:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t gotify/server:$(shell echo $(VERSION) | cut -d '.' -f -1) \
	    -t ghcr.io/gotify/server:latest \
		-t ghcr.io/gotify/server:${VERSION} \
		-t ghcr.io/gotify/server:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t ghcr.io/gotify/server:$(shell echo $(VERSION) | cut -d '.' -f -1) \
		-t gotify/server-arm64:latest \
		-t gotify/server-arm64:${VERSION} \
		-t gotify/server-arm64:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t gotify/server-arm64:$(shell echo $(VERSION) | cut -d '.' -f -1) \
		-t ghcr.io/gotify/server-arm64:latest \
		-t ghcr.io/gotify/server-arm64:${VERSION} \
		-t ghcr.io/gotify/server-arm64:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t ghcr.io/gotify/server-arm64:$(shell echo $(VERSION) | cut -d '.' -f -1) \
		-t gotify/server-arm7:latest \
		-t gotify/server-arm7:${VERSION} \
		-t gotify/server-arm7:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t gotify/server-arm7:$(shell echo $(VERSION) | cut -d '.' -f -1) \
		-t ghcr.io/gotify/server-arm7:latest \
		-t ghcr.io/gotify/server-arm7:${VERSION} \
		-t ghcr.io/gotify/server-arm7:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t ghcr.io/gotify/server-arm7:$(shell echo $(VERSION) | cut -d '.' -f -1) \
		-t gotify/server-riscv64:latest \
		-t gotify/server-riscv64:${VERSION} \
		-t gotify/server-riscv64:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t gotify/server-riscv64:$(shell echo $(VERSION) | cut -d '.' -f -1) \
		-t ghcr.io/gotify/server-riscv64:latest \
		-t ghcr.io/gotify/server-riscv64:${VERSION} \
		-t ghcr.io/gotify/server-riscv64:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t ghcr.io/gotify/server-riscv64:$(shell echo $(VERSION) | cut -d '.' -f -1) \
		--build-arg RUN_TESTS=$(DOCKER_TEST_LEVEL) \
		--build-arg GO_VERSION=$(shell cat GO_VERSION) \
		--build-arg LD_FLAGS="$$LD_FLAGS" \
		--platform linux/amd64,linux/arm64,linux/386,linux/arm/v7,linux/riscv64 \
		-f docker/Dockerfile .

build-docker: build-docker-multiarch

_build_within_docker: OUTPUT = gotify-app
_build_within_docker:
	${DOCKER_GO_BUILD} -o ${OUTPUT}

build-js:
	(cd ui && NODE_OPTIONS="${NODE_OPTIONS}" yarn build)

build-linux-amd64:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-amd64 make _build_within_docker OUTPUT=${BUILD_DIR}/gotify-linux-amd64

build-linux-386:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-386 make _build_within_docker OUTPUT=${BUILD_DIR}/gotify-linux-386

build-linux-arm-7:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-arm-7 make _build_within_docker OUTPUT=${BUILD_DIR}/gotify-linux-arm-7

build-linux-arm64:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-arm64 make _build_within_docker OUTPUT=${BUILD_DIR}/gotify-linux-arm64

build-linux-riscv64:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-riscv64 make _build_within_docker OUTPUT=${BUILD_DIR}/gotify-linux-riscv64

build-windows-amd64:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-windows-amd64 make _build_within_docker OUTPUT=${BUILD_DIR}/gotify-windows-amd64.exe

build-windows-386:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-windows-386 make _build_within_docker OUTPUT=${BUILD_DIR}/gotify-windows-386.exe

build: build-linux-arm-7 build-linux-amd64 build-linux-386 build-linux-arm64 build-linux-riscv64 build-windows-amd64 build-windows-386

.PHONY: test-race test-coverage test check-go check-js verify-swagger check download-tools update-swagger package-zip build-docker build-js build
