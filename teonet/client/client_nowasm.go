// Teonet wasm client to use in wasm applications

//go:build !wasm

package client

import (
	"fmt"

	"github.com/teonet-go/teonet"
)

type Teonet struct {
	*teonet.Teonet
}

// New starts Teonet client
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

type APIClient struct{ *teonet.APIClient }

func (teo *Teonet) NewAPIClient(addr string) (cli *APIClient, err error) {
	apicli, err := teo.Teonet.NewAPIClient(addr)
	cli = &APIClient{apicli}
	return
}
