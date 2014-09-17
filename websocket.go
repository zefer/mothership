package main

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

// Time allowed to write a message to the peer.
const (
	// Time allowed to write a message to the peer.
	writeWait = 3 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type connection struct {
	ws *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

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
		h.unregister <- c
		c.ws.Close()
	}
}

func broadcastStatus() {
	// Get mpd status.
	b, err := mpdStatus()
	if err != nil {
		glog.Errorln(err)
		return
	}
	h.broadcast <- b
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		glog.Infof("Received pong")
		return nil
	})
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		h.broadcast <- message
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
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
		case <-ticker.C:
			glog.Infof("Sending ping")
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				glog.Errorf("Ping error %v: ", err)
				return
			}
		}
	}
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Error(err)
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writePump()
	sendStatus(c)
	c.readPump()
}
