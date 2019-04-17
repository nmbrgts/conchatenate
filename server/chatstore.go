package main

import (
	"strings"
	"sync"
)

// ChatStore is mutex proted string and the main storage for the chat app.
// TODO: Flesh out an interface around ChatStore
// TODO: Move to external storage.
type ChatStore struct {
	committedChat string
	sep           string
	activeMsgs    []string
	indexMap      map[int]int
	mux           sync.Mutex
}

// SWrite is a safe write method that appends the string it is given to
// the store.
func (cs *ChatStore) SWrite(id int, s string) {
	defer cs.mux.Unlock()
	// static id for now, later SWrite will take an id as a unique
	// Web Socket Handler identifier
	cs.mux.Lock()
	if cs.indexMap == nil {
		cs.indexMap = make(map[int]int)
	}
	ix, ok := cs.indexMap[id]
	if !ok {
		cs.activeMsgs = append(cs.activeMsgs, s)
		cs.indexMap[id] = len(cs.activeMsgs) - 1
		return
	}
	cs.activeMsgs[ix] += s
}

// SRead is a safe read mothode that returns the current chat string.
// It should be safe to read without locks, but this will provide an
// interface for refactoring later
func (cs *ChatStore) SRead() string {
	cs.mux.Lock()
	committedChat := cs.committedChat
	activeMsgs := cs.activeMsgs
	sep := cs.sep
	cs.mux.Unlock()
	return committedChat + strings.Join(activeMsgs, sep)
}

// StoreWorker is the process that handles writing to and broacasting
// the updated store value.
// TODO: Add process pooling
// TODO: Refector to use buffered channel?
// TODO: Refactor from single function to interface
func StoreWorker(store *ChatStore, send chan string) chan string {
	receive := make(chan string)
	go func() {
		for {
			msg := <-receive
			store.SWrite(1, msg)
			send <- store.SRead()
		}
	}()
	return receive
}
