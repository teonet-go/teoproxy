// Teonet wasm client to use in wasm applications

package client

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	ws "github.com/teonet-go/teoproxy/ws/client"
	"github.com/teonet-go/teoproxy/ws/command"
)

type Teonet struct {
	ws *ws.WsClient // Websocket client
	id uint32       // Packet id
}

// Start Teonet client
func New(appShort string, onReconnected func()) (teo *Teonet, err error) {
	teo = new(Teonet)
	teo.ws = ws.NewWsClient(
	// func(message []byte) bool {
	// 	cmd := command.NewEmpty()
	// 	err = cmd.UnmarshalBinary(message)
	// 	if err != nil {
	// 		log.Println("Can't unmarshal teonet proxy server command, error:",
	// 			err, string(message))
	// 		return false
	// 	}
	// 	log.Println("Got Teonet proxy server command:", cmd.Cmd.String(),
	// 		string(cmd.Data))

	// 	return false
	// },
	)
	err = teo.ws.Connect(onReconnected)
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

func (teo *Teonet) NewAPIClient(peer string) (cli *APIClient, err error) {
	cmd := command.New(command.NewApiClient, []byte(peer))
	data, _ := cmd.MarshalBinary()
	teo.ws.SendMessage(data)
	cli = &APIClient{teo: teo, addr: peer}
	return
}

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
		log.Println("Got id", cmd.Id)
		if cmd.Id != id {
			return false
		}
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

type APIClient struct {
	teo  *Teonet
	addr string
}

func (api *APIClient) Address() string {
	return api.addr
}

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
