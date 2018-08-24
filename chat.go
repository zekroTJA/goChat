package main

import (
	"log"

	"github.com/boltdb/bolt"
)

const (
	HISTORY_CAP = 200
)

// Chat collects and handles all 
// websocket conenctions
type Chat struct {
	Sockets map[*WebSocket][]string
	History []*Event
	DB      *bolt.DB
}

// NewChat creates a new instance pointer of Chat
func NewChat(DB *bolt.DB) *Chat {
	chat := &Chat{Sockets: make(map[*WebSocket][]string), DB: DB}
	return chat
}

// Register registers a new websocket connection
func (c *Chat) Register(socket *WebSocket) {
	log.Printf("[SOCKET CONNECTED]")
	c.Sockets[socket] = make([]string, 2)
}

// Unregister unregisters a disconnected client
// and closes clients channels
func (c *Chat) Unregister(socket *WebSocket, conerr ...bool) {
	log.Printf("[SOCKET DISCONNECTED]")
	if action, ok := socket.Events["disconnected"]; ok && len(conerr) == 0 {
		action(&Event{
			Name: "disconnected", 
			Data: map[string]interface{}{
				"name": c.Sockets[socket][0],
				"nclients": len(c.Sockets),
			},
		})
	}
	// close(socket.Out)
	// close(socket.In)
	delete(c.Sockets, socket)
	socket.Conn.Close()
}

// Broadcast sends an event to all
// connected clients
func (c *Chat) Broadcast(message []byte) {
	for s, _ := range c.Sockets {
		select {
		case s.Out <- message:
		default:
			c.Unregister(s, true)
		}
	}
}

func (c *Chat) AppendHistory(event *Event) {
	if len(c.History) > HISTORY_CAP {
		c.History = append(c.History[len(c.History)-HISTORY_CAP:], event)
		return
	}
	c.History = append(c.History, event)
}