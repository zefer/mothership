package websocket

import (
	"time"

	websock "github.com/gorilla/websocket"
	"gopkg.in/airbrake/glog.v1"
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

type Conn struct {
	ws *websock.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

var (
	upgrader = websock.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func (c *Conn) Send(b []byte) {
	select {
	case c.send <- b:
	default:
		h.unregister <- c
		c.ws.Close()
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Conn) readPump() {
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
func (c *Conn) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *Conn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websock.CloseMessage, []byte{})
				return
			}
			// Send string to client.
			if err := c.write(websock.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			glog.Infof("Sending ping")
			if err := c.write(websock.PingMessage, []byte{}); err != nil {
				glog.Errorf("Ping error %v: ", err)
				return
			}
		}
	}
}
