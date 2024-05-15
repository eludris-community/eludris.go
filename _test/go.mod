module github.com/eludris-community/eludris.go/_test

go 1.21.3

replace github.com/eludris-community/eludris.go/v2 => ../

require (
	github.com/apex/log v1.9.0
	github.com/eludris-community/eludris.go/v2 v2.0.0-20240426121112-70fcb0ccb363
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/eludris-community/eludris-api-types.go/v2 v2.0.3-0.20240515230257-8861b88c0ad0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sasha-s/go-csync v0.0.0-20240107134140-fcbab37b09ad // indirect
	golang.org/x/net v0.25.0 // indirect
)
