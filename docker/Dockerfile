ARG DEBIAN=stable

# --- JS Builder ---

FROM debian:${DEBIAN} AS js-builder

ARG BUILD_JS=1
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

FROM debian:${DEBIAN} AS builder

ARG BUILD_JS=1
ARG GO_BUILD_FLAGS=-mod=readonly -a -installsuffix cgo
ENV DEBIAN_FRONTEND=noninteractive

ADD https://raw.githubusercontent.com/travis-ci/gimme/master/gimme /usr/local/bin/gimme 

RUN chmod +x /usr/local/bin/gimme # compatiblity

RUN apt-get update && apt-get install -yq --no-install-recommends \
    build-essential \
    ca-certificates \
    curl \
    git

COPY . /src/gotify
COPY --from=js-builder /target /ui-build

RUN if [ "$BUILD_JS" = "1" ]; then \
    cp -r --update /ui-build /src/gotify/ui/build; \
    fi

RUN cd /src/gotify && \
    read -r GIMME_GO_VERSION < GO_VERSION && \
    echo "Using Go version: $GIMME_GO_VERSION" && \
    gimme ${GIMME_GO_VERSION} && \
    mkdir -p /target/app && \
    bash -c ". ~/.gimme/envs/go${GIMME_GO_VERSION}.env && go build ${GO_BUILD_FLAGS} -o /target/app/gotify-app"

FROM debian:${DEBIAN}-slim

# build-time configurable port
ARG GOTIFY_SERVER_EXPOSE=80
ENV GOTIFY_SERVER_PORT=$GOTIFY_SERVER_EXPOSE

WORKDIR /app

RUN export DEBIAN_FRONTEND=noninteractive && apt-get update && apt-get install -yq --no-install-recommends \
    tzdata \
    curl \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s CMD curl --fail http://localhost:$GOTIFY_SERVER_PORT/health || exit 1
EXPOSE $GOTIFY_SERVER_EXPOSE

COPY --from=builder /target /

ENTRYPOINT ["./gotify-app"]
