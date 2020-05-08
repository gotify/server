module github.com/gotify/server/v2

require (
	github.com/Southclaws/configor v1.0.0 // indirect
	github.com/fortytw2/leaktest v1.3.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/gzip v0.0.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.2.0
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/gobuffalo/envy v1.9.0 // indirect
	github.com/gobuffalo/logger v1.0.3 // indirect
	github.com/gobuffalo/packd v1.0.0 // indirect
	github.com/gobuffalo/packr/v2 v2.7.1
	github.com/golang/protobuf v1.4.1 // indirect
	github.com/gorilla/websocket v1.4.0
	github.com/gotify/configor v1.0.2-0.20190112111140-7d9c7c7e6233
	github.com/gotify/location v0.0.0-20170722210143-03bc4ad20437
	github.com/gotify/plugin-api v1.0.0
	github.com/h2non/filetype v1.0.10
	github.com/jinzhu/gorm v1.9.11
	github.com/lib/pq v1.5.2 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/robfig/cron v0.0.0-20180505203441-b41be1df6967
	github.com/rogpeppe/go-internal v1.5.2 // indirect
	github.com/stretchr/testify v1.5.1
	golang.org/x/crypto v0.0.0-20200429183012-4b2356b1ed79
	golang.org/x/sys v0.0.0-20200501145240-bc7a7d42d5c3 // indirect
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.12.0

go 1.13
