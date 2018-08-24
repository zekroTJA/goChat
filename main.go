package main

import (
	"os"
	"log"
	"time"
	"strings"
	"net/http"

	"github.com/boltdb/bolt"
)

func main() {

	args := os.Args[1:]
	port := "7777"
	if len(args) > 0 {
		port = args[0]
	}

	db, err := bolt.Open("goChatDb.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Database connection etablished")

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("accounts"))
		exists := bucket.Get([]byte("__exists"))
		if exists == nil {
			bucket, err = tx.CreateBucket([]byte("accounts"))
			if err != nil {
				log.Fatal("Error creating database bucket:", err)
			}
			err = bucket.Put([]byte("__exists"), []byte("1"))
			log.Println("'accounts' bucket created")
		}
		return nil
	})

	// Setting up new Chat instance
	chat := NewChat(db)

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
			ws.SetHandler("login", func(event *Event) {
				uname := event.Data.(map[string]string)["username"]
				// passwd := event.Data.(map[string]string)["password"]
				// for _, u := range chat.Sockets {
				// 	if u[0] == uname {
				// 		log.Println(u)
				// 		go func() {
				// 			ws.Out <- (&Event{
				// 				Name: "connect_rejected",
				// 				Data: nil,
				// 			}).Raw()
				// 		}()
				// 		return
				// 	}
				// }
				chat.Sockets[ws] = []string{ uname, UtilGetRandomColor() }
				chat.Broadcast((&Event{
					Name: "connected", 
					Data: map[string]interface{}{
						"name": chat.Sockets[ws][0],
						"nclients": len(chat.Sockets),
						"history": chat.History,
					},
				}).Raw())
			})

			// CHAT MESSAGE EVENT
			// -> Attach username to message
			// -> Broadcast the chat message to all users
			ws.SetHandler("message", func(event *Event) {
				username := chat.Sockets[ws][0]
				color := chat.Sockets[ws][1]
				if len(strings.Trim(event.Data.(string), " \t")) < 1 {
					return
				}
				event.Data = map[string]interface{}{
					"username": username,
					"timestamp": time.Now().Unix(),
					"color": color,
					"message": strings.Replace(event.Data.(string), "\\n", "\n", -1),
				}
				chat.Broadcast(event.Raw())
				chat.AppendHistory(event)
			})

			// DISCONNECT EVENT
			// -> Broadcast to all clients that 
			//    user has disconnected
			ws.SetHandler("disconnected", func(event *Event) {
				chat.Broadcast(event.Raw())
			})
		}
	})

	log.Printf("Listening on port %s ...", port)
	err = http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
