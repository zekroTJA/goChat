package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocket struct {
	Conn *websocket.Conn
	Send chan []byte
}

func (ws *WebSocket) Reader() {
	for {
		_, message, err := ws.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[ERROR] %v", err)
			}
			break
		}
		log.Printf("[MSG] %s", string(message))
	}
}

func (ws *WebSocket) Writer() {
	for {
		select {
		case message, ok := <-ws.Send:
			if !ok {
				ws.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
		}
	}
}

func NewWebSocket(w http.ResponseWriter, r *http.Request) *WebSocket {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	ws := &WebSocket{Conn: conn, Send: make(chan []byte)}
	go ws.Reader()
	return ws
}
