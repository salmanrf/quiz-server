package wshandler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aidarkhanov/nanoid"
	"github.com/gorilla/websocket"
	"github.com/salmanrf/svelte-go-quiz-server/api/common"
	"github.com/salmanrf/svelte-go-quiz-server/api/quiz"
	"github.com/salmanrf/svelte-go-quiz-server/api/room/dto"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

var rooms = make(map[string]*Room)

func Create(createDto dto.CreateRoomDto) (*Room, error) {
	quiz, exists := quiz.Get(createDto.QuizCode)
	
	if !exists {
		fmt.Println("Error Creating room: Quiz doesnt exist.")
		
		return nil, errors.New("quiz doesn't exist")
	}
	
	alphabet := nanoid.DefaultAlphabet
	id, err := nanoid.Generate(alphabet, 10)

	if err != nil {
		return nil, err
	}
	
	newRoom := &Room{
		Title: createDto.Title,
		Code: id,
		Quota: createDto.Quota,
		QuestionTimeout: createDto.QuestionTimeout,
		NextQuestionInterval: createDto.NextQuestionInterval,
		quiz: *quiz,
		state: RoomWaitingNext,
		members: make(map[*Member]Member),
		broadcast: make(chan []byte),
		register: make(chan *Member),
		unregister: make(chan *Member),
		event: make(chan *Message[interface{}]),
	}

	rooms[newRoom.Code] = newRoom

	go newRoom.Init()
	
	return newRoom, nil
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

	if !exists {
		fmt.Println("Room doesn't exist")
		return errors.New("room doesn't exist")
	}
	
	if room.UsedQuota + 1 > room.Quota {
		fmt.Println("Room has reached maximum quota")	
		return errors.New("room has reached maximum quota")	
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
	role := 1

	if room.UsedQuota > 0 {
		role = 2
	}

	// ? Increment current room used quota
	room.UsedQuota++
	
	newMember := &Member{
		Role: role,
		Username: joinDto.Username,
		conn: ws,
		send: make(chan []byte),
		room: room,
		Id: id,
	}

	defer func() {
		message := &Message[interface{}]{
			Event: "register", 
			sender: newMember,
		}
		
		room.event <- message
	}()

	go newMember.ReadPump()
	go newMember.WritePump()
	
	return nil
}

func Start(r *Room, _ *Message[interface{}]) {
	interval := r.NextQuestionInterval // ? Interval to next question in milliseconds
	timeout := r.QuestionTimeout // ? Question timeout in milliseconds

	r.state = RoomWaitingNext
	<- time.After(time.Duration(interval) * time.Millisecond)
	
	for _, question := range r.quiz.Questions {
		r.state = RoomWaiting
		
		fmt.Println("question", question)

		<- time.After(time.Duration(timeout) * time.Millisecond)
		
		<- time.After(time.Duration(interval) * time.Millisecond)

		r.state = RoomWaitingNext
	}

	fmt.Println("Finish")
}

func OnAnswer(r *Room, message *Message[interface{}]) {
	
}

func OnJoin(r *Room, message *Message[interface{}]) {
	new_member := message.sender
	
	if new_member.Role == 1 {
		r.admin = new_member
	} else {
		r.members[message.sender] = *message.sender
	}
		
	room := Room{
		Code: r.Code,
		Title: r.Title,
		Quota: r.Quota,
		UsedQuota: r.UsedQuota,
	}
	
	join_message := Message[JoinMessage] {
		Event: "room_join",
		Data: JoinMessage{
			Id: new_member.Id, 
			Username: new_member.Username, 
			Role: new_member.Role,
			Room: room,
		},
	}

	new_member.conn.WriteJSON(join_message)

	join_message.Event = "member_join"
	
	r.admin.conn.WriteJSON(join_message)
	
	for m := range r.members {
		if m == new_member {
			continue
		}

		m.conn.WriteJSON(join_message)
	}
}