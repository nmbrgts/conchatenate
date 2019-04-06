package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
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

func TestWebSocket(t *testing.T) {
	buildTestServer := func(t *testing.T) (*websocket.Conn, chan chan string, chan string) {
		t.Helper()
		register := make(chan chan string)
		broadcast := make(chan string)
		testServer := httptest.NewServer(http.HandlerFunc(BuildWSHandler(register, broadcast)))
		wsUrl := "ws" + strings.TrimPrefix(testServer.URL, "http")
		ws, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
		if err != nil {
			t.Fatal(err)
		}
		return ws, register, broadcast
	}
	t.Run(
		"Websocket handler should register a channel that it writes content from",
		func(t *testing.T) {
			want := "hallo, this is dog"
			ws, register, _ := buildTestServer(t)
			channel := <-register
			go func() { channel <- want }()
			_, got, err := ws.ReadMessage()
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != want {
				t.Errorf(
					"Expected WS to send \"%s\" instead, got \"%s\"",
					want, got,
				)
			}
		},
	)
	t.Run(
		"Websocket handler should forward messages to the broadcast chanel",
		func(t *testing.T) {
			want := "hallo, this is dog"
			ws, register, broadcast := buildTestServer(t)
			<-register
			err := ws.WriteMessage(websocket.TextMessage, []byte(want))
			if err != nil {
				t.Fatal(err)
			}
			got := "hallo, this is not dog"
			select {
			case got = <-broadcast:
			case <-time.After(2 * time.Second):
				t.Errorf("Expected WS handler to broadcast, but it never did")
				return
			}
			if string(got) != want {
				t.Errorf("Expected WS handler to broadcast \"%s\" instead, got \"%s\"", want, got)
			}
		},
	)
}

func TestStoreWorker(t *testing.T) {
	t.Run(
		"StoreWorker should update store based in channel",
		func(t *testing.T) {
			want := "hallo, this is dog"
			store := ChatStore{}
			broadcast := make(chan string)
			recieve := StoreWorker(&store, broadcast)
			recieve <- want
			got := store.chat
			if got != want {
				t.Errorf("Expected WS handler to broadcast \"%s\" instead, got \"%s\"", want, got)
			}
		},
	)
	t.Run(
		"StoreWorker should breoadcast store's chat value",
		func(t *testing.T) {
			want := "hallo, this is dog"
			store := ChatStore{}
			broadcast := make(chan string)
			recieve := StoreWorker(&store, broadcast)
			recieve <- want
			got := <-broadcast
			want = store.chat
			if got != want {
				t.Errorf("Expected store value \"%s\" to mat broadcast value \"%s\"", want, got)
			}
		},
	)
}
