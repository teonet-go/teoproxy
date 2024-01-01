// Teonet wasm client to use in wasm applications

package client

import (
	"log"
	"sync/atomic"

	ws "github.com/teonet-go/teoproxy/ws/client"
	"github.com/teonet-go/teoproxy/ws/command"
)

type Teonet struct {
	ws *ws.WsClient // Websocket client
	id uint32       // Packet id
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

// getNextID gets next packet id
func (teo *Teonet) getNextID() uint32 {
	return atomic.AddUint32(&teo.id, 1)
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

func (teo *Teonet) SendTo(apiCmd string, apiData []byte) (id uint32, err error) {
	cmd := command.New(command.SendTo, []byte(apiCmd))
	cmd.Id = teo.getNextID()
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	id = cmd.Id
	return
}

func (teo *Teonet) WaitFrom(peer string, id uint32) (data []byte, err error) {
	// TODO: implement this function code
	return
}
