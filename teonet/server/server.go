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

	// Init api clients object
	teo.initAPIClients()

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
	teo.WsServer = ws.New(teo.processMessage)

	return
}

// processMessage processes a websocket message received from a client.
// It decodes the base64 encoded message, unmarshals the teonet command,
// processes the command by calling processCommand, and writes the response
// back to the client.
func (teo *TeonetServer) processMessage(conn *websocket.Conn, message []byte) {

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

// processCommand processes a Teonet command received from a client.
// It handles different command types like Connect, Disconnect etc.
// Returns the response data and error.
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
		teo.Lock()
		defer teo.Unlock()
		addr := string(cmd.Data)
		if err = teo.ConnectTo(addr); err != nil {
			err = fmt.Errorf("can't connect to peer %s, error: %s", addr, err)
			log.Println(err)
			return
		}
		str := fmt.Sprintf("Connected to peer %s", addr)
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
		str := fmt.Sprintf("Connected to peer %s api", addr)
		data = []byte(str)
		log.Println(str)

	// Process SendTo command
	case command.ApiSendTo:

		// Split commands data to peer name and api command
		splitData := strings.Split(string(cmd.Data), ",")
		if len(splitData) < 3 {
			err = fmt.Errorf("wrong command data: %s", cmd.Cmd.String())
			return
		}

		apiPeerName := splitData[0]
		apiCommand := splitData[1]
		apiCommandData := cmd.Data[len(apiPeerName)+1+len(apiCommand)+1:]

		log.Println("Send api command:", string(apiCommand), " to peer:",
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
			log.Println("Got response from peer, len:", len(data), " err:", err)
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
		log.Println("Unknown command:", err)
	}

	return
}

// APIClients stores a map of APIClient instances, keyed by peer name.
// It uses a RWMutex for concurrent access control.
type APIClients struct {
	m map[string]*teonet.APIClient
	*sync.RWMutex
}

// initAPIClients initializes the apiClients field of the TeonetServer.
// It creates a new APIClients instance to store API client connections
// in a concurrent map, protected by an RWMutex.
func (teo *TeonetServer) initAPIClients() {
	teo.apiClients = &APIClients{
		m:       make(map[string]*teonet.APIClient),
		RWMutex: &sync.RWMutex{},
	}
}

// Add adds a new APIClient instance to the APIClients map,
// keyed by the provided name. It locks the map during the update
// to prevent concurrent map writes. It first checks if a client
// already exists for the given name and returns immediately if
// so to avoid overwriting the existing client.
func (cli *APIClients) Add(name string, api *teonet.APIClient) {
	cli.Lock()
	defer cli.Unlock()

	if _, ok := cli.m[name]; ok {
		return
	}

	cli.m[name] = api
}

// Remove removes the APIClient instance for the given name from the APIClients
// map. It locks the map during the update to prevent concurrent map writes.
func (cli *APIClients) Remove(name string) {
	cli.Lock()
	defer cli.Unlock()
	delete(cli.m, name)
}

// Get retrieves the APIClient instance for the given name from the
// APIClients map. It locks the map for reading during the lookup to
// prevent concurrent map access. The second return value indicates
// if a client was found. This is an exported method that is part of
// the APIClients API.
func (cli *APIClients) Get(name string) (api *teonet.APIClient, ok bool) {
	cli.RLock()
	defer cli.RUnlock()
	api, ok = cli.m[name]
	return
}

// Exists checks if an APIClient with the given name exists in the APIClients
// map. It calls the Get method and checks if it returned a client.
func (cli *APIClients) Exists(name string) bool {
	_, ok := cli.Get(name)
	return ok
}
