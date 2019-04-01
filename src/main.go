package main

import (
	"fmt"
	"sync"
)

type ChatStore struct {
	chat string
	mux sync.Mutex
}

func (cs *ChatStore) SWrite (s string) {
	defer cs.mux.Unlock()
	cs.mux.Lock()
	cs.chat = cs.chat + s
}

func main() {
	fmt.Println("nothing to see here, yet...")
}
