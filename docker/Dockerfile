FROM frolvlad/alpine-glibc:glibc-2.28
WORKDIR /app
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ADD gotify-app /app/
ENTRYPOINT ["./gotify-app"]
