package main

import (
	"log"
	"time"

	"github.com/bwmarrin/snowflake"
)

const (
	HISTORY_CAP = 200
)

// Chat collects and handles all
// websocket conenctions
type Chat struct {
	Sockets     map[*WebSocket]*Author
	Sessions    map[string]*Session
	Users       map[string]string
	History     []*Event
	TempHistory map[int64]map[int64]bool
	AccMgr      *AccountManager
}

type Session struct {
	ID     int64   `json:"id"`
	Author *Author `json:"author"`
	IPAddr string  `json:"ipaddr"`
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
		Sockets:     make(map[*WebSocket]*Author),
		Sessions:    make(map[string]*Session),
		Users:       make(map[string]string),
		AccMgr:      accMgr,
		TempHistory: make(map[int64]map[int64]bool),
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
	if action, ok := socket.Events["disconnected"]; ok && len(conerr) == 0 && c.Sockets[socket] != nil {
		action(&Event{
			Name: "disconnected",
			Data: map[string]interface{}{
				"name":     c.Sockets[socket].Username,
				"nclients": len(c.Sockets),
			},
		})
	}
	delete(c.Sockets, socket)
	socket.Conn.Close()
}

func (c *Chat) Login(author *Author, addr string) {
	node, _ := snowflake.NewNode(100)
	id := node.Generate().Int64()
	session := &Session{
		Author: author,
		ID:     id,
		IPAddr: addr,
	}
	c.Sessions[addr] = session
	time.AfterFunc(5*time.Minute, func() {
		delete(c.Sessions, addr)
	})
}

// Broadcast sends an event to all
// connected clients
func (c *Chat) Broadcast(message []byte) {
	for s := range c.Sockets {
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

func (c *Chat) FetchMessage(ident func(*Event) bool) (*Event, int) {
	for i, msg := range c.History {
		if ident(msg) {
			return msg, i
		}
	}
	return nil, -1
}

func (c *Chat) GetMessageByID(id int64) (*Event, int) {
	return c.FetchMessage(func(e *Event) bool {
		return (e.Data.(*Message).ID == id)
	})
}

func (c *Chat) DeleteMessageByID(id int64) *Event {
	msg, i := c.GetMessageByID(id)
	if msg == nil {
		return nil
	}
	c.History = append(c.History[:i], c.History[i+1:]...)
	return msg
}

func (c *Chat) EnqueueTempHistory(id int64) {
	now := time.Now().UnixNano()
	if len(c.TempHistory[id]) == 0 {
		c.TempHistory[id] = map[int64]bool{
			now: true,
		}
	} else {
		c.TempHistory[id][now] = true
	}
	time.AfterFunc(10*time.Second, func() {
		delete(c.TempHistory[id], now)
	})
}

func (c *Chat) TempHistoryLength(id int64) int {
	if c.TempHistory[id] == nil {
		return 0
	}
	return len(c.TempHistory[id])
}
