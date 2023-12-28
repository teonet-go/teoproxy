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
	ws.Connect()
	ws.SendMessages("Hello, server!")

	w.ShowAndRun()
}
