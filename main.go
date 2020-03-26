package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

type ChatText struct {
	Identify string `json:"identify"`
	RoomName string `json:"roomName"`
	UserName string `json:"userName"`
	Text     string `json:"text"`
}

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "join-room", func(s socketio.Conn, msg string) {
		chatText := ChatText{}
		err := json.Unmarshal([]byte(msg), &chatText)
		if err != nil {
			log.Fatal(err)
		}
		server.JoinRoom("/", chatText.RoomName, s)
		server.BroadcastToRoom("/", chatText.RoomName, "user-connect", msg)
	})
	server.OnEvent("/", "msg", func(s socketio.Conn, msg string) {
		chatText := ChatText{}
		err := json.Unmarshal([]byte(msg), &chatText)
		if err != nil {
			log.Fatal(err)
		}
		server.BroadcastToRoom("/", chatText.RoomName, "reply", msg)
	})
	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("assets")))
	log.Println("Start at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
