ARG BUILDKIT_SBOM_SCAN_CONTEXT=true
# Suppress warning about invalid variable expansion
ARG GO_VERSION=PLEASE_PROVIDE_GO_VERSION
ARG DEBIAN=sid-slim

# Hack to normalize platform to match the chosed build image
# Get the gotify/build image tag
ARG __TARGETPLATFORM_DASHES=${TARGETPLATFORM/\//-}
ARG __TARGETPLATFORM_GO_NOTATION=${__TARGETPLATFORM_DASHES/arm\/v7/arm-7}

# --- JS Builder ---

FROM --platform=${BUILDPLATFORM} node:23 AS js-builder

ARG BUILD_JS=0

COPY ./Makefile /src/gotify/Makefile
COPY ./ui /src/gotify/ui

RUN if [ "$BUILD_JS" = "1" ]; then \
    (cd /src/gotify/ui && yarn install) && \
    (cd /src/gotify && make build-js) \
    else \
    mkdir -p /src/gotify/ui/build; \
    fi

# --- Go Builder ---

FROM --platform=${BUILDPLATFORM} gotify/build:${GO_VERSION}-${__TARGETPLATFORM_GO_NOTATION} AS builder

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG BUILD_JS=0
ARG RUN_TESTS=0 # 0=never, 1=native only
ARG LD_FLAGS=""
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -yq --no-install-recommends \
    ca-certificates \
    git

COPY . /src/gotify
COPY --from=js-builder /src/gotify/ui/build /ui-build

RUN if [ "$BUILD_JS" = "1" ]; then \
    cp -r --update /ui-build /src/gotify/ui/build; \
    fi

RUN cd /src/gotify && \
    if [ "$RUN_TESTS" = "1" ] && [ "$BUILDPLATFORM" = "$TARGETPLATFORM" ]; then \
    go test -v ./...; \
    fi && \
    LD_FLAGS=${LD_FLAGS} make OUTPUT=/target/app/gotify-app _build_within_docker

FROM debian:${DEBIAN}

# Build-time configurables
ARG GOTIFY_SERVER_EXPOSE=80
ENV GOTIFY_SERVER_PORT=$GOTIFY_SERVER_EXPOSE

WORKDIR /app

RUN export DEBIAN_FRONTEND=noninteractive && apt-get update && apt-get install -yq --no-install-recommends \
    tzdata \
    curl \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s CMD curl --fail http://localhost:$GOTIFY_SERVER_PORT/health || exit 1
EXPOSE $GOTIFY_SERVER_EXPOSE

COPY --from=builder /target /

ENTRYPOINT ["./gotify-app"]
