package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
)

func main() {

	args := os.Args[1:]
	port := "7777"
	if len(args) > 0 {
		port = args[0]
	}

	accMgr, err := NewAccountManager("./accDataBase.json")
	if err != nil {
		panic(err)
	}
	defer accMgr.Save()

	// Setting up new Chat instance
	chat := NewChat(accMgr)

	authorNode, _ := snowflake.NewNode(0)
	messageNode, _ := snowflake.NewNode(1)

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
				dataMap := event.Data.(map[string]interface{})
				uname := dataMap["username"].(string)
				passwd := dataMap["password"].(string)
				acc, valid := accMgr.Check(uname, passwd)
				if acc == nil {
					accMgr.Register(uname, passwd)
					accMgr.Save()
				} else if !valid {
					go func() {
						ws.Out <- (&Event{
							Name: "connect_rejected",
							Data: nil,
						}).Raw()
					}()
					return
				}
				chat.Sockets[ws] = &Author{
					ID:       authorNode.Generate().Int64(),
					Color:    UtilGetRandomColor(),
					Username: uname,
				}
				chat.Users[chat.Sockets[ws].Username] = chat.Sockets[ws].Color
				chat.Broadcast((&Event{
					Name: "connected",
					Data: map[string]interface{}{
						"name":     chat.Sockets[ws].Username,
						"nclients": len(chat.Sockets),
						"clients":  chat.Users,
						"history":  chat.History,
					},
				}).Raw())
			})

			ws.SetHandler("checkUsername", func(event *Event) {
				uname := event.Data.(string)
				accExists := accMgr.Get(uname) != nil
				log.Println(accExists)
				go func() {
					ws.Out <- (&Event{
						Name: "usernameState",
						Data: accExists,
					}).Raw()
				}()
			})

			// CHAT MESSAGE EVENT
			// -> Attach username to message
			// -> Broadcast the chat message to all users
			ws.SetHandler("message", func(event *Event) {
				// username := chat.Sockets[ws].Username
				// color := chat.Sockets[ws].Color
				if len(strings.Trim(event.Data.(string), " \t")) < 1 {
					return
				}
				event.Data = &Message{
					Author:    chat.Sockets[ws],
					Content:   strings.Replace(event.Data.(string), "\\n", "\n", -1),
					Timestamp: time.Now().Unix(),
					ID:        messageNode.Generate().Int64(),
				}
				chat.Broadcast(event.Raw())
				chat.AppendHistory(event)
			})

			ws.SetHandler("deleteMessage", func(event *Event) {
				msgID := int64(event.Data.(float64))
				msg, i := chat.GetMessageByID(msgID)
				if msg == nil {
					return
				}
				chat.History = append(chat.History[:i], chat.History[i+1:]...)
				eventOut := &Event{
					Name: "messageDeleted",
					Data: msg,
				}
				chat.Broadcast(eventOut.Raw())
			})

			// DISCONNECT EVENT
			// -> Broadcast to all clients that
			//    user has disconnected
			ws.SetHandler("disconnected", func(event *Event) {
				dataMap := event.Data.(map[string]interface{})
				uname := dataMap["name"].(string)
				delete(chat.Users, uname)
				dataMap["clients"] = chat.Users
				event.Data = dataMap
				chat.Broadcast(event.Raw())
			})
		}
	})

	log.Printf("Listening on port %s ...", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
