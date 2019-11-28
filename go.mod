module github.com/gotify/server

require (
	github.com/Southclaws/configor v1.0.0 // indirect
	github.com/fortytw2/leaktest v1.3.0
	github.com/gin-contrib/gzip v0.0.1
	github.com/gin-gonic/gin v1.5.0
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/gobuffalo/packr v1.22.0
	github.com/gorilla/websocket v1.4.0
	github.com/gotify/configor v1.0.2-0.20190112111140-7d9c7c7e6233
	github.com/gotify/location v0.0.0-20170722210143-03bc4ad20437
	github.com/gotify/plugin-api v1.0.0
	github.com/h2non/filetype v1.0.10
	github.com/jinzhu/gorm v1.9.11
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/robfig/cron v0.0.0-20180505203441-b41be1df6967
	github.com/stretchr/testify v1.4.0
	golang.org/x/crypto v0.0.0-20191128160524-b544559bb6d1
	golang.org/x/sys v0.0.0-20191128015809-6d18c012aee9 // indirect
	gopkg.in/go-playground/validator.v9 v9.30.2
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.12.0

go 1.13
