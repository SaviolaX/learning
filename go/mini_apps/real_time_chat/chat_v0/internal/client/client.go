package client

import (
	"chatV0/internal/hub"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = time.Second * 10
	pongWait       = time.Second * 60
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
	name string
}

func NewClient(conn *websocket.Conn, h *hub.Hub, name string) *Client {
	c := &Client{
		conn: conn,
		name: name,
		send: make(chan []byte, 256),
	}

	h.Register(c)

	go c.readPump(h)
	go c.writePump()

	return c
}

func (c *Client) Send(msg []byte) bool {
	select {
	case c.send <- msg:
		return true
	default:
		return false
	}
}

func (c *Client) Close() {
	close(c.send)
}

func (c *Client) readPump(h *hub.Hub) {
	defer func() {
		h.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("unexpected error from [%s]: %v", c.name, err)
			}
			return
		}

		message = []byte(c.name + ": " + string(message))

		h.Broadcast(message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			n := len(c.send)
			for range n {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
