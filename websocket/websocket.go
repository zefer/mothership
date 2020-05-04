package websocket

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func init() {
	go h.run()
}

// Sends the message to all active websocket connections.
func Broadcast(m []byte) {
	h.broadcast <- m
}

// Registers a listener func that will be called each time we receive a new
// websocket connection.
func OnConnect(l func(*Conn)) {
	h.listeners = append(h.listeners, l)
}

// Provides a websocket upgrade handler which accepts a browser websocket
// connection and registers them in the hub of active connections.
func Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	c := &Conn{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writePump()
	c.readPump()
}
