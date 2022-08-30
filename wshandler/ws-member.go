package wshandler

import (
	"log"

	"github.com/gorilla/websocket"
)

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

		m.room.broadcast <- message
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

