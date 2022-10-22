package wshandler

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Member struct {
	room *Room
	Role int `json:"role"`
	Id string `json:"id"`
	Username string `json:"username"`
	conn *websocket.Conn
	send chan[]byte
}

func (m *Member) ReadPump() {
	defer func() {
		m.room.unregister <- m
		m.conn.Close()
	}()

	for {
		_, message, err := m.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var payload Message[interface{}]
		
		json.Unmarshal(message, &payload)

		payload.sender = m

		m.room.event <- &payload
	}
}

func (m *Member) WritePump() {
	defer func() {
		m.conn.Close()
	}()
	
	for {
		select {
		case message, ok := <- m.send:
			
			if !ok {
				m.conn.Close()
				return
			}
			
			m.conn.WriteMessage(1, message)
		}
	}
}