package main

import (
	"log"
)

const (
	HISTORY_CAP = 200
)

// Chat collects and handles all
// websocket conenctions
type Chat struct {
	Sockets map[*WebSocket]*Author
	Users   map[string]string
	History []*Event
	AccMgr  *AccountManager
}

type Author struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Color    string `json:"color"`
}

type Message struct {
	ID        int64   `json:"id"`
	Content   string  `json:"content"`
	Author    *Author `json:"author"`
	Timestamp int64   `json:"timestamp"`
}

// NewChat creates a new instance pointer of Chat
func NewChat(accMgr *AccountManager) *Chat {
	chat := &Chat{
		Sockets: make(map[*WebSocket]*Author),
		Users:   make(map[string]string),
		AccMgr:  accMgr,
	}
	return chat
}

// Register registers a new websocket connection
func (c *Chat) Register(socket *WebSocket) {
	log.Printf("[SOCKET CONNECTED]")
	c.Sockets[socket] = &Author{}
}

// Unregister unregisters a disconnected client
// and closes clients channels
func (c *Chat) Unregister(socket *WebSocket, conerr ...bool) {
	log.Printf("[SOCKET DISCONNECTED]")
	if action, ok := socket.Events["disconnected"]; ok && len(conerr) == 0 {
		action(&Event{
			Name: "disconnected",
			Data: map[string]interface{}{
				"name":     c.Sockets[socket].Username,
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
