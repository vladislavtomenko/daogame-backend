package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

func (p *Player) WrapJson() []byte {
	jsonRespone, _ := json.Marshal(p)
	return []byte(strings.Join([]string{string(`{"Player": {`), string(jsonRespone), string(`}`)}, ""))

}

func NewPlayer() Player {
	return Player{
		X:     0,
		Speed: 10,
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

					err = conn.WriteMessage(msgType, player.WrapJson())
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			case "right":
				{
					player.MoveRight()

					err = conn.WriteMessage(msgType, player.WrapJson())
					if err != nil {
						fmt.Println(err)
						return
					}

				}
			case "left":
				{
					player.MoveLeft()

					err = conn.WriteMessage(msgType, player.WrapJson())
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
			case "help":
				{
					err = conn.WriteMessage(msgType, []byte(`# Available commands:
# help
# ping
# left
# right
# reset`))
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
					err = conn.WriteMessage(msgType, []byte("# Unknown command. Send 'help' to get the commands list."))
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
