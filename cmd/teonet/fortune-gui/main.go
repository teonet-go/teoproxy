// Copyright 2023-2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example:This a GUI application in Go.
//
// The main function is the entry point of the program. It initializes a
// connection to the Teonet network and Teofortune service, and then launches
// the graphical user interface (GUI).
//
// It first tries to connect to Teonet and Teofortune by calling newTeofortune,
// passing the application name and Teofortune address. newTeofortune returns a
// teofortune object that contains the Teonet client and Teofortune API client.
// If there is an error connecting, it prints the error and exits the program.
//
// If the connection is successful, it calls the newGui method on the teofortune
// object. This will create and display the GUI window.
//
// The main purpose of main is to initialize the networking connections and
// launch the user interface. It takes no direct inputs. Its output is to start
// up the GUI application.
//
// It achieves this by first initializing the back-end networking using
// newTeofortune. This sets up the ability to connect to Teonet and Teofortune.
// Once that is done, it launches the front-end GUI by calling newGui.
//
// The main logic flow is:
//
// - Initialize networking
// - If no errors, launch GUI
// - If errors, print error and exit
//
// So in summary, main handles setting up the app and starting it up,
// delegating the specific networking and GUI code to other functions. This is
// a common structure for main functions - initialize, then launch.
package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/teonet-go/teoproxy/teonet/client"
)

const (
	appName    = "Teonet fortune golang GUI application"
	appShort   = "teofortunegui"
	appVersion = "0.0.3"
	teoFortune = "8agv3IrXQk7INHy5rVlbCxMWVmOOCoQgZBF"
)

// main is the entry point of the program.
//
// It then connects to the Teonet and Teofortune server using the newTeofortune function.
// If there is an error connecting to Teonet, it prints an error message and returns.
// Finally, it creates and runs the gui interface using the newGui function.
func main() {

	// Connect to Teonet and Teofortune server
	teo, err := newTeofortune(appShort, teoFortune)
	if err != nil {
		log.Println("can't connect to Teonet, error:", err)
		return
	}

	// Create and run gui interface
	teo.newGui()
}

// teofortune contains teonet data and holds methods to start gui, process
// teonet connection and teofortune api
type teofortune struct {
	*client.Teonet                   // Teonet object
	addr           string            // Teofortune address
	client         *client.APIClient // Teofortune api client
}

// newTeofortune initializes a new teofortune instance with a Teonet client
// and API client to connect to the Teonet network and Teofortune service.
// It takes the application short name and Teofortune server address as arguments.
//
// It creates a new teofortune instance and sets the Teofortune address.
// It initializes the Teonet client with an onConnected callback.
// In the callback it connects to the Teonet network and Teofortune service,
// and initializes the Teofortune API client.
// It returns the teofortune instance and any error.
//
// This allows creating a Teofortune client instance connected to Teonet.
func newTeofortune(appShort, teoFortune string) (teo *teofortune, err error) {

	teo = new(teofortune)
	teo.addr = teoFortune

	// On Teonet client connected or reconnected
	onConnected := func() {

		log.Println("Teonet client connected.")

		// Connect to Teonet
		err = teo.Connect()
		if err != nil {
			// time.Sleep(1 * time.Second)
			// goto connect
			err = fmt.Errorf("can't connect to Teonet, error: " + err.Error())
			return
		}

		// Connect to teoFortune server(peer)
		if err = teo.ConnectTo(teo.addr); err != nil {
			err = fmt.Errorf("can't connect to 'fortune', error: %s" + err.Error())
			return
		}

		// Connet to fortune api
		if teo.client, err = teo.NewAPIClient(teo.addr); err != nil {
			err = fmt.Errorf("can't connect to 'fortune' api, error: %s", err.Error())
			return
		}

	}

	// Start Teonet client
	teo.Teonet, err = client.New(appShort, onConnected)
	if err != nil {
		err = fmt.Errorf("can't init Teonet, error: " + err.Error())
		return
	}
	onConnected()

	return
}

// newGui creates and displays the GUI window for the fortune application.
// It creates the window, sets up the widgets, handles user interaction
// to refresh the fortune message, and displays the window.
func (teo *teofortune) newGui() {
	a := app.New()
	w := a.NewWindow("Teofortune")

	label := widget.NewLabel("Fortune message from Teofortune server:")
	fmsg, _ := teo.fortune()
	message := widget.NewLabel(fmsg)
	w.SetContent(container.NewVBox(
		label,
		widget.NewButton("Show next", func() {
			fmsg, _ := teo.fortune()
			message.SetText(fmsg)
		}),
		message,
	))

	w.Resize(fyne.Size{Width: 600, Height: 600})
	w.ShowAndRun()
}

// fortune gets fortune messsage from teofortune microservice.
func (teo *teofortune) fortune() (msg string, err error) {

	// Get fortune message from teofortune microservice
	id, err := teo.client.SendTo("fortb", nil)
	if err != nil {
		return
	}
	log.Println("Send id", id, "ApiSendTo")
	data, err := teo.WaitFrom(teo.client.Address(), uint32(id))
	if err != nil {
		return
	}

	msg = string(data)
	return
}
