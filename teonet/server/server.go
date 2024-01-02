package server

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/teonet-go/teonet"
	"github.com/teonet-go/teoproxy/ws/command"
	ws "github.com/teonet-go/teoproxy/ws/server"
)

type TeonetServer struct {
	*ws.WsServer
	*teonet.Teonet
	*teonet.APIClient
}

func New(appShort string) (teo *TeonetServer, err error) {
	teo = &TeonetServer{}

	// Create websocket server
	teo.WsServer = ws.New(teo.processMessage)

	// Start Teonet client
	teo.Teonet, err = teonet.New(appShort)
	if err != nil {
		return
	}

	// Connect to Teonet
	err = teo.Connect()
	if err != nil {
		err = fmt.Errorf("can't connect to Teonet, error: " + err.Error())
		return
	}

	return
}

func (teo *TeonetServer) processMessage(conn *websocket.Conn, message []byte) {

	// Process message logic here
	// log.Println("Received message (teonet proxy server):", len(message), string(message))

	// decode message base64
	message, err := base64.StdEncoding.DecodeString(string(message))
	if err != nil {
		log.Println("Can't decode message base64, error:", err)
		return
	}

	// Check teonet command
	cmd := &command.TeonetCmd{}
	err = cmd.UnmarshalBinary(message)
	if err != nil {
		log.Println("Can't unmarshal teonet command, error:", err, string(message))
		return
	}
	log.Println("Got Teonet proxy client command:", cmd.Id, cmd.Cmd.String(),
		string(cmd.Data))

	// Process command
	data, err := teo.processCommand(cmd)
	if err != nil {
		log.Println("process command, error:", err)
		return
	}

	// Write response to client
	cmd.Data, cmd.Err = data, err
	data, _ = cmd.MarshalBinary()
	if err = conn.WriteMessage(websocket.TextMessage,
		[]byte(base64.StdEncoding.EncodeToString(data))); err != nil {
		log.Println("Can't write message to client, error:", err)
	}
}

func (teo *TeonetServer) processCommand(cmd *command.TeonetCmd) (data []byte,
	err error) {

	switch cmd.Cmd {

	// Process Connect command
	case command.Connect:
		data = []byte("Connected to Teonet")

	// Process Disconnect command
	case command.Disconnect:
		// TODO: Add your code here

	// Process ConnectTo peer command
	case command.ConnectTo:
		addr := string(cmd.Data)
		if err = teo.ConnectTo(addr); err != nil {
			err = fmt.Errorf("can't connect to peer %s, error: %s", addr, err)
			log.Println(err)
			return
		}
		data = []byte(fmt.Sprintf("Connected to peer %s", addr))

	// Process NewAPIClient command
	case command.NewAPIClient:
		addr := string(cmd.Data)
		if teo.APIClient, err = teo.NewAPIClient(addr); err != nil {
			err = fmt.Errorf("can't connect to peer %s api, error: %s", addr,
				err.Error())
			return
		}
		data = []byte(fmt.Sprintf("Connected to peer %s api", addr))

	// Process SendTo command
	case command.SendTo:
		type apiAnswer struct {
			data []byte
			err  error
		}
		w := make(chan apiAnswer, 1)
		teo.APIClient.SendTo(string(cmd.Data), nil, func(data []byte, err error) {
			log.Println("Got response from peer, len:", len(data), " err:", err)
			w <- apiAnswer{data, err}
		})
		answer := <-w
		data, err = answer.data, answer.err

	// Unknown command
	default:
		err = fmt.Errorf("unknown command: %s", cmd.Cmd.String())
		log.Println("Unknown command:", err)
	}
	return
}
