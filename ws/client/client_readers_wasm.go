// Copyright 2023-2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm

// Package client provides client-side WebSocket functionality.

package client

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
)

// Readers holds the map of ReaderFunc callbacks and a RWMutex for concurrent access.
type Readers struct {
	m map[string]ReaderFunc
	*sync.RWMutex
}
type ReaderFunc func(message []byte) bool

// newReaders initializes the readers for the WsClient.
//
// The processMessage parameter is a variadic parameter that accepts one or
// more ReaderFunc functions. These functions are used to process incoming
// messages in the WsClient.
//
// This function does not return any values.
func (ws *WsClient) newReaders(processMessage ...ReaderFunc) {
	ws.Readers = &Readers{m: make(map[string]ReaderFunc), RWMutex: new(sync.RWMutex)}
	for _, p := range processMessage {
		ws.Readers.AddReader(p)
	}
	return
}

// AddReader adds a new reader to the Readers map and returns the generated id.
//
// The function takes a ReaderFunc as a parameter, which is the function to be
// executed when a message is received by the reader.
// It returns the generated readers id as a string.
func (r *Readers) AddReader(processMessage ReaderFunc) (id string) {

	r.Lock()
	defer r.Unlock()

	// newId generates random string with n bytes
	newId := func(n int) string {
		b := make([]byte, n)
		rand.Read(b)
		return base64.URLEncoding.EncodeToString(b)
	}

	// Create and check new id
	for {
		id = newId(16)
		if _, ok := r.m[id]; !ok {
			break
		}
	}

	// Add new reader to map
	r.m[id] = processMessage

	return
}

// RemoveReader removes a reader from the Readers struct by the given id.
//
// Parameters:
//   - id: The id of the reader to be removed.
func (r *Readers) RemoveReader(id string) {
	r.Lock()
	defer r.Unlock()
	delete(r.m, id)
}

// getReaders retrieves a reader function by its ID from the Readers struct.
//
// Parameters:
// - id: the ID of the reader to retrieve.
//
// Returns:
// - reader: the reader function associated with the given ID.
// - ok: a boolean indicating whether the reader was found.
func (r *Readers) getReaders(id string) (reader ReaderFunc, ok bool) {
	r.RLock()
	defer r.RUnlock()
	reader, ok = r.m[id]
	return
}

// processReaders processes the readers in the Readers struct.
//
// It takes a message of type []byte as a parameter.
// It does not return any value.
func (r *Readers) processReaders(message []byte) {
	r.RLock()
	defer r.RUnlock()
	for _, reader := range r.m {
		if reader(message) {
			break
		}
	}
}
