package main

import "log"

// Chat collects and handles all 
// websocket conenctions
type Chat struct {
	Sockets map[*WebSocket]string
}

// NewChat creates a new instance pointer of Chat
func NewChat() *Chat {
	chat := &Chat{Sockets: make(map[*WebSocket]string)}
	return chat
}

// Register registers a new websocket connection
func (c *Chat) Register(socket *WebSocket) {
	log.Printf("[SOCKET CONNECTED]")
	c.Sockets[socket] = ""
}

// Unregister unregisters a disconnected client
// and closes clients channels
func (c *Chat) Unregister(socket *WebSocket, conerr ...bool) {
	log.Printf("[SOCKET DISCONNECTED]")
	if action, ok := socket.Events["disconnected"]; ok && len(conerr) == 0 {
		action(&Event{
			Name: "disconnected", 
			Data: map[string]interface{}{
				"name": c.Sockets[socket],
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