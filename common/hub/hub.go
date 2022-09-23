package hub

import (
	"github.com/gorilla/websocket"
	"github.com/yoyo-inc/yoyo/common/logger"
	"net/http"
	"time"
)

var hub Hub

type Hub struct {
	upgrader  websocket.Upgrader
	clients   map[*Client]interface{}
	broadcast chan string
	register  chan *Client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					if err := client.conn.Close(); err != nil {
						logger.Error(err)
					}
					delete(h.clients, client)
				}
			}
		}
	}
}

type Client struct {
	conn *websocket.Conn
	hub  *Hub
	send chan string
}

func (c *Client) WriteDump() {
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {
			return
		}
		delete(c.hub.clients, c)
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// The hub closed the channel.
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					logger.Error(err)
					return
				}
				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					logger.Error(err)
					return
				}

				w.Write([]byte(message))

				n := len(c.send)
				for i := 0; i < n; i++ {
					w.Write([]byte{'\n'})
					w.Write([]byte(<-c.send))
				}

				if err = w.Close(); err != nil {
					return
				}
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Error(err)
				return
			}
		}
	}
}

func Setup() {
	hub = Hub{
		upgrader: websocket.Upgrader{
			HandshakeTimeout: 60,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients:   make(map[*Client]interface{}),
		broadcast: make(chan string),
		register:  make(chan *Client),
	}

	go hub.Run()
}

func SendMessage(message string) {
	hub.broadcast <- message
}

func Register(r *http.Request, w http.ResponseWriter) {
	conn, err := hub.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	client := &Client{
		conn: conn,
		hub:  &hub,
		send: make(chan string, 10),
	}

	hub.register <- client

	go client.WriteDump()
}
