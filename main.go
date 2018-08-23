package main

import (
	"log"
	"net/http"
)

func main() {

	chat := NewChat()

	http.Handle("/", http.FileServer(http.Dir("./assets")))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if ws, err := NewWebSocket(chat, w, r); err == nil {
			chat.Register(ws)

			// Username Input handler
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

			ws.SetHandler("message", func(event *Event) {
				username := chat.Sockets[ws]
				event.Data = username + ": " + event.Data.(string)
				chat.Broadcast(event.Raw())
			})

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
