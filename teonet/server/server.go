// Copyright 2023-2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Server package contains the server implementation for the Teonet proxy.
package server

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/teonet-go/teomon"
	"github.com/teonet-go/teonet"
	"github.com/teonet-go/teoproxy/ws/command"
	ws "github.com/teonet-go/teoproxy/ws/server"
)

// init initializes the Go program.
//
// It sets the log flags to include the standard date and time format, as well
// as microseconds.
func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

// TeonetServer is the main server type that contains the core components.
// It has mutexes for synchronization, the websocket server,
// the Teonet client, and API clients. As an exported type, it is part of the
// public API.
type TeonetServer struct {
	*sync.Mutex
	*ws.WsServer
	*teonet.Teonet
	apiClients *APIClients
	stream     *StreamAnswer
}

// TeonetMonitor contains monitoring information to send to the Teonet monitor.
type TeonetMonitor struct {
	Addr         string
	AppName      string
	AppShort     string
	AppVersion   string
	TeoVersion   string
	AppStartTime time.Time
}

// New creates a new TeonetServer instance. It initializes the mutex, API clients,
// Teonet client, and websocket server. The appShort parameter specifies the
// application name. The monitor parameter optionally configures connecting to a
// Teonet monitor for metrics reporting. It returns the TeonetServer instance
// and any error. As an exported function, this serves as the main constructor for
// the TeonetServer type.
func New(appShort string, monitor *TeonetMonitor) (teo *TeonetServer, err error) {
	teo = &TeonetServer{Mutex: new(sync.Mutex)}

	// Init apiClients and stream objects
	teo.newAPIClients()

	// Start Teonet client
	teo.Teonet, err = teonet.New(appShort, teo.reader)
	if err != nil {
		return
	}

	// Connect to Teonet
	err = teo.Connect()
	if err != nil {
		err = fmt.Errorf("can't connect to Teonet, error: " + err.Error())
		return
	}

	// Connect to monitor
	if monitor != nil && len(monitor.Addr) > 0 {
		teomon.Connect(teo.Teonet, monitor.Addr, teomon.Metric{
			AppName:      monitor.AppName,
			AppShort:     monitor.AppShort,
			AppVersion:   monitor.AppVersion,
			TeoVersion:   teonet.Version,
			AppStartTime: time.Now(),
		})
	}
	log.Println("Connected to monitor")

	// Create websocket server
	teo.WsServer = ws.New(
		// On websocket disconnect func
		func(conn *websocket.Conn) {
			teo.stream.RemoveConn(conn)
		},
		// On websocket message functions
		teo.processMessage,
	)

	return
}

func (teo *TeonetServer) reader(c *teonet.Channel, p *teonet.Packet, e *teonet.Event) bool {

	if e.Err != nil {
		return false
	}

	if e.Event != teonet.EventData {
		return false
	}

	log.Printf("got packet in reader: id: %d, data len: %d, from %s", p.ID(),
		len(p.Data()), c)

	peer := c.String()
	streem := strings.Split(string(p.Data()), "/")[0]
	if conns, ok := teo.stream.Get(peer, streem); ok {
		cmd := &command.TeonetCmd{}
		cmd.Data = p.Data()
		data, _ := cmd.MarshalBinary()
		for _, conn := range conns {
			if err := teo.WriteMessage(conn, websocket.TextMessage,
				[]byte(base64.StdEncoding.EncodeToString(data))); err != nil {
				log.Printf("can't write message to client %p, error: %s", conn,
					err)
			} else {
				log.Printf("message sent to client %p, data len: %v", conn,
					len(data))
			}
		}
	}

	return false
}

func (teo *TeonetServer) WriteMessage(conn *websocket.Conn, messageType int,
	data []byte) error {
	teo.Lock()
	defer teo.Unlock()
	return conn.WriteMessage(messageType, data)
}

// processMessage processes a websocket message received from a client.
// It decodes the base64 encoded message, unmarshals the teonet command,
// processes the command by calling processCommand, and writes the response
// back to the client.
func (teo *TeonetServer) processMessage(conn *websocket.Conn, message []byte) {

	// decode message base64
	message, err := base64.StdEncoding.DecodeString(string(message))
	if err != nil {
		log.Println("can't decode message base64, error:", err)
		return
	}

	// Check teonet command
	cmd := &command.TeonetCmd{}
	err = cmd.UnmarshalBinary(message)
	if err != nil {
		log.Println("can't unmarshal teonet command, error:", err, string(message))
		return
	}
	log.Println("got Teonet proxy client command:", cmd.Id, cmd.Cmd.String(),
		string(cmd.Data))

	// Process command
	data, err := teo.processCommand(cmd, conn)
	if err != nil {
		log.Println("process command, error:", err)
		return
	}

	// Write response to client
	cmd.Data, cmd.Err = data, err
	data, _ = cmd.MarshalBinary()
	if err = teo.WriteMessage(conn, websocket.TextMessage,
		[]byte(base64.StdEncoding.EncodeToString(data))); err != nil {
		log.Println("can't write message to client, error:", err)
	}
}

// processCommand processes a Teonet command received from a client.
// It handles different command types like Connect, Disconnect etc.
// Returns the response data and error.
func (teo *TeonetServer) processCommand(cmd *command.TeonetCmd, conn *websocket.Conn) (data []byte,
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
		teo.Lock()
		defer teo.Unlock()
		addr := string(cmd.Data)
		if err = teo.ConnectTo(addr); err != nil {
			err = fmt.Errorf("can't connect to peer %s, error: %s", addr, err)
			log.Println(err)
			return
		}
		str := fmt.Sprintf("connected to peer %s", addr)
		data = []byte(str)
		log.Println(str)

	// Process NewAPIClient command
	case command.NewApiClient:
		addr := string(cmd.Data)
		if !teo.apiClients.Exists(addr) {
			var cli *teonet.APIClient
			if cli, err = teo.NewAPIClient(addr); err != nil {
				err = fmt.Errorf("can't connect to peer %s api, error: %s",
					addr, err.Error())
				return
			}
			teo.apiClients.Add(addr, cli)
		}
		str := fmt.Sprintf("connected to peer %s api", addr)
		data = []byte(str)
		log.Println(str)

	// Process sream command
	case command.Stream:
		// Split commands data to peer name and api command
		splitData := strings.Split(string(cmd.Data), ",")
		if len(splitData) < 2 {
			err = fmt.Errorf("wrong command data: %s", string(cmd.Data))
			return
		}

		peer := splitData[0]
		stream := splitData[1]

		teo.stream.Add(peer, stream, conn)
		log.Printf("stream added, conn: %p, stream: %v", conn, teo.stream)

	// Process SendTo command
	case command.ApiSendTo:

		// Split commands data to peer name and api command
		splitData := strings.Split(string(cmd.Data), ",")
		if len(splitData) < 3 {
			err = fmt.Errorf("wrong command data: %s", string(cmd.Data))
			return
		}

		apiPeerName := splitData[0]
		apiCommand := splitData[1]
		apiCommandData := cmd.Data[len(apiPeerName)+1+len(apiCommand)+1:]

		log.Println("send api command:", string(apiCommand), " to peer:",
			apiPeerName, " data len:", len(apiCommandData))

		// Api answer struct
		type apiAnswer struct {
			data []byte
			err  error
		}
		w := make(chan apiAnswer, 1)
		// Get api client by name
		api, ok := teo.apiClients.Get(apiPeerName)
		if !ok {
			err = fmt.Errorf(
				"can't get api client, error: has not connected to peer api %s",
				apiPeerName,
			)
			return
		}
		// Send request to api peer
		api.SendTo(apiCommand, apiCommandData, func(data []byte, err error) {
			log.Println("got response from peer, len:", len(data), " err:", err)
			w <- apiAnswer{data, err}
		})

		// Get answer from api peer or timeout
		var answer apiAnswer
		select {
		case answer = <-w:
		case <-time.After(5 * time.Second):
			answer = apiAnswer{nil, fmt.Errorf("timeout")}
		}
		data, err = answer.data, answer.err

	// Unknown command
	default:
		err = fmt.Errorf("unknown command: %s", cmd.Cmd.String())
		log.Println("unknown command:", err)
	}

	return
}

// newAPIClients creates and initializes new APIClients and new StreamAnswer.
func (teo *TeonetServer) newAPIClients() {
	teo.apiClients = new(APIClients).Init()
	teo.stream = new(StreamAnswer).Init()
}
