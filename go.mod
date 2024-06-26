module github.com/eludris-community/eludris.go/v2

go 1.19

replace github.com/eludris-community/eludris.go/v2 => ./

require (
	github.com/apex/log v1.9.0
	github.com/eludris-community/eludris-api-types.go/v2 v2.0.2
	github.com/gorilla/websocket v1.5.0
	github.com/sasha-s/go-csync v0.0.0-20210812194225-61421b77c44b
)

require github.com/pkg/errors v0.8.1 // indirect
