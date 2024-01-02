//go:build wasm

package client

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"sync"
	"syscall/js"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

// WsClient is javascript websocket client to use in wasm application.
type WsClient struct {
	js.Value
	*Readers
}

func NewWsClient(processMessage ...ReaderFunc) *WsClient {
	ws := &WsClient{}
	ws.newReaders(processMessage...)
	return ws
}

func (ws *WsClient) Connect() (err error) {
	done := make(chan struct{}, 0)

	// Create a JavaScript WebSocket object
	js.Global().Set("socket", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			log.Println("Invalid number of arguments")
			return nil
		}

		url := args[0].String()
		log.Println("Url:", url)

		// Create a WebSocket connection
		ws.Value = js.Global().Get("WebSocket").New(url)

		// WebSocket open event handler
		ws.Value.Set("onopen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			log.Println("WebSocket connection established.")
			// Send a message through the WebSocket
			done <- struct{}{}
			return nil
		}))

		// WebSocket message event handler
		ws.Value.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			// Handle incoming messages from the server
			message := args[0].Get("data").String()
			// log.Println("Got message from server:", message)
			data, err := base64.StdEncoding.DecodeString(message)
			if err != nil {
				log.Println("Can't decode message base64, error:", err)
				return nil
			}

			// Process message
			ws.processReaders(data)

			return nil
		}))

		// WebSocket close event handler
		ws.Value.Set("onclose", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			log.Println("WebSocket connection closed.")
			// Handle WebSocket connection closure
			return nil
		}))

		// WebSocket error event handler
		ws.Value.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// message := args[0].Get("error").String()
			log.Println("WebSocket error:", args[0])
			// Handle WebSocket errors
			return nil
		}))

		return nil
	}))

	// Call the JavaScript function to create the WebSocket connection
	js.Global().Call("socket", "ws://localhost:8081/ws")

	<-done

	log.Println("WebSocket connection done.")
	return
}

// SendMessage sends a message to the websocket server
func (ws *WsClient) SendMessage(message []byte) {
	ws.Value.Call("send", base64.StdEncoding.EncodeToString(message))
	log.Println("Send message to server:", message)
}

// func (*WsClient) receiveMessages(conn *websocket.Conn) {
// 	for {
// 		// Read a message from the server
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Error receiving message from WebSocket server:", err)
// 			return
// 		}

// 		// Process the received message
// 		log.Println("Received message from server:", string(message))
// 	}
// }

type Readers struct {
	m map[string]ReaderFunc
	*sync.RWMutex
}
type ReaderFunc func(message []byte) bool

func (ws *WsClient) newReaders(processMessage ...ReaderFunc) {
	ws.Readers = &Readers{m: make(map[string]ReaderFunc), RWMutex: new(sync.RWMutex)}
	for _, p := range processMessage {
		ws.Readers.AddReader(p)
	}
	return
}

func (r *Readers) AddReader(processMessage ReaderFunc) (id string) {

	// newId generates random string with n bytes
	newId := func(n int) string {
		b := make([]byte, n)
		rand.Read(b)
		return base64.URLEncoding.EncodeToString(b)
	}

	r.Lock()

	// Create and check new id
	for {
		id = newId(16)
		if _, ok := r.m[id]; !ok {
			break
		}
	}

	r.m[id] = processMessage
	r.Unlock()

	return
}

func (r *Readers) RemoveReader(id string) {
	r.Lock()
	delete(r.m, id)
	r.Unlock()
}

func (r *Readers) getReaders(id string) (reader ReaderFunc, ok bool) {
	r.RLock()
	reader, ok = r.m[id]
	r.RUnlock()
	return
}

func (r *Readers) processReaders(message []byte) {
	r.RLock()
	defer r.RUnlock()
	for _, reader := range r.m {
		if reader(message) {
			break
		}
	}
}
