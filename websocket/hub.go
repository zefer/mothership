package websocket

// Maintains active WebSocket connections and broadcasts messages
type hub struct {
	// Registered connections.
	connections map[*Connection]bool

	// Inbound messages from the connections.
	Broadcast chan []byte

	// Register requests from the connections.
	register chan *Connection

	// Unregister requests from connections.
	unregister chan *Connection

	listeners []func(*Connection)
}

var Hub = hub{
	Broadcast:   make(chan []byte),
	register:    make(chan *Connection),
	unregister:  make(chan *Connection),
	connections: make(map[*Connection]bool),
	listeners:   []func(*Connection){},
}

func (h *hub) OnConnect(l func(*Connection)) {
	h.listeners = append(h.listeners, l)
}

func (h *hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			for _, l := range h.listeners {
				l(c)
			}
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.Broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}
