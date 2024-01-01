// Teonet wasm client to use in wasm applications

package client

import (
	"log"

	ws "github.com/teonet-go/teoproxy/ws/client"
	"github.com/teonet-go/teoproxy/ws/command"
)

type Teonet struct {
	ws *ws.WsClient
}

// Start Teonet client
func New(appShort string) (teo *Teonet, err error) {
	teo = new(Teonet)
	teo.ws = ws.NewWsClient(func(message []byte) {
		cmd := command.NewEmpty()
		err = cmd.UnmarshalBinary(message)
		if err != nil {
			log.Println("Can't unmarshal teonet proxy server command, error:",
				err, string(message))
			return
		}
		log.Println("Got Teonet proxy server command:", cmd.Cmd.String(),
			/* cmd.Data,  */ string(cmd.Data))
	})
	err = teo.ws.Connect()
	return
}

// Connect to Teonet
func (teo *Teonet) Connect() (err error) {
	cmd := command.New(command.Connect, nil)
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	return
}

func (teo *Teonet) Disconnect() (err error) {
	cmd := command.New(command.Disconnect, nil)
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	return
}

func (teo *Teonet) ConnectTo(peer string) (err error) {
	cmd := command.New(command.ConnectTo, []byte(peer))
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	return
}

func (teo *Teonet) NewAPIClient(peer string) (err error) {
	cmd := command.New(command.NewAPIClient, []byte(peer))
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	return
}

func (teo *Teonet) SendTo(apiCmd string, apiData []byte) (err error) {
	cmd := command.New(command.SendTo, []byte(apiCmd))
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	return
}
