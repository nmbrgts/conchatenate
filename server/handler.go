package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// we'll use this later to register an ID for the handler
var registrar = Registrar{}

// BuildWSHandler builds our websocket handler with provided send and receive channels.
// TODO: Interface?
func BuildWSHandler(
	register chan chan string,
	broadcast chan Message,
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
		id := registrar.Register()
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
				broadcast <- IdMessage{id, string(message)}
			}
		}()
	}
}
