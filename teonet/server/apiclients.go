// Copyright 2023-2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// APIClients module of the Server package.

package server

import (
	"sync"

	"github.com/teonet-go/teonet"
)

// APIClients stores a map of APIClient instances, keyed by peer name.
// It uses a RWMutex for concurrent access control.
type APIClients struct {
	m map[string]*teonet.APIClient
	*sync.RWMutex
}

// Init initializes the APIClients map and mutex.
func (cli *APIClients) Init() *APIClients {
	cli.m = make(map[string]*teonet.APIClient)
	cli.RWMutex = new(sync.RWMutex)
	return cli
}

// Add adds a new APIClient instance to the APIClients map,
// keyed by the provided name. It locks the map during the update
// to prevent concurrent map writes. It first checks if a client
// already exists for the given name and returns immediately if
// so to avoid overwriting the existing client.
func (cli *APIClients) Add(name string, api *teonet.APIClient) {
	cli.Lock()
	defer cli.Unlock()

	if _, ok := cli.m[name]; ok {
		return
	}

	cli.m[name] = api
}

// Remove removes the APIClient instance for the given name from the APIClients
// map. It locks the map during the update to prevent concurrent map writes.
func (cli *APIClients) Remove(name string) {
	cli.Lock()
	defer cli.Unlock()
	delete(cli.m, name)
}

// Get retrieves the APIClient instance for the given name from the
// APIClients map. It locks the map for reading during the lookup to
// prevent concurrent map access. The second return value indicates
// if a client was found. This is an exported method that is part of
// the APIClients API.
func (cli *APIClients) Get(name string) (api *teonet.APIClient, ok bool) {
	cli.RLock()
	defer cli.RUnlock()
	api, ok = cli.m[name]
	return
}

// Exists checks if an APIClient with the given name exists in the APIClients
// map. It calls the Get method and checks if it returned a client.
func (cli *APIClients) Exists(name string) bool {
	_, ok := cli.Get(name)
	return ok
}
