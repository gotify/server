FROM frolvlad/alpine-glibc:glibc-2.31
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata curl
ADD gotify-app /app/
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s CMD curl --fail http://localhost:80/health || exit 1
EXPOSE 80
ENTRYPOINT ["./gotify-app"]
