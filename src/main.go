package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type ChatStore struct {
	chat string
	mux  sync.Mutex
}

func (cs *ChatStore) SWrite(s string) {
	defer cs.mux.Unlock()
	cs.mux.Lock()
	cs.chat = cs.chat + s
}

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
				for c, _ := range chans {
					go func(c chan string) { c <- message }(c)
				}
			}
		}
	}()
	return register, broadcast
}

func BuildWSHandler(
	register chan chan string,
	broadcast chan string,
) func(http.ResponseWriter, *http.Request) {
	upgrader := websocket.Upgrader{}
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		// write routine will foward values from this chan
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
	fmt.Println("Nothing to see here, yet...")
}
