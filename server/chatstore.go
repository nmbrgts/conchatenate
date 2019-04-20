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
	sync.Mutex
}

// SWrite is a safe write method that appends the string it is given to
// the store.
func (cs *ChatStore) SWrite(id int, s string) {
	defer cs.Unlock()
	// static id for now, later SWrite will take an id as a unique
	// Web Socket Handler identifier
	cs.Lock()
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

// SRead is a safe read method that returns the current chat string.
// It should be safe to read without locks, but this will provide an
// interface for refactoring later
func (cs *ChatStore) SRead() string {
	cs.Lock()
	committedChat := cs.committedChat
	activeMsgs := cs.activeMsgs
	sep := cs.sep
	cs.Unlock()
	return committedChat + strings.Join(activeMsgs, sep)
}

// ShiftCursor is a concurrency safe method that moves the current sender ID
// "cursor" to a new entry at the end of the chat.
func (cs *ChatStore) ShiftCursor(id int) {
	defer cs.Unlock()
	cs.Lock()
	cs.indexMap[id] = len(cs.activeMsgs)
	cs.activeMsgs = append(cs.activeMsgs, "")
}

// StoreWorker is the process that handles writing to and broacasting
// the updated store value.
// TODO: Add process pooling
// TODO: Refector to use buffered channel?
// TODO: Refactor from single function to interface
func StoreWorker(store *ChatStore, send chan string) chan Message {
	receive := make(chan Message)
	go func() {
		for {
			msg := <-receive
			id, _ := msg.GetSenderId()
			content, gotContent := msg.GetContent()
			if gotContent {
				store.SWrite(id, content)
			} else {
				store.ShiftCursor(id)
			}
			send <- store.SRead()
		}
	}()
	return receive
}
