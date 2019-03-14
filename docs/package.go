// Package docs Gotify REST-API.
//
// This is the documentation of the Gotify REST-API.
//
// # Authentication
// In Gotify there are two token types:
// __clientToken__: a client is something that receives message and manages stuff like creating new tokens or delete messages. (f.ex this token should be used for an android app)
// __appToken__: an application is something that sends messages (f.ex. this token should be used for a shell script)
//
// The token can be either transmitted through a header named `X-Gotify-Key` or a query parameter named `token`.
// There is also the possibility to authenticate through basic auth, this should only be used for creating a clientToken.
//
// \---
//
// Found a bug or have some questions? [Create an issue on GitHub](https://github.com/gotify/server/issues)
//
//     Schemes: http, https
//     Host: localhost
//     Version: 2.0.1
//     License: MIT https://github.com/gotify/server/blob/master/LICENSE
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     SecurityDefinitions:
//        appTokenQuery:
//           type: apiKey
//           name: token
//           in: query
//        clientTokenQuery:
//           type: apiKey
//           name: token
//           in: query
// 	      appTokenHeader:
//           type: apiKey
//           name: X-Gotify-Key
//           in: header
// 	      clientTokenHeader:
//           type: apiKey
//           name: X-Gotify-Key
//           in: header
//        basicAuth:
//           type: basic
//
// swagger:meta
package docs
