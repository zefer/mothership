package websocket

// Maintains active WebSocket connections and broadcasts messages
type hub struct {
	// Registered connections.
	connections map[*Conn]bool

	// Inbound messages from the connections.
	Broadcast chan []byte

	// Register requests from the connections.
	register chan *Conn

	// Unregister requests from connections.
	unregister chan *Conn

	listeners []func(*Conn)
}

var Hub = hub{
	Broadcast:   make(chan []byte),
	register:    make(chan *Conn),
	unregister:  make(chan *Conn),
	connections: make(map[*Conn]bool),
	listeners:   []func(*Conn){},
}

func (h *hub) OnConnect(l func(*Conn)) {
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
