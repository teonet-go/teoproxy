package server

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/teonet-go/teonet"
	"github.com/teonet-go/teoproxy/ws/command"
	ws "github.com/teonet-go/teoproxy/ws/server"
)

type TeonetServer struct {
	*ws.WsServer
	*teonet.Teonet
	apiClients *APIClients
}

func New(appShort string) (teo *TeonetServer, err error) {
	teo = &TeonetServer{}

	// Init api clients object
	teo.initAPIClients()

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
		data = []byte(fmt.Sprintf("Connected to peer %s api", addr))

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

		log.Println("Send api command:", string(apiCommand), " to peer:", apiPeerName, " data len:", len(apiCommandData))

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

type APIClients struct {
	m map[string]*teonet.APIClient
	*sync.RWMutex
}

func (teo *TeonetServer) initAPIClients() {
	teo.apiClients = &APIClients{
		m:       make(map[string]*teonet.APIClient),
		RWMutex: &sync.RWMutex{},
	}
}

func (cli *APIClients) Add(name string, api *teonet.APIClient) {
	cli.Lock()
	defer cli.Unlock()

	if _, ok := cli.m[name]; ok {
		return
	}

	cli.m[name] = api
}

func (cli *APIClients) Remove(name string) {
	cli.Lock()
	defer cli.Unlock()
	delete(cli.m, name)
}

func (cli *APIClients) Get(name string) (api *teonet.APIClient, ok bool) {
	cli.RLock()
	defer cli.RUnlock()
	api, ok = cli.m[name]
	return
}

func (cli *APIClients) Exists(name string) bool {
	_, ok := cli.Get(name)
	return ok
}
