// Teonet wasm client to use in wasm applications

package client

import (
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
	teo.ws.SendMessage([]byte{command.Connect})
	return
}

func (teo *Teonet) Disconnect() (err error) {
	teo.ws.SendMessage([]byte{command.Dsconnect})
	return
}

func (teo *Teonet) ConnectTo(peer string) (err error) {
	data := append([]byte{command.ConnectTo}, []byte(peer)...)
	teo.ws.SendMessage(data)
	return
}

func (teo *Teonet) NewAPIClient(peer string) (api *ws.WsClient, err error) {
	data := append([]byte{command.NewAPIClient}, []byte(peer)...)
	teo.ws.SendMessage(data)
	return
}
