package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/teonet-go/teoproxy/ws/client"
)

//go:generate fyne package -os wasm

func main() {
	app := app.New()
	w := app.NewWindow("Hello")
	w.SetContent(widget.NewLabel("Hello Fyne!"))
	w.Resize(fyne.NewSize(200, 200))

	ws := client.NewWsClient()
	err := ws.Connect()
	if err == nil {
		ws.SendMessage([]byte{1, 2, 3, 4, 5})
	}

	w.ShowAndRun()
}
