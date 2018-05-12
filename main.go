package main

import (
	"encoding/json"
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

type GameMap struct {
	X int `json:"X"`
}

func NewGameMap() GameMap {
	return GameMap{
		X: 0,
	}
}

func main() {
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		gameMap := NewGameMap()

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(msg))

			switch string(msg) {
			case "reset":
				{
					gameMap.X = 0

					jsonRespone, _ := json.Marshal(gameMap)
					err = conn.WriteMessage(msgType, []byte(string(jsonRespone)))
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			case "right":
				{
					gameMap.X = gameMap.X + 40

					jsonRespone, _ := json.Marshal(gameMap)
					err = conn.WriteMessage(msgType, []byte(string(jsonRespone)))
					if err != nil {
						fmt.Println(err)
						return
					}

				}
			case "left":
				{
					gameMap.X = gameMap.X - 40
					if gameMap.X < 0 {
						gameMap.X = 0
					}

					jsonRespone, _ := json.Marshal(gameMap)
					err = conn.WriteMessage(msgType, []byte(string(jsonRespone)))
					if err != nil {
						fmt.Println(err)
						return
					}

				}
			case "ping":
				{
					err = conn.WriteMessage(msgType, []byte("pong"))
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			case "close":
				{
					conn.Close()
					return
				}
			default:
				{
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
