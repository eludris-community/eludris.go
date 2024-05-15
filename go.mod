module github.com/eludris-community/eludris.go/v2

go 1.19

replace github.com/eludris-community/eludris.go/v2 => ./

require (
	github.com/apex/log v1.9.0
	github.com/eludris-community/eludris-api-types.go/v2 v2.0.3-0.20240515230257-8861b88c0ad0
	github.com/gorilla/websocket v1.5.1
	github.com/sasha-s/go-csync v0.0.0-20240107134140-fcbab37b09ad
)

require (
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.25.0 // indirect
)
