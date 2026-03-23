package hub

type Client interface {
	Send(msg []byte) bool
	Close()
}

type Hub struct {
	clients    map[Client]bool
	broadcast  chan []byte
	register   chan Client
	unregister chan Client
}

func New() *Hub {
	return &Hub{
		clients:    make(map[Client]bool),
		register:   make(chan Client),
		unregister: make(chan Client),
		broadcast:  make(chan []byte, 256),
	}
}

func (h *Hub) Register(c Client)    { h.register <- c }
func (h *Hub) Unregister(c Client)  { h.unregister <- c }
func (h *Hub) Broadcast(msg []byte) { h.broadcast <- msg }

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true

		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				c.Close()
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				if !c.Send(msg) {
					delete(h.clients, c)
					c.Close()
				}
			}
		}
	}
}
