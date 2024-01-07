// Copyright 2023-2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm

// Teonet wasm client to use in wasm applications.
package client

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	ws "github.com/teonet-go/teoproxy/ws/client"
	"github.com/teonet-go/teoproxy/ws/command"
)

// Teonet represents a Teonet client instance. It contains:
// - ws: Websocket client
// - id: Packet id
type Teonet struct {
	ws *ws.WsClient // Websocket client
	id uint32       // Packet id
}

// New creates a new Teonet client instance. It initializes the websocket
// client and connects to the Teonet proxy server. The appShort string is
// used for logging. The onReconnected callback is invoked when the
// websocket reconnects after a disconnect. It returns a pointer to the
// Teonet client and an error.
func New(appShort string, onReconnected func()) (teo *Teonet, err error) {
	teo = new(Teonet)
	teo.ws = ws.NewWsClient(
		// Common reader. It process Id 0 command answers.
		func(message []byte) bool {
			cmd := command.NewEmpty()
			err = cmd.UnmarshalBinary(message)
			if err != nil {
				log.Println("Can't unmarshal teonet proxy server command, error:",
					err, string(message))
				return false
			}
			if cmd.Id == 0 {
				log.Println("Recv id", cmd.Id, cmd.Cmd.String())
				return true
			}

			return false
		},
	)
	err = teo.ws.Connect(onReconnected)
	return
}

// getNextID atomically increments the id field by 1 and returns
// the incremented value. This provides each packet sent via the Teonet
// client with a unique id.
func (teo *Teonet) getNextID() uint32 {
	return atomic.AddUint32(&teo.id, 1)
}

// Connect sends a Connect command to the Teonet proxy server
// to establish a connection. It marshals the command into a
// binary format and sends it via the websocket client.
// It returns any error from marshaling or sending the command.
func (teo *Teonet) Connect() (err error) {
	cmd := command.New(command.Connect, nil)
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	return
}

// Disconnect sends a Disconnect command to the Teonet proxy server
// to close the connection. It marshals the command into a
// binary format and sends it via the websocket client.
// It returns any error from marshaling or sending the command.
func (teo *Teonet) Disconnect() (err error) {
	cmd := command.New(command.Disconnect, nil)
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	return
}

// ConnectTo sends a ConnectTo command with the provided peer name to the Teonet
// proxy server to establish a connection to that peer. It marshals the command
// into a binary format and sends it via the websocket client.
// It returns any error from marshaling or sending the command.
func (teo *Teonet) ConnectTo(peer string) (err error) {
	cmd := command.New(command.ConnectTo, []byte(peer))
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	return
}

// NewAPIClient creates a new APIClient instance that can be used to make
// API calls to the peer specified in the peer parameter. It sends a
// NewApiClient command to the Teonet proxy server to establish the API
// connection, and returns a pointer to the new APIClient instance.
func (teo *Teonet) NewAPIClient(peer string) (cli *APIClient, err error) {
	cmd := command.New(command.NewApiClient, []byte(peer))
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	cli = &APIClient{teo: teo, addr: peer}
	return
}

// WaitFrom waits to receive a response with the given ID from the specified
// peer. It adds a reader callback to the websocket client that waits for a
// matching response, with a timeout. It returns the response data and any
// error. This allows waiting for async responses to requests sent to peers.
func (teo *Teonet) WaitFrom(peer string, id uint32) (data []byte, err error) {
	var readerId string
	type resultData struct {
		data []byte
		err  error
	}
	w := make(chan resultData, 1)
	readerId = teo.ws.AddReader(func(message []byte) bool {

		cmd := command.NewEmpty()
		err = cmd.UnmarshalBinary(message)
		if err != nil {
			log.Println("Can't unmarshal teonet proxy server command, error:",
				err, string(message))
			return false
		}
		if cmd.Id != id {
			return false
		}
		log.Println("Recv id", cmd.Id, cmd.Cmd.String())
		log.Println("Got Teonet proxy server command:", cmd.Cmd.String(),
			string(cmd.Data))

		go teo.ws.RemoveReader(readerId)
		w <- resultData{cmd.Data, nil}

		return true
	})

	// Get answer from server or timeout
	var answer resultData
	select {
	case answer = <-w:
	case <-time.After(5 * time.Second):
		answer = resultData{nil, fmt.Errorf("timeout")}
	}
	data, err = answer.data, answer.err

	return
}

// APIClient is a client for making API calls to peers over a Teonet connection.
// It contains the Teonet client instance and address of the peer to send
// requests to.
type APIClient struct {
	teo  *Teonet
	addr string
}

// Address returns the address of the peer this APIClient is configured to send
// requests to.
func (api *APIClient) Address() string {
	return api.addr
}

// SendTo sends an API command and data to the configured peer address.
// It returns a unique command ID and error. The apiCmd and apiData are
// joined into a single byte slice that is sent as the command data.
func (api *APIClient) SendTo(apiCmd string, apiData []byte) (id uint32, err error) {
	data := []byte(api.Address() + "," + apiCmd + ",")
	data = append(data, apiData...)
	cmd := command.New(command.ApiSendTo, data)
	cmd.Id = api.teo.getNextID()
	data, _ = cmd.MarshalBinary()
	api.teo.ws.SendMessage(data)
	id = cmd.Id
	return
}
