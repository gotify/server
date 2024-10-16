ARG BUILDKIT_SBOM_SCAN_CONTEXT=true
ARG DEBIAN_STABLE_ARCH=stable-slim
ARG DEBIAN_EXPERIMENTAL_ARCH=sid-slim
# suppress warnung about invalid variable expansion
ARG GO_VERSION=PLEASE_PROVIDE_GO_VERSION

# Hack to normalize platform to match the chosed build image
# Get the gotify/build image tag
ARG __TARGETPLATFORM_DASHES=${TARGETPLATFORM/\//-}
ARG __TARGETPLATFORM_GO_NOTATION=${__TARGETPLATFORM_DASHES/arm\/v7/arm-7}
# Set the unstable tag if the target platform is not stable
ARG __TARGETPLATFORM_IS_STABLE=${TARGETPLATFORM/linux\/riscv64/}
ARG __DEBIAN_HAS_STABLE_TAG=${__TARGETPLATFORM_IS_STABLE/?*/${DEBIAN_STABLE_ARCH}} # if the last variable is not empty, set it to the debian stable tag
ARG DEBIAN=${__DEBIAN_HAS_STABLE_TAG:-${DEBIAN_EXPERIMENTAL_ARCH}} # else set it to the debian experimental tag

# --- JS Builder ---

FROM --platform=${BUILDPLATFORM} debian:${DEBIAN} AS js-builder

ARG BUILD_JS=0
ARG NODE_OPTIONS
ENV DEBIAN_FRONTEND=noninteractive

COPY . /src/gotify

RUN if [ "$BUILD_JS" = "1" ]; then \
    apt-get update && apt-get install -yq --no-install-recommends \
    curl \
    git \
    nodejs \
    npm && \
    \
    NODE_OPTIONS_DEFAULT=$(if node --help | grep -q -- "--openssl-legacy-provider"; then echo --openssl-legacy-provider; fi) && \
    export NODE_OPTIONS=${NODE_OPTIONS:-$NODE_OPTIONS_DEFAULT} && \
    echo "Using NODE_OPTIONS=$NODE_OPTIONS" && \
    cd /src/gotify/ui && \
    \
    npm install -g yarn && \
    \
    yarn install && \
    yarn build && \
    \
    cp -r /src/gotify/ui/build /target; \
    \
    else \
    mkdir -p /target; \
    fi

# --- Go Builder ---

FROM --platform=${BUILDPLATFORM} gotify/build:${GO_VERSION}-${__TARGETPLATFORM_GO_NOTATION} AS builder

ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG BUILD_JS=0
ARG RUN_TESTS=0 # 0=never, 1=native only
ARG GO_BUILD_FLAGS=-mod=readonly -a -installsuffix cgo
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -yq --no-install-recommends \
    ca-certificates \
    git

COPY . /src/gotify
COPY --from=js-builder /target /ui-build

RUN if [ "$BUILD_JS" = "1" ]; then \
    cp -r --update /ui-build /src/gotify/ui/build; \
    fi

RUN cd /src/gotify && \
    mkdir -p /target/app && \
    if [ "$RUN_TESTS" = "1" ] && [ "$BUILDPLATFORM" = "$TARGETPLATFORM" ]; then \
    go test -v ./...; \
    fi && \
    go build ${GO_BUILD_FLAGS} -o /target/app/gotify-app

FROM debian:${DEBIAN}

# build-time configurables
ARG GOTIFY_SERVER_EXPOSE=80
ARG GOTIFY_UID=0
ENV GOTIFY_SERVER_PORT=$GOTIFY_SERVER_EXPOSE

WORKDIR /app

RUN export DEBIAN_FRONTEND=noninteractive && apt-get update && apt-get install -yq --no-install-recommends \
    tzdata \
    curl \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    if [ "$GOTIFY_UID" != "0" ]; then \
    useradd -u $GOTIFY_UID -r gotify && \
    chown gotify:gotify /app; \
    fi

USER $GOTIFY_UID

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s CMD curl --fail http://localhost:$GOTIFY_SERVER_PORT/health || exit 1
EXPOSE $GOTIFY_SERVER_EXPOSE

COPY --from=builder /target /

ENTRYPOINT ["./gotify-app"]
