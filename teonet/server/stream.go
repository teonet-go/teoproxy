// Copyright 2023-2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Answer stream module of the Server package.

package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

// StreamAnswer stores a map of stream answers, keyed by peer name and stream
// name. The stream name has next format: "stream_name/96", where stream_name
// is the stream name and /96 is the stream parameters.
type StreamAnswer struct {
	m map[string]map[string][]*websocket.Conn
	*sync.RWMutex
}

// Init initializes the StreamAnswer map and mutex.
func (s *StreamAnswer) Init() *StreamAnswer {
	s.m = make(map[string]map[string][]*websocket.Conn)
	s.RWMutex = new(sync.RWMutex)
	return s
}

// Add adds a new stream answer to the StreamAnswer map, keyed by the provided
// peer name and stream name. It locks the map during the update to prevent
// concurrent map writes. It first checks if a stream answer already exists
// for the given peer and stream name and returns immediately if so to avoid
// overwriting the existing answer.
func (s *StreamAnswer) Add(peer, stream string, conn *websocket.Conn) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.m[peer]; !ok {
		s.m[peer] = make(map[string][]*websocket.Conn)
	}
	// if _, ok := s.m[peer][stream]; ok {
	// 	return
	// }
	if !s.existsUnsafe(peer, stream, conn) {
		s.m[peer][stream] = append(s.m[peer][stream], conn)
	}
}

// existsUnsafe checks if the connection already exists in the given peer and 
// stream array.
func (s *StreamAnswer) existsUnsafe(peer, stream string, conn *websocket.Conn) bool {
	if _, ok := s.m[peer]; !ok {
		return false
	}
	// if _, ok := s.m[peer][stream]; !ok {
	// 	return false
	// }
	for _, c := range s.m[peer][stream] {
		if c == conn {
			return true
		}
	}
	return false
}

// RemoveConn removes the stream answer by connection from the StreamAnswer map. It
// locks the map during the update to prevent concurrent map writes. 
func (s *StreamAnswer) RemoveConn(conn *websocket.Conn) {
	s.Lock()
	defer s.Unlock()

	for peer, streams := range s.m {
		for stream, conns := range streams {
			for i, c := range conns {
				if c == conn {
					s.m[peer][stream] = append(conns[:i], conns[i+1:]...)
					if len(s.m[peer][stream]) == 0 {
						delete(s.m[peer], stream)
					}
					break
				}
			}
		}
	}
}

// RemovePeer removes all stream answers for the given peer from the
// StreamAnswer map. It locks the map during the update to prevent concurrent
// map writes.
func (s *StreamAnswer) RemovePeer(peer string) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, peer)
}

// Get retrieves the stream answers for the given peer and stream name from the
// StreamAnswer map. It locks the map for reading during the lookup to prevent
// concurrent map access. The second return value indicates if a stream answer
// was found. This is an exported method that is part of the StreamAnswer API.
func (s *StreamAnswer) Get(peer, stream string) (conn []*websocket.Conn, ok bool) {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.m[peer]; !ok {
		return nil, false
	}
	if _, ok := s.m[peer][stream]; !ok {
		return nil, false
	}
	return s.m[peer][stream], true
}
