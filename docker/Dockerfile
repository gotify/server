FROM frolvlad/alpine-glibc:glibc-2.28
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
ADD gotify-app /app/
EXPOSE 80
ENTRYPOINT ["./gotify-app"]
