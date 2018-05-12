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

type Player struct {
	X     int `json:"X"`
	Speed int `json:"Speed"`
}

func (p *Player) MoveLeft() {
	p.X = p.X - p.Speed
	if p.X < 0 {
		p.X = 0
	}
	return
}

func (p *Player) MoveRight() {
	p.X = p.X + p.Speed
	return
}

func (p *Player) ResetLocation() {
	p.X = 0
	return
}

func NewPlayer() Player {
	return Player{
		X:     0,
		Speed: 40,
	}
}

func main() {
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		player := NewPlayer()

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
					player.ResetLocation()

					jsonRespone, _ := json.Marshal(player)
					err = conn.WriteMessage(msgType, []byte(string(jsonRespone)))
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			case "right":
				{
					player.MoveRight()

					jsonRespone, _ := json.Marshal(player)
					err = conn.WriteMessage(msgType, []byte(string(jsonRespone)))
					if err != nil {
						fmt.Println(err)
						return
					}

				}
			case "left":
				{
					player.MoveLeft()

					jsonRespone, _ := json.Marshal(player)
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
