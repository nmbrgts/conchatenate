package main

import (
	"fmt"
	"sync"
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
		var chans []chan string
		for {
			select {
			case reg := <-register:
				// register new channel
				chans = append(chans, reg)
			case message := <-broadcast:
				// broadcast to registered channels w/o blocking
				for _, c := range chans {
					go func(c chan string) { c <- message }(c)
				}
			}
		}
	}()
	return register, broadcast
}

func main() {
	fmt.Println("Nothing to see here, yet...")
}
