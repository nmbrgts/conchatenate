// contains code for handling broadcasting events
package main

// Broadcaster is the process that handles sending outgoing messages to all
// websocket processes.
// TODO: Interface?
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
				for c := range chans {
					go func(c chan string) { c <- message }(c)
				}
			}
		}
	}()
	return register, broadcast
}
