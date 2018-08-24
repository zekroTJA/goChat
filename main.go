package main

import (
	"log"
	"net/http"
)

func main() {

	// Setting up new Chat instance
	chat := NewChat()

	// Delivering client (website) content to root address.
	http.Handle("/", http.FileServer(http.Dir("./assets")))

	// Setting up websocket conenction path to /ws.
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if ws, err := NewWebSocket(chat, w, r); err == nil {
			chat.Register(ws)

			// USERNAME INPUT EVENT
			// -> Checks if name is not connected yet
			//    -> else send 'connect_reject' event
			// -> Broadcast to all clients that user has connected
			ws.SetHandler("username", func(event *Event) {
				uname := event.Data.(string)
				for _, u := range chat.Sockets {
					if u == uname {
						go func() {
							ws.Out <- (&Event{
								Name: "connect_rejected",
								Data: nil,
							}).Raw()
						}()
						return
					}
				}
				chat.Sockets[ws] = uname
				chat.Broadcast((&Event{
					Name: "connected", 
					Data: map[string]interface{}{
						"name": chat.Sockets[ws],
						"nclients": len(chat.Sockets),
					},
				}).Raw())
			})

			// CHAT MESSAGE EVENT
			// -> Attach username to message
			// -> Broadcast the chat message to all users
			ws.SetHandler("message", func(event *Event) {
				username := chat.Sockets[ws]
				event.Data = username + ": " + event.Data.(string)
				chat.Broadcast(event.Raw())
			})

			// DISCONNECT EVENT
			// -> Broadcast to all clients that 
			//    user has disconnected
			ws.SetHandler("disconnected", func(event *Event) {
				chat.Broadcast(event.Raw())
			})
		}
	})

	log.Println("Listening...")
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal(err)
	}
}
