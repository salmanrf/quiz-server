package wshandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aidarkhanov/nanoid"
	"github.com/gorilla/websocket"
	"github.com/salmanrf/svelte-go-quiz-server/api/common"
	"github.com/salmanrf/svelte-go-quiz-server/api/room/dto"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

var rooms = make(map[string]*Room)

func Create(createDto dto.CreateRoomDto) (*Room, error) {
	// Upgrader.CheckOrigin = func (*http.Request) bool {return true}

	// ws, err := Upgrader.Upgrade(w, r, nil)

	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	alphabet := nanoid.DefaultAlphabet
	id, err := nanoid.Generate(alphabet, 10)

	if err != nil {
		// ws.WriteJSON(api.ApiResponse[interface{}]{Message: "Encountered Internal Error."})
		return nil, err
	}
	
	newRoom := &Room{
		Code: id,
		members: make(map[*Member]Member),
		broadcast: make(chan []byte),
		register: make(chan *Member),
		unregister: make(chan *Member),
	}

	rooms[newRoom.Code] = newRoom

	go newRoom.Init()
	
	return newRoom, nil
	
	// id, err = nanoid.Generate(alphabet, 10)
	
	// if err != nil {
	// 	ws.WriteJSON(api.ApiResponse[interface{}]{Message: "Encountered Internal Error."})
	// 	return nil, err
	// }
	
	// roomAdmin := &Member{
	// 	Id: id,
	// 	Role: 1,
	// 	room: newRoom,
	// 	conn: ws,
	// 	Username: createDto.AdminName,
	// }
	
	// go newRoom.Init()

	// fmt.Println("newRoom", newRoom)
	// fmt.Println("roomAdmin", roomAdmin)
	
	// newRoom.register <- roomAdmin
	
	// return newRoom, err
}

func Find() []Room {
	result := make([]Room, len(rooms))
	
	index := 0
	for _, r := range rooms {
		result[index] = *r
		
		index++
	}

	return result
}

func Join(joinDto dto.JoinRoomDto, w http.ResponseWriter, r * http.Request) (error) {
	room, exists := rooms[joinDto.RoomCode]

	fmt.Println("room", room)
	
	if !exists {
		fmt.Println("room doesn't exist")
		return errors.New("room doesn't exist")
	}
	
	Upgrader.CheckOrigin = func(*http.Request) bool {return true}
	
	ws, err := Upgrader.Upgrade(w, r, nil)
	
	if err != nil {
		fmt.Println("err", err)
		return errors.New("error upgrading to websocket")
	}
	
	alphabet := nanoid.DefaultAlphabet
	id, err := nanoid.Generate(alphabet, 10)
	
	if err != nil {
		ws.WriteJSON(common.ApiResponse[interface{}]{Message: err.Error()})
		ws.Close()
	}
	
	// ? The first user to enter the room is the admin
	role := 0

	if room.UsedQuota > 0 {
		role = 2
	} else {
		role = 1
	}
	
	newMember := &Member{
		Role: role,
		Username: joinDto.Username,
		conn: ws,
		send: make(chan []byte),
		room: room,
		Id: id,
	}

  room.register <- newMember

	go newMember.ReadPump()
	go newMember.WritePump()

	return nil
}

func (room *Room) Init() {
	for {
		select {
		case member := <- room.register:
			room.members[member] = *member

			join_message := JoinMessage{
				Username: member.Username, 
				RoomCode: room.Code,
			}

			json, err := json.Marshal(join_message)
			
			if err != nil {
				return
			}
			
			member.send <- json

		case member := <- room.unregister:
			_, exists := room.members[member]

			if exists {
				delete(room.members, member)
				close(member.send)
			}
		case message := <- room.broadcast:
			for m := range room.members {
				m.send <- message
			}
		}
	}
}