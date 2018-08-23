package main

import "log"

type Chat struct {
	Sockets map[*WebSocket]string
}

func NewChat() *Chat {
	chat := &Chat{Sockets: make(map[*WebSocket]string)}
	return chat
}

func (c *Chat) Register(socket *WebSocket) {
	log.Printf("[SOCKET CONNECTED]")
	c.Sockets[socket] = ""
}

func (c *Chat) Unregister(socket *WebSocket) {
	log.Printf("[SOCKET DISCONNECTED]")
	if action, ok := socket.Events["disconnected"]; ok {
		action(&Event{
			Name: "disconnected", 
			Data: map[string]interface{}{
				"name": c.Sockets[socket],
				"nclients": len(c.Sockets),
			},
		})
	}
	close(socket.Out)
	close(socket.In)
	delete(c.Sockets, socket)
}

func (c *Chat) Broadcast(message []byte) {
	for s, _ := range c.Sockets {
		select {
		case s.Out <- message:
		default:
			c.Unregister(s)
		}
	}
}