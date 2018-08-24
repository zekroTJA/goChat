package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	MAX_MSG_SIZE = 5000
)

// The Upgrader will be used to upgrade a HTTP conenction
// to a websocket connection with specified preferences.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket containing the websocket conenction,
// an instance pointer of Chat, which handles all
// websocket connections, an in channel for received
// messages and an out channel for sending messages.
// Also it contains a map containing all registered 
// event functions.
type WebSocket struct {
	Chat   *Chat
	Conn   *websocket.Conn
	Out    chan []byte
	In     chan []byte
	Events map[string]EventHandler
}

// NewWebSocket creates a new instance of WebSocket, upgrades the HTTP 
// connection to a websocket connection and runs the Reader and
// Writer go routines.
func NewWebSocket(chat *Chat, w http.ResponseWriter, r *http.Request) (*WebSocket, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ERROR | SOCKET CONNECT] %v", err)
		return nil, err
	}
	// conn.SetWriteDeadline(time.Now().Add(MSG_TIMEOUT))
	ws := &WebSocket{
		Chat: chat,
		Conn: conn, 
		Out: make(chan []byte), 
		In: make(chan []byte), 
		Events: make(map[string]EventHandler),
	}
	go ws.Reader()
	go ws.Writer()
	return ws, nil
}

// Reader routine will wait for incomming messages.
// They will be tried to be parsed to an Event object.
// Then, the Events map will be checked if an action
// is registered for the Event, which will be executed
// with the events data as argument.
func (ws *WebSocket) Reader() {
	defer func() {
		ws.Chat.Unregister(ws)
		ws.Conn.Close()
	}()
	ws.Conn.SetReadLimit(MAX_MSG_SIZE)
	for {
		_, message, err := ws.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[ERROR] %v", err)
			}
			break
		}
		event, err := NewEventFromRaw(message)
		if err != nil {
			log.Printf("[ERROR | MSG] %v", err)
		} else {
			log.Printf("[MSG] %v", event)
		}
		if action, ok := ws.Events[event.Name]; ok {
			action(event)
		}
	}
}

// Writer routines waits for new Events in 
// Out channel to be send to the conencted client
// as binary data.
func (ws *WebSocket) Writer() {
	for {
		select {
		case message, ok := <-ws.Out:
			if !ok {
				ws.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := ws.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			w.Close()
		}
	}
}

// SetHandler registers a new Event by name of the
// event and a function, which will be executed if 
// the event fires, which will be get passed the data
// of the event.
func (ws *WebSocket) SetHandler(event string, action EventHandler) *WebSocket {
	ws.Events[event] = action
	return ws
}