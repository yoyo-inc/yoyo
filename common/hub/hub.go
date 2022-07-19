package hub

import (
	"github.com/gorilla/websocket"
	"github.com/yoyo-inc/yoyo/common/logger"
	"net/http"
)

type Hub struct {
	upgrader    websocket.Upgrader
	connections []*websocket.Conn
	broadcast   chan interface{}
	register    chan *websocket.Conn
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.connections = append(h.connections, conn)
		case message := <-h.broadcast:
			for _, conn := range h.connections {
				err := conn.WriteJSON(message)
				if err != nil {
					return
				}
			}
		}
	}
}

var hub Hub

func Setup() {
	hub = Hub{
		upgrader: websocket.Upgrader{
			HandshakeTimeout: 60,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		connections: make([]*websocket.Conn, 0, 100),
		broadcast:   make(chan interface{}),
		register:    make(chan *websocket.Conn),
	}

	go hub.Run()
}

func BroadcastMessage(message interface{}) {
	hub.broadcast <- message
}

func Register(r *http.Request, w http.ResponseWriter) {
	conn, err := hub.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}

	defer conn.Close()

	hub.register <- conn
}
