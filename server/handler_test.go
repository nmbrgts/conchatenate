package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWebSocket(t *testing.T) {
	buildTestServer := func(t *testing.T) (*websocket.Conn, chan chan string, chan Message) {
		t.Helper()
		register := make(chan chan string)
		broadcast := make(chan Message)
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
			case msg := <-broadcast:
				got, _ = msg.GetContent()
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
