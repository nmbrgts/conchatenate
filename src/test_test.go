package main

import (
	"testing"
)


func TestChatStore(t *testing.T) {
	t.Run(
		"Length should equal total length of all concurrent SWrite values",
		func (t *testing.T) {
			chats := ChatStore{}
			want := 1000
			done := make(chan bool)
			for i := 0; i < want; i++ {
				go func (done chan bool) {
					chats.SWrite("a")
					done <- true
				}(done)
			}
			for i := 0; i < want; i ++ {
				<- done
			}
			got := len(chats.chat)
			if want != got {
				t.Errorf("Expected %d elements, got %d", want, got)
			}
	})
}
