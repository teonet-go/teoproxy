package server

import (
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
	// log.Println("Received message (teonet proxy server):", message, string(message))

	// Check teonet command
	cmd := &command.TeonetCmd{}
	err := cmd.UnmarshalBinary(message)
	if err != nil {
		log.Println("Can't unmarshal teonet command, error:", err)
		return
	}
	log.Println("Get Teonet proxy command:", cmd.Cmd.String(), cmd.Data, string(cmd.Data))

	// Process command
	teo.processCommand(cmd)

	// Write response to client
	err = conn.WriteMessage(websocket.TextMessage, []byte("Message received"))
	if err != nil {
		log.Println("Can't write message to client, error:", err)
	}

}

func (teo *TeonetServer) processCommand(cmd *command.TeonetCmd) (data []byte, err error) {
	switch cmd.Cmd {
	case command.Connect:
		// Process Connect command
		// TODO: Add your code here
	case command.Disconnect:
		// Process Disconnect command
		// TODO: Add your code here
	case command.ConnectTo:
		// Process ConnectTo command
		// TODO: Add your code here
		// Connect to teoFortune server(peer)
		addr := string(cmd.Data)
		if err = teo.ConnectTo(addr); err != nil {
			err = fmt.Errorf("can't connect to 'fortune', error: %s" + err.Error())
		}
	case command.NewAPIClient:
		// Process NewAPIClient command
		// TODO: Add your code here
	default:
		err = fmt.Errorf("unknown command: %s", cmd.Cmd.String())
		log.Println("Unknown command:", err)
	}
	return
}
