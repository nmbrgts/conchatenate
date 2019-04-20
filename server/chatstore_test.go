package main

import (
	"strings"
	"testing"
)

func TestChatStore(t *testing.T) {
	t.Run(
		"ChatStore should return an emtpy list when initialized",
		func(t *testing.T) {
			want := ""
			store := ChatStore{}
			got := store.SRead()
			if got != want {
				t.Errorf(
					"Expected store.SRead() to return \"%s\", but got \"%s\" instead",
					want, got)
			}
		},
	)
	t.Run(
		"ChatStore should be read/write invariant with one writer",
		func(t *testing.T) {
			want := "hallo, this is dog"
			store := ChatStore{}
			store.SWrite(1, want)
			got := store.SRead()
			if got != want {
				t.Errorf(
					"Expected store.SRead() to return \"%s\", but got \"%s\" instead",
					want, got)
			}
		},
	)
	t.Run(
		"Length should equal total length of all concurrent SWrite values",
		func(t *testing.T) {
			chats := ChatStore{}
			want := 1000
			done := make(chan bool)
			for i := 0; i < want; i++ {
				go func(done chan bool) {
					chats.SWrite(1, "a")
					done <- true
				}(done)
			}
			for i := 0; i < want; i++ {
				<-done
			}
			got := len(chats.SRead())
			if want != got {
				t.Errorf("Expected %d elements, got %d", want, got)
			}
		})
	t.Run(
		"Ordering of writes should be maintained between registered write IDs",
		func(t *testing.T) {
			chars := []string{"A", "B", "C"}
			store := ChatStore{}
			want := ""
			for id, char := range chars {
				want += strings.Repeat(char, 101)
				store.SWrite(id, char)
			}
			done := make(chan bool, 300)
			for i := 0; i < 300; i++ {
				go func(id int, done chan bool) {
					store.SWrite(id, chars[id])
					done <- true
				}(i%3, done)
			}
			for i := 0; i < 300; i++ {
				<-done
			}
			got := store.SRead()
			if want != got {
				t.Errorf("Expected ordering of elements to match between iteratively built strin:\nwant:\n\"%s\"\ngot:\n\"%s\"",
					want, got)
			}
		},
	)
	t.Run(
		"ChatStore injects sep member between entries",
		func(t *testing.T) {
			msg := "hallo, this is dog"
			sep := "<||>"
			want := msg + sep + msg
			store := ChatStore{sep: sep}
			store.SWrite(0, msg)
			store.SWrite(1, msg)
			got := store.SRead()
			if want != got {
				t.Errorf("Expect ChatSore to read out \"%s\", but got \"%s\"",
					want, got)
			}
		},
	)
	t.Run(
		"`ShiftCursor` method moves sender id to a new index in the chat",
		func(t *testing.T) {
			want := "a.a.a.a.a.a.a.a.a.a"
			senderId := 1
			msg := "a"
			store := ChatStore{sep: "."}
			for i := 0; i < 9 ; i ++ {
				store.SWrite(senderId, msg)
				store.ShiftCursor(senderId)
			}
			store.SWrite(senderId, msg)
			got := store.SRead()
			if want != got {
				t.Errorf("Expect ChatSore to read out \"%s\", but got \"%s\"",
					want, got)
			}
		},
	)
}

func TestStoreWorker(t *testing.T) {
	t.Run(
		"StoreWorker should update store from messages sent to input channel",
		func(t *testing.T) {
			want := "hallo, this is dog"
			store := ChatStore{}
			broadcast := make(chan string)
			receive := StoreWorker(&store, broadcast)
			receive <- PlainMessage{want}
			got := store.SRead()
			if got != want {
				t.Errorf("Expected store value to be \"%s\", instead got \"%s\"", want, got)
			}
		},
	)
	t.Run(
		"StoreWorker should breoadcast store's chat value",
		func(t *testing.T) {
			want := "hallo, this is dog"
			store := ChatStore{}
			broadcast := make(chan string)
			receive := StoreWorker(&store, broadcast)
			receive <- PlainMessage{want}
			got := <-broadcast
			want = store.SRead()
			if got != want {
				t.Errorf("Expected broadcasted value to be \"%s\", instead got \"%s\"", want, got)
			}
		},
	)
	t.Run(
		"StoreWorker progressex sender's cursor in response to NilMessage",
		func (t *testing.T) {
			msg := "hallo, this is dog"
			sep := "->0<-"
			want := msg + sep + msg

		},
	)
}
