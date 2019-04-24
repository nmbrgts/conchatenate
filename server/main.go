// Conchatenate: soon to be the world's worst overly concurrent chat app
package main

import (
	"log"
	"net/http"
	"sync"
)

// Registrar is a shared data structure for IDing unique handler processes.
type Registrar struct {
	count int
	sync.Mutex
}

// Register returns a unique int to the process that calls it.
// This value is used internally as an identifier for websocket handlers
func (r *Registrar) Register() int {
	r.Lock()
	id := r.count
	r.count++
	r.Unlock()
	return id
}

func main() {
	store := ChatStore{}
	register, broadcast := Broadcaster()
	receive := StoreWorker(&store, broadcast)
	http.HandleFunc("/chat", BuildWSHandler(register, receive))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/testpage.html")
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
