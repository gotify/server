# Gotify Server
[![Build Status][badge-travis]][travis] [![codecov][badge-codecov]][codecov] [![Go Report Card][badge-go-report]][go-report] [![Swagger Valid][badge-swagger]][swagger] [![Api Docs][badge-api-docs]][api-docs] [![latest release version][badge-release]][release]

   * [Motivation](#motivation)
   * [Features](#features)
   * [Installation](#installation)
     * [Docker](#docker)
     * [Binary](#binary)
   * [Configuration](#configuration)
      * [File](#file)
      * [Environment](#environment)
      * [Database](#database)
   * [Building](#building)
      * [Cross-Platform](#cross-platform)
   * [Tests](#tests)
   * [Versioning](#versioning)
   * [License](#license)

## Motivation
We wanted a simple server for sending and receiving messages (in real time per websocket). For this, not many open source projects existed and most of the existing ones were abandoned. Also, a requirement was that it can be self-hosted. We know there are many free and commercial push services out there.

## Features
* API (see [api docs][api-docs]) for
  * sending messages
  * receiving messages per websocket
  * user management
  * client/device & application management
* *[In Progress]* Web-UI
* *[In Progress]* Android-App -> [gotify/android](https://github.com/gotify/android)

## Installation

### Docker
The docker image is available on docker hub at [gotify/server](https://hub.docker.com/r/gotify/server/).

``` bash
docker run -p 80:80 gotify/server
```
Also there is a specific docker image for arm-7 processors (raspberry pi), named [gotify/server-arm7](https://hub.docker.com/r/gotify/server-arm7/).
``` bash
docker run -p 80:80 gotify/server-arm7
```

### Binary
Visit the [releases page](https://github.com/gotify/server/releases) and download the zip for your OS.

## Configuration
### File
``` yml
server:
  port: 80 # the port for the http server
  ssl:
    enabled: false # if https should be enabled
    redirecttohttps: true # redirect to https if site is accessed by http
    port: 443 # the https port
    certfile: # the cert file (leave empty when using letsencrypt)
    certkey: # the cert key (leave empty when using letsencrypt)
    letsencrypt:
      enabled: false # if the certificate should be requested from letsencrypt
      accepttos: false # if you accept the tos from letsencrypt
      cache: certs # the directory of the cache from letsencrypt
      hosts: # the hosts for which letsencrypt should request certificates
      - mydomain.tld
      - myotherdomain.tld
database: # for database see (configure database section)
  dialect: sqlite3
  connection: gotify.db
defaultuser: # on database creation, gotify creates an admin user
  name: admin # the username of the default user
  pass: admin # the password of the default user
passstrength: 10 # the bcrypt password strength (higher = better but also slower)
```

### Environment
``` bash
GOTIFY_SERVER_PORT=80
GOTIFY_SERVER_SSL_ENABLED=false
GOTIFY_SERVER_SSL_REDIRECTTOHTTPS=true
GOTIFY_SERVER_SSL_PORT=443
GOTIFY_SERVER_SSL_CERTFILE=
GOTIFY_SERVER_SSL_CERTKEY=
GOTIFY_SERVER_SSL_LETSENCRYPT_ENABLED=false
GOTIFY_SERVER_SSL_LETSENCRYPT_ACCEPTTOS=false
GOTIFY_SERVER_SSL_LETSENCRYPT_CACHE=certs
# lists are a little weird but do-able (:
GOTIFY_SERVER_SSL_LETSENCRYPT_HOSTS=- mydomain.tld\n- myotherdomain.tld
GOTIFY_DATABASE_DIALECT=sqlite3
GOTIFY_DATABASE_CONNECTION=gotify.db
GOTIFY_DEFAULTUSER_NAME=admin
GOTIFY_DEFAULTUSER_PASS=admin
GOTIFY_PASSSTRENGTH=10
```

### Database
| Dialect   | Connection                                                           |
| :-------: | :------------------------------------------------------------------: |
| sqlite3   | `path/to/database.db`                                                |
| mysql     | `gotify:secret@/gotifydb?charset=utf8&parseTime=True&loc=Local `     |
| postgres  | `host=localhost port=3306 user=gotify dbname=gotify password=secret` |

## Building

The app can be built with the default golang build command.
``` bash
go build app.go
```

### Cross-Platform
The project has a CGO reference (because of sqlite3), therefore a GCO cross compiler is needed for compiling for other platforms. We use [karalabe/xgo](https://github.com/karalabe/xgo) for this, xgo is a bundle of docker containers for building go apps.
``` bash
VERSION=mybuild1 make build-binary
```

## Tests
The tests can be executed with:
``` bash
make test
# or
go test ./...
```

## Versioning
We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/gotify/server/tags).

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

 [badge-api-docs]: https://img.shields.io/badge/api-docs-blue.svg
 [badge-swagger]: https://img.shields.io/swagger/valid/2.0/https/raw.githubusercontent.com/gotify/server/master/docs/spec.json.svg
 [badge-go-report]: https://goreportcard.com/badge/github.com/gotify/server
 [badge-codecov]: https://codecov.io/gh/gotify/server/branch/master/graph/badge.svg
 [badge-travis]: https://travis-ci.org/gotify/server.svg?branch=master
 [badge-release]: https://img.shields.io/github/release/gotify/server.svg
 [release]: https://github.com/gotify/server/releases/latest
 [travis]: https://travis-ci.org/gotify/server
 [codecov]: https://codecov.io/gh/gotify/server
 [go-report]: https://goreportcard.com/report/github.com/gotify/server
 [swagger]: https://github.com/gotify/server/blob/master/docs/spec.json
 [api-docs]: https://gotify.github.io/api-docs/