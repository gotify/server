LICENSE_DIR=./licenses/
BUILD_DIR=./build
DOCKER_DIR=./docker/
SHELL := /bin/bash
GO_VERSION=`cat GO_VERSION`
DOCKER_BUILD_IMAGE=gotify/build
DOCKER_WORKDIR=/proj
DOCKER_RUN=docker run --rm -v "$$PWD/.:${DOCKER_WORKDIR}" -v "`go env GOPATH`/pkg/mod/.:/go/pkg/mod:ro" -w ${DOCKER_WORKDIR}
DOCKER_GO_BUILD=go build -mod=readonly -a -installsuffix cgo -ldflags "$$LD_FLAGS"

test: test-coverage test-race test-js
check: check-go check-swagger check-js
check-ci: check-swagger check-js

require-version:
	if [ -n ${VERSION} ] && [[ $$VERSION == "v"* ]]; then echo "The version may not start with v" && exit 1; fi
	if [ -z ${VERSION} ]; then echo "Need to set VERSION" && exit 1; fi;

test-race:
	go test -race ./...

test-coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...

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
	go install github.com/go-swagger/go-swagger/cmd/swagger@v0.26.1

embed-static:
	go run hack/packr/packr.go

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

build-docker-amd64: require-version
	cp ${BUILD_DIR}/gotify-linux-amd64 ./docker/gotify-app
	cd ${DOCKER_DIR} && \
		docker build \
		-t gotify/server:latest \
		-t gotify/server:${VERSION} \
		-t gotify/server:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t gotify/server:$(shell echo $(VERSION) | cut -d '.' -f -1) \
		-t ghcr.io/gotify/server:latest \
		-t ghcr.io/gotify/server:${VERSION} \
		-t ghcr.io/gotify/server:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t ghcr.io/gotify/server:$(shell echo $(VERSION) | cut -d '.' -f -1) .
	rm ${DOCKER_DIR}gotify-app

build-docker-arm-7: require-version
	cp ${BUILD_DIR}/gotify-linux-arm-7 ./docker/gotify-app
	cd ${DOCKER_DIR} && \
		docker build -f Dockerfile.armv7 \
		-t gotify/server-arm7:latest \
		-t gotify/server-arm7:${VERSION} \
		-t gotify/server-arm7:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t gotify/server-arm7:$(shell echo $(VERSION) | cut -d '.' -f -1) \
		-t ghcr.io/gotify/server-arm7:latest \
		-t ghcr.io/gotify/server-arm7:${VERSION} \
		-t ghcr.io/gotify/server-arm7:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t ghcr.io/gotify/server-arm7:$(shell echo $(VERSION) | cut -d '.' -f -1) .
	rm ${DOCKER_DIR}gotify-app

build-docker-arm64: require-version
	cp ${BUILD_DIR}/gotify-linux-arm64 ./docker/gotify-app
	cd ${DOCKER_DIR} && \
		docker build -f Dockerfile.arm64 \
		-t gotify/server-arm64:latest \
		-t gotify/server-arm64:${VERSION} \
		-t gotify/server-arm64:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t gotify/server-arm64:$(shell echo $(VERSION) | cut -d '.' -f -1) \
		-t ghcr.io/gotify/server-arm64:latest \
		-t ghcr.io/gotify/server-arm64:${VERSION} \
		-t ghcr.io/gotify/server-arm64:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t ghcr.io/gotify/server-arm64:$(shell echo $(VERSION) | cut -d '.' -f -1) .
	rm ${DOCKER_DIR}gotify-app

build-docker-riscv64: require-version
	cp ${BUILD_DIR}/gotify-linux-riscv64 ./docker/gotify-app
	cd ${DOCKER_DIR} && \
		docker build -f Dockerfile.riscv64 \
		-t gotify/server-riscv64:latest \
		-t gotify/server-riscv64:${VERSION} \
		-t gotify/server-riscv64:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t gotify/server-riscv64:$(shell echo $(VERSION) | cut -d '.' -f -1) \
		-t ghcr.io/gotify/server-riscv64:latest \
		-t ghcr.io/gotify/server-riscv64:${VERSION} \
		-t ghcr.io/gotify/server-riscv64:$(shell echo $(VERSION) | cut -d '.' -f -2) \
		-t ghcr.io/gotify/server-riscv64:$(shell echo $(VERSION) | cut -d '.' -f -1) .
	rm ${DOCKER_DIR}gotify-app

build-docker: build-docker-amd64 build-docker-arm-7 build-docker-arm64 build-docker-riscv64

build-js:
	(cd ui && yarn build)

build-linux-amd64:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-amd64 ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/gotify-linux-amd64 ${DOCKER_WORKDIR}

build-linux-386:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-386 ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/gotify-linux-386 ${DOCKER_WORKDIR}

build-linux-arm-7:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-arm-7 ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/gotify-linux-arm-7 ${DOCKER_WORKDIR}

build-linux-arm64:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-arm64 ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/gotify-linux-arm64 ${DOCKER_WORKDIR}

build-linux-riscv64:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-linux-riscv64 ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/gotify-linux-riscv64 ${DOCKER_WORKDIR}

build-windows-amd64:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-windows-amd64 ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/gotify-windows-amd64.exe ${DOCKER_WORKDIR}

build-windows-386:
	${DOCKER_RUN} ${DOCKER_BUILD_IMAGE}:$(GO_VERSION)-windows-386 ${DOCKER_GO_BUILD} -o ${BUILD_DIR}/gotify-windows-386.exe ${DOCKER_WORKDIR}

build: build-linux-arm-7 build-linux-amd64 build-linux-386 build-linux-arm64 build-linux-riscv64 build-windows-amd64 build-windows-386

.PHONY: test-race test-coverage test check-go check-js verify-swagger check download-tools update-swagger package-zip build-docker build-js build
