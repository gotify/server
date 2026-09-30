package docs

import "github.com/gin-gonic/gin"

var ui = `
<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.29.1/swagger-ui.css" integrity="sha512-eZpfl9qlKnbDvlJ2brfdx3nhlP1FMsA23w65motxKdYsUcfMcdO2bcLPr7mXhvyzmDZwuzYCJKrl/sEo1ditVQ==" crossorigin="anonymous" referrerpolicy="no-referrer" >
    <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }
      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }
      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>

    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.29.1/swagger-ui-bundle.js" integrity="sha512-iDdi6WwSTimAFh1NhPyFZpHJWr16m/PupHztElqQL+gQoQDavUATsP9hcvgs9Yci+EUA//WvZYZsCTDhRDKc3g==" crossorigin="anonymous" referrerpolicy="no-referrer"> </script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.29.1/swagger-ui-standalone-preset.js" integrity="sha512-oGq+mP+3XUyYDm9vsprsd6YdhOokdkcLNUZnmEE6HJzG0F9fIaB8kgGZsjdfTWN8w/1f2kN4a/0xkla0hsetcQ==" crossorigin="anonymous" referrerpolicy="no-referrer"> </script>
    <script>
    function getBaseURL() {
      var path = window.location.pathname
      path = path.substr(0, path.lastIndexOf('/')+1)
      return window.location.host + path;
    }
    window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
        url: "swagger?base="+encodeURIComponent(getBaseURL()),
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      })
      // End Swagger UI call region
      window.ui = ui
    }
  </script>
  </body>
</html>
`

// UI serves the swagger ui.
func UI(ctx *gin.Context) {
	ctx.Writer.WriteString(ui)
}
