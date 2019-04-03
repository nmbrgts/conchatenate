package main

import (
	"testing"
	"time"
)

func TestChatStore(t *testing.T) {
	t.Run(
		"Length should equal total length of all concurrent SWrite values",
		func(t *testing.T) {
			chats := ChatStore{}
			want := 1000
			done := make(chan bool)
			for i := 0; i < want; i++ {
				go func(done chan bool) {
					chats.SWrite("a")
					done <- true
				}(done)
			}
			for i := 0; i < want; i++ {
				<-done
			}
			got := len(chats.chat)
			if want != got {
				t.Errorf("Expected %d elements, got %d", want, got)
			}
		})
}

func TestBroadcast(t *testing.T) {
	t.Run(
		"All channels registered with broadcast should receive broadcasted messages",
		func(t *testing.T) {
			want := "hallo, this is dog"
			register, broadcast := Broadcaster()
			// register, _ := Broadcaster()
			chans := []chan string{
				make(chan string),
				make(chan string),
				make(chan string),
			}
			for _, chan_ := range chans {
				register <- chan_
			}
			broadcast <- want
			gots := make(chan string)
			for _, chan_ := range chans {
				go func(c chan string) { gots <- <-c }(chan_)
			}
			for i := 0; i < len(chans); i++ {
				got := <-gots
				if want != got {
					t.Errorf(
						"Expected all channels to return \"%s\" but, channel %d returned \"%s\"",
						want, i, got,
					)
				}
			}
		})
	t.Run(
		"Gracefully handle multiple registrations of the same channel",
		func(t *testing.T) {
			want := "hallo, this is dog"
			register, broadcast := Broadcaster()
			chan_ := make(chan string)
			register <- chan_
			register <- chan_
			broadcast <- want
			respCount := 0
		L:
			for {
				select {
				case <-chan_:
					respCount++
				case <-time.After(2 * time.Second):
					break L
				}
			}
			if respCount > 1 {
				t.Errorf(
					"Expected doubly registered channel to receive on broadcast, but it received %d broadcasts",
					respCount,
				)
			}
		},
	)
}
