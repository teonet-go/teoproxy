// Copyright 2023-2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !wasm

// Teonet client to use in native applications.
package client

import (
	"fmt"

	"github.com/teonet-go/teonet"
)

// Teonet is a wrapper struct that embeds a teonet.Teonet client.
// It allows extending the base teonet.Teonet client with additional methods.
type Teonet struct {
	*teonet.Teonet
}

// New starts Teonet client and returns a new Teonet instance.
//
// appShort - short application name
// onReconnected - callback that will be called after successful reconnection
//
// Returns a new Teonet instance and error if any.
func New(appShort string, onReconnected func()) (teo *Teonet, err error) {

	teo = new(Teonet)

	// Start Teonet client
	teo.Teonet, err = teonet.New(appShort)
	if err != nil {
		err = fmt.Errorf("can't init Teonet, error: " + err.Error())
		return
	}

	return
}

// APIClient wraps a teonet.APIClient to provide additional methods.
type APIClient struct{ *teonet.APIClient }

// NewAPIClient creates a new instance of the APIClient struct and returns it
// along with any error that occurred.
//
// Parameters:
// - addr: The address to connect to.
//
// Returns:
// - cli: A pointer to the APIClient struct.
// - err: Any error that occurred during the creation of the APIClient.
func (teo *Teonet) NewAPIClient(addr string) (cli *APIClient, err error) {
	apicli, err := teo.Teonet.NewAPIClient(addr)
	cli = &APIClient{apicli}
	return
}
