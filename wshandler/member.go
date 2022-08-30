package wshandler

import "github.com/gorilla/websocket"

type Member struct {
	room *Room
	Role int `json:"role"`
	Id string `json:"id"`
	Username string `json:"username"`
	conn *websocket.Conn
	send chan[]byte
}