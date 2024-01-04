//go:build wasm

package client

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"syscall/js"
	"time"
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

func (ws *WsClient) Connect(onReconnected func()) (err error) {
	var connected bool
	done := make(chan struct{}, 0)

	// wsScheme returns "ws" or "wss" depending on the given URL scheme
	wsScheme := func(httpScheme string) string {
		if httpScheme == "https" {
			return "wss"
		}
		return "ws"
	}

	// Get the current URL and parse it to create the WebSocket URL
	href := js.Global().Get("location").Get("href")
	u, err := url.Parse(href.String())
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("%s://%s/ws", wsScheme(u.Scheme), u.Host)
	log.Println("Connect to websocket:", url)

	// Call the JavaScript function to create the WebSocket connection
	connect := func() {
		js.Global().Call("socket", url)
	}

	// Create a JavaScript WebSocket object
	js.Global().Set("socket", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			log.Println("Invalid number of arguments")
			return nil
		}

		url := args[0].String()
		log.Println("Connect to websocket:", url)

		// Create a WebSocket connection
		ws.Value = js.Global().Get("WebSocket").New(url)

		// WebSocket open event handler
		ws.Value.Set("onopen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			log.Println("WebSocket connection established.")
			if !connected {
				connected = true
				done <- struct{}{}
			} else {
				onReconnected()
			}
			return nil
		}))

		// WebSocket message event handler
		ws.Value.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			// Handle incoming messages from the server
			message := args[0].Get("data").String()
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
			time.Sleep(1 * time.Second)

			// Reconnect
			connect()
			return nil
		}))

		// WebSocket error event handler
		ws.Value.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			message := args[0].Get("error").String()
			log.Println("WebSocket error:", message)
			return nil
		}))

		return nil
	}))

	// Call the JavaScript function to create the WebSocket connection
	connect()

	// Wait for the WebSocket connection to be established or timeout
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		err = fmt.Errorf("timeout")
	}

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
