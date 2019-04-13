// Conchatenate: soon to be the world's worst overly concurrent chat app
package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"strings"
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
	cs.mux.Unlock()
	return committedChat + strings.Join(activeMsgs, ",")
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

// Broadcaster is the process that handles sending outgoing messages to all
// websocket processes.
// TODO: Interface?
func Broadcaster() (chan chan string, chan string) {
	register := make(chan chan string)
	broadcast := make(chan string)
	go func() {
		chans := make(map[chan string]bool)
		for {
			select {
			case reg := <-register:
				// register new channel
				chans[reg] = true
			case message := <-broadcast:
				// broadcast to registered channels w/o blocking
				for c := range chans {
					go func(c chan string) { c <- message }(c)
				}
			}
		}
	}()
	return register, broadcast
}

// BuildWSHandler builds our websocket handler with provided send and receive channels.
// TODO: Interface?
func BuildWSHandler(
	register chan chan string,
	broadcast chan string,
) func(http.ResponseWriter, *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// write routine will forward values from this chan
		receiver := make(chan string)
		register <- receiver
		// write routine
		go func() {
			for {
				message := <-receiver
				err := conn.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					break
				}
			}
		}()
		// read routine
		go func() {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					break
				}
				broadcast <- string(message)
			}
		}()
	}
}

func main() {
	store := ChatStore{}
	register, broadcast := Broadcaster()
	receive := StoreWorker(&store, broadcast)
	http.HandleFunc("/chat", BuildWSHandler(register, receive))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/static/testpage.html")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
