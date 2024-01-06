# Teoproxy

Teonet proxy client server packages to connect golang wasm applications with [Teonet](https://github.com/teonet-go) peers.

[![GoDoc](https://godoc.org/github.com/teonet-go/teoproxy?status.svg)](https://godoc.org/github.com/teonet-go/teoproxy/)
[![Go Report Card](https://goreportcard.com/badge/github.com/teonet-go/teoproxy)](https://goreportcard.com/report/github.com/teonet-go/teoproxy)

Teoproxy provides a websocket client server packages that can be used to connect wasm application with it own web server which runs Teonet and connects to teonet peers used in wasm application.

<p align="center">
<img src="https://github.com/teonet-go/.github/blob/main/profile/microservices.jpg?raw=true" />
</p>

## Getting started

There is main example in the `cmd/teonet/fortune-gui` folder. The example shows how to connect to Teonet peers with wasm application, and send and receive messages from wasm application to Teonet "fortune" peer.

There is [complex Teonet example](https://github.com/teonet-go/.github/blob/main/profile/complex.md) which use mach teonet applications to get fortune messages from fortune teonet server. This teoproxy example do the same teonet function:

- connect to Teonet
- connect to Teonet "fortune" peer
- connect to "fortune" api
- request and receive "fortune" messages from Teonet server

To do this we use fyne-ios package to make simple gui application and teoproxy/teonet client and server packages to connect wasm application with Teonet peers.

Go to [fortune-gui](cmd/teonet/fortune-gui) (cmd/teonet/fortune-gui) folder, and run native example:

```bash
go run main.go
```

To create web server for this application, go to the [fortune-gui/serve](cmd/teonet/fortune-gui/serve) (cmd/teonet/fortune-gui/serve) and execute next commands:

```bash

# Install fyne executible (if not installed)
go install fyne.io/fyne/v2/cmd/fyne@latest

# Build web package (or you can use `go generate` command to build
# and run this web server)
fyne package -os wasm --sourceDir ../

# Run web server
go run .

```

By default web server runs on `localhost:8081` port. So you can open this url in your browser: [http://localhost:8081](http://localhost:8081) and see the fortune-gui application in your browser.

## License

[BSD](LICENSE)
