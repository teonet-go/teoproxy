//go:build wasm

package client

import (
	"fmt"
	"log"
	"syscall/js"

	"github.com/gorilla/websocket"
)

// WsClient is javascript websocket client to use in wasm application.
type WsClient struct {
	js.Value
}

func NewWsClient() *WsClient {
	return &WsClient{}
}

func (ws *WsClient) Connect() {
	done := make(chan struct{}, 0)

	// Create a JavaScript WebSocket object
	js.Global().Set("socket", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			fmt.Println("Invalid number of arguments")
			return nil
		}

		url := args[0].String()
		fmt.Println("Url:", url)

		// Create a WebSocket connection
		ws.Value = js.Global().Get("WebSocket").New(url)

		// WebSocket open event handler
		ws.Value.Set("onopen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			fmt.Println("WebSocket connection established.")
			// Send a message through the WebSocket
			ws.Value.Call("send", "Hello, server! (inside js)")
			done <- struct{}{}
			return nil
		}))

		// WebSocket message event handler
		ws.Value.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			message := args[0].Get("data").String()
			fmt.Println("Message from server:", message)
			// Handle incoming messages from the server
			return nil
		}))

		// WebSocket close event handler
		ws.Value.Set("onclose", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			fmt.Println("WebSocket connection closed.")
			// Handle WebSocket connection closure
			return nil
		}))

		// WebSocket error event handler
		ws.Value.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			// message := args[0].Get("error").String()
			fmt.Println("WebSocket error:", args[0])
			// Handle WebSocket errors
			return nil
		}))

		return nil
	}))

	// Call the JavaScript function to create the WebSocket connection
	js.Global().Call("socket", "ws://localhost:8080/ws")

	<-done

	fmt.Println("WebSocket connection done.")
}

func (ws *WsClient) SendMessages(message string) {
	// Send a message to the websocket server
	ws.Value.Call("send", message)
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
