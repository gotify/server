FROM frolvlad/alpine-glibc:glibc-2.28

ENV GID 999
ENV UID 999

WORKDIR /app

RUN addgroup -g ${GID} -S gotify && \
    adduser -u ${UID} -S gotify -G gotify && \
    apk add --no-cache ca-certificates tzdata

COPY --chown=gotify:gotify gotify-app /app/

EXPOSE 8080

ENTRYPOINT ["./gotify-app"]

USER gotify