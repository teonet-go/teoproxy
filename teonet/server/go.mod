module github.com/teonet-go/teoproxy/teonet/server

go 1.21.5

// replace github.com/teonet-go/teoproxy/ws/client => ../../ws/client/
// replace github.com/teonet-go/teoproxy/ws/command => ../../ws/command/
// replace github.com/teonet-go/teoproxy/ws/server => ../../ws/server/

require (
	github.com/gorilla/websocket v1.5.1
	github.com/teonet-go/teomon v0.5.14
	github.com/teonet-go/teonet v0.6.6
	github.com/teonet-go/teoproxy/ws/command v0.0.0-00010101000000-000000000000
	github.com/teonet-go/teoproxy/ws/server v0.0.0-00010101000000-000000000000
)

require (
	github.com/denisbrodbeck/machineid v1.0.1 // indirect
	github.com/google/uuid v1.4.0 // indirect
	github.com/kirill-scherba/bslice v0.0.2 // indirect
	github.com/kirill-scherba/stable v0.0.8 // indirect
	github.com/teonet-go/tru v0.0.18 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
)
