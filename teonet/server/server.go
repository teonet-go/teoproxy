package server

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/teonet-go/teoproxy/ws/command"
	ws "github.com/teonet-go/teoproxy/ws/server"
)

type TeonetServer struct {
	*ws.WsServer
}

func New() (teo *TeonetServer) {
	teo = &TeonetServer{}
	teo.WsServer = ws.New(teo.processMessage)
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

	// Write response to client
	err = conn.WriteMessage(websocket.TextMessage, []byte("Message received"))
	if err != nil {
		log.Println("Can't write message to client, error:", err)
	}

}
