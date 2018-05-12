package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			switch string(msg) {
			case "ping":
				{
					fmt.Println("ping")
					err = conn.WriteMessage(msgType, []byte("pong"))
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			case "close":
				{
					conn.Close()
					fmt.Println(string(msg))
					return
				}
			default:
				{
					fmt.Println("unknown command")
					err = conn.WriteMessage(msgType, []byte("Unknown command"))
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		}
	})
	http.ListenAndServe(":3000", nil)
}
