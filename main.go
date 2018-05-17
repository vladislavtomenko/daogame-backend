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

type Map struct {
	Size    int      `json:"size"`
	Objects []Object `json:"objects"`
}

// GetImpassableObjectsInRange returns a list of the impassable objects in the range
func (m *Map) GetImpassableObjectsInRange(x1 int, x2 int) []Object {
	var objList []Object
	for _, obj := range m.Objects {
		if x1 <= obj.X && x2 >= obj.X && obj.Passable == false {
			objList = append(objList, obj)
		}
	}
	return objList
}

// GetObjectsInRange returns a list of the all objects in the range
func (m *Map) GetObjectsInRange(x1 int, x2 int) []Object {
	var objList []Object
	for _, obj := range m.Objects {
		if x1 <= obj.X && x2 >= obj.X {
			objList = append(objList, obj)
		}
	}
	return objList
}

func (m *Map) WrapJson() []byte {
	json, _ := json.Marshal(m)
	return []byte(strings.Join([]string{`{"map":`, string(json), `}`}, ""))
}

type Object struct {
	Type     string `json:"type"`
	X        int    `json:"x"`
	Size     int    `json:"size"`
	Passable bool   `json:"passable"`
}

type Player struct {
	X     int `json:"x"`
	Speed int `json:"speed"`
}

func (p *Player) MoveLeft(gameMap *Map) {

	objList := gameMap.GetImpassableObjectsInRange(p.X-p.Speed, p.X)
	if len(objList) > 0 {
		newX := p.X - p.Speed
		for _, obj := range objList {
			if newX < obj.X {
				newX = obj.X + 1
			}
		}
		p.X = newX
	} else {
		p.X = p.X - p.Speed
		if p.X < 0 {
			p.X = 0
		}
	}

}

func (p *Player) MoveRight(gameMap *Map) {

	objList := gameMap.GetImpassableObjectsInRange(p.X, p.X+p.Speed)
	if len(objList) > 0 {
		newX := p.X + p.Speed
		for _, obj := range objList {
			if newX > obj.X {
				newX = obj.X - 1
			}
		}
		p.X = newX
	} else {
		p.X = p.X + p.Speed
		if p.X > gameMap.Size {
			p.X = gameMap.Size
		}
	}
}

// ResetLocation set player location to 0
func (p *Player) ResetLocation() {
	p.X = 0
	return
}

func (p *Player) WrapJson() []byte {
	json, _ := json.Marshal(p)
	return []byte(strings.Join([]string{`{"player":`, string(json), `}`}, ""))

}

func NewPlayer() Player {
	return Player{
		X:     0,
		Speed: 20,
	}
}

func NewRandomMap() Map {
	return Map{
		Size: 1000,
		Objects: []Object{
			Object{
				Type:     "fence",
				X:        939,
				Size:     5,
				Passable: false,
			},
			Object{
				Type:     "fence",
				X:        934,
				Size:     5,
				Passable: false,
			},
			Object{
				Type:     "balloon",
				X:        999,
				Size:     1,
				Passable: true,
			},
			Object{
				Type:     "balloon",
				X:        400,
				Size:     1,
				Passable: true,
			},
			Object{
				Type:     "balloon",
				X:        119,
				Size:     1,
				Passable: true,
			},
		},
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
		gameMap := NewRandomMap()

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}

			//fmt.Println(string(msg)) //// DEBUG:

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
			case "player":
				{
					err = conn.WriteMessage(msgType, player.WrapJson())
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			case "right":
				{
					player.MoveRight(&gameMap)

					err = conn.WriteMessage(msgType, player.WrapJson())
					if err != nil {
						fmt.Println(err)
						return
					}

				}
			case "left":
				{
					player.MoveLeft(&gameMap)

					err = conn.WriteMessage(msgType, player.WrapJson())
					if err != nil {
						fmt.Println(err)
						return
					}

				}
			case "map":
				{
					err = conn.WriteMessage(msgType, gameMap.WrapJson())
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
# help - show this message
# map - get the game map
# ping - ping backend
# left - move player to the left
# right - move player to the right
# player - get player obj
# reset - reset player location
# map`))
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
