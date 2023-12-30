// Teonet wasm client to use in wasm applications

package client

import (
	"fmt"

	ws "github.com/teonet-go/teoproxy/ws/client"
	"github.com/teonet-go/teoproxy/ws/command"
)

type Teonet struct {
	ws *ws.WsClient
}

// Start Teonet client
func New(appShort string) (teo *Teonet, err error) {
	teo = new(Teonet)
	teo.ws = ws.NewWsClient()
	err = teo.ws.Connect()
	return
}

// Connect to Teonet
func (teo *Teonet) Connect() (err error) {
	cmd := &command.TeonetCmd{Cmd: command.Connect}
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	return
}

func (teo *Teonet) Disconnect() (err error) {
	cmd := &command.TeonetCmd{Cmd: command.Disconnect}
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	return
}

func (teo *Teonet) ConnectTo(peer string) (err error) {
	cmd := &command.TeonetCmd{Cmd: command.ConnectTo, Data: []byte(peer)}
	data, _ := cmd.MarshalBinary()
	fmt.Println("SendMessage ConnectTo:", data)
	teo.ws.SendMessage(data)
	return
}

func (teo *Teonet) NewAPIClient(peer string) (api *ws.WsClient, err error) {
	cmd := &command.TeonetCmd{Cmd: command.NewAPIClient, Data: []byte(peer)}
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	return
}
