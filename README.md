# Gotify Server
[![Build Status][badge-travis]][travis] [![codecov][badge-codecov]][codecov] [![Go Report Card][badge-go-report]][go-report] [![Swagger Valid][badge-swagger]][swagger] [![Api Docs][badge-api-docs]][api-docs] [![latest release version][badge-release]][release]

<img align="right" src="logo.png" />

   * [Motivation](#motivation)
   * [Features](#features)
   * [Installation](#installation)
   * [Configuration](#configuration)
   * [Setup Dev Environment](#setup-dev-environment)
   * [Add Message Examples](#add-message-examples)
   * [Building](#building)
   * [Tests](#tests)
   * [Versioning](#versioning)
   * [License](#license)

## Motivation
We wanted a simple server for sending and receiving messages (in real time per web socket). For this, not many open source projects existed and most of the existing ones were abandoned. Also, a requirement was that it can be self-hosted. We know there are many free and commercial push services out there.

## Features
* REST-API for
  * sending messages
  * receiving messages per websocket
  * user management
  * client/device & application management
* [REST-API Documentation][api-docs] (also available at `/docs`)
* Web-UI
<img alt="Gotify UI screenshot" src="ui.png" />

* Android-App -> [gotify/android](https://github.com/gotify/android)

[<img src="https://play.google.com/intl/en_gb/badges/images/generic/en_badge_web_generic.png" alt="Get it on Google Play" width="150" />][playstore]
[<img src="https://f-droid.org/badge/get-it-on.png" alt="Get it on F-Droid" width="150"/>][fdroid]

Google Play and the Google Play logo are trademarks of Google LLC.

## Installation

### Docker
The docker image is available on docker hub at [gotify/server][docker-normal].

``` bash
$ docker run -p 80:80 gotify/server
```
Also there is a specific docker image for arm-7 processors (raspberry pi), named [gotify/server-arm7][docker-arm7].
``` bash
$ docker run -p 80:80 gotify/server-arm7
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
      cache: data/certs # the directory of the cache from letsencrypt
      hosts: # the hosts for which letsencrypt should request certificates
      - mydomain.tld
      - myotherdomain.tld
database: # for database see (configure database section)
  dialect: sqlite3
  connection: data/gotify.db
defaultuser: # on database creation, gotify creates an admin user
  name: admin # the username of the default user
  pass: admin # the password of the default user
passstrength: 10 # the bcrypt password strength (higher = better but also slower)
uploadedimagesdir: data/images # the directory for storing uploaded images
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
GOTIFY_UPLOADEDIMAGESDIR=images
```

### Add Message Examples

You can obtain an application-token from the apps tab inside the UI or using the REST-API (`GET /application`)

NOTE: Assuming Gotify is running on `http://localhost:8008`.

**curl**
```bash
  curl -X POST "http://localhost:8008/message?token=<token-from-application>" -F "title=My Title" -F "message=This is my message"
```

**python**

```python
import requests #pip install requests
resp = requests.post('http://localhost:8008/message?token=<token-from-application>', json={
    "message": "Well hello there.",
    "priority": 2,
    "title": "This is my title"
})
```

**golang**

```go
package main

import (
        "net/http"
        "net/url"
)

func main() {
    http.PostForm("http://localhost:8008/message?<token-from-application>", url.Values{"message": {"My Message"}, "title": {"My Title"}})
}
```


### Database
| Dialect   | Connection                                                           |
| :-------: | :------------------------------------------------------------------: |
| sqlite3   | `path/to/database.db`                                                |
| mysql     | `gotify:secret@/gotifydb?charset=utf8&parseTime=True&loc=Local `     |
| postgres  | `host=localhost port=3306 user=gotify dbname=gotify password=secret` |

## Setup Dev Environment

### Setup Server
Download go dependencies with [golang/dep](https://github.com/golang/dep).
```
$ dep ensure
```
Run golang server.
```
$ go run app.go
```

### Setup UI
*Commands must be executed inside the ui directory.*

Download dependencies with [npm](https://github.com/npm/npm).
``` bash
$ npm install
```
Star the UI development server.
``` bash
$ npm start
```
Open `http://localhost:3000` inside your favorite browser.

The UI requires a Gotify server running on `localhost:80` this can be adjusted inside the [ui/src/index.tsx](ui/src/index.tsx).

## Building

### Build Server
``` bash
$ go build app.go
```

### Build UI
``` bash
$ npm run build
```

### Cross-Platform
The project has a CGO reference (because of sqlite3), therefore a GCO cross compiler is needed for compiling for other platforms. See [.travis.yml](.travis.yml) on how we do that.

## Tests
The tests can be executed with:
``` bash
$ make test
# or
$ go test ./...
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
 [docker-normal]: https://hub.docker.com/r/gotify/server/
 [docker-arm7]: https://hub.docker.com/r/gotify/server-arm7/
 [playstore]: https://play.google.com/store/apps/details?id=com.github.gotify
 [fdroid]: https://f-droid.org/de/packages/com.github.gotify/
