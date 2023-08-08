module github.com/eludris-community/eludris.go/v2

go 1.19

replace github.com/eludris-community/eludris.go/v2 => ./

require (
	github.com/apex/log v1.9.0
	github.com/eludris-community/eludris-api-types.go v0.0.0-20230325172402-dd501e701d08
	github.com/gorilla/websocket v1.5.0
	github.com/mitchellh/mapstructure v1.5.0
)

require github.com/pkg/errors v0.8.1 // indirect
