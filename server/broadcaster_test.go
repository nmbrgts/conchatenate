package main

import (
	"testing"
	"time"
)

func TestBroadcaster(t *testing.T) {
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
			for _, channel := range chans {
				register <- channel
			}
			broadcast <- want
			gots := make(chan string)
			for _, channel := range chans {
				go func(c chan string) { gots <- <-c }(channel)
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
			channel := make(chan string)
			register <- channel
			register <- channel
			broadcast <- want
			respCount := 0
		L:
			for {
				select {
				case <-channel:
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
