package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/teonet-go/teoproxy/teonet/client"
)

//go:generate fyne package -os wasm

const appShort = "client-gui"

func main() {
	app := app.New()
	w := app.NewWindow("Hello")
	w.SetContent(widget.NewLabel("Hello Fyne!"))
	w.Resize(fyne.NewSize(200, 200))

	// Test connect to WebSocket server
	// ws := ws.NewWsClient()
	// err := ws.Connect()
	// if err == nil {
	// 	ws.SendMessage([]byte{1, 2, 3, 4, 5})
	// }

	const peer = "8agv3IrXQk7INHy5rVlbCxMWVmOOCoQgZBF" // teoFortune peer

	// Start Teonet proxy client
	teo, err := client.New(appShort)
	if err != nil {
		log.Fatal("can't initialize Teonet, error:", err)
	}
	// Connect to Teonet using proxy server
	err = teo.Connect()
	if err != nil {
		log.Fatal("can't connect to Teonet, error:", err)
	}
	// Connect to teoFortune server(peer)
	if err = teo.ConnectTo(peer); err != nil {
		log.Fatal("can't connect to peer, error:", err)
	}
	// Connet to fortune api
	api, err := teo.NewAPIClient(peer)
	if err != nil {
		log.Fatal("can't connect to peer api, error:", err)
		return
	}

	log.Println("connected to Teonet", api)

	w.ShowAndRun()
}
