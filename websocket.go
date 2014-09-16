package main

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

// Time allowed to write a message to the peer.
const writeWait = 3 * time.Second

type connection struct {
	ws *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

var conns = map[*connection]bool{}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func sendStatus(c *connection) {
	// Get mpd status.
	b, err := mpdStatus()
	if err != nil {
		return
	}
	// Send it.
	select {
	case c.send <- b:
	default:
		close(c.send)
		delete(conns, c)
	}
}

func broadcastStatus() {
	// Get mpd status.
	b, err := mpdStatus()
	if err != nil {
		glog.Errorln(err)
		return
	}
	// Broadcast it.
	for c := range conns {
		select {
		case c.send <- b:
		default:
			close(c.send)
			delete(conns, c)
		}
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *connection) writePump() {
	defer func() {
		if _, ok := conns[c]; ok {
			delete(conns, c)
		}
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			// Send mpd status to client.
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			glog.Error(err)
		}
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws}
	go c.writePump()
	conns[c] = true
	sendStatus(c)
}
