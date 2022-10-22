package wshandler

import (
	"fmt"

	quiz_models "github.com/salmanrf/svelte-go-quiz-server/api/quiz/models"
)

// ? Room States:
// ? waiting-question, waiting-next, paused, ended
type Room struct {
	Title string `json:"title"`
	Code string `json:"code"`
	Quota int `json:"quota"`
	UsedQuota int `json:"used_quota"`
	NextQuestionInterval int `json:"next_question_interval"`
	QuestionTimeout int `json:"question_timeout"`
	members map[*Member]Member;
	admin *Member;
	state string;
	current_question int;
	broadcast chan []byte;
	register chan *Member;
	unregister chan *Member;
	event chan *Message[interface{}]
	quiz quiz_models.Quiz
}

var (
	RoomWaiting = "waiting"
	RoomWaitingNext = "waiting-next"
	RoomPaused = "waiting"
	RoomEnded = "ended"
)

var RoomHandlers = make(map[string]func(*Room, *Message[interface{}]))

func (room *Room) Init() {
	RoomHandlers["start"] = Start
	RoomHandlers["answer"] = OnAnswer
	RoomHandlers["register"] = OnJoin
	
	for {
		select {
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
		case message := <- room.event:
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("Unexpected Error", err)
				}

				room.Init()
			}()
			
			handler, ok := RoomHandlers[message.Event]

			fmt.Println("handler", ok)

			if !ok {
				response := Message[map[string]interface{}]{Event: "error", Data: make(map[string]interface{})}
				response.Data["message"] = "Invalid Message Event"

				write_err := message.sender.conn.WriteJSON(response)

				if write_err != nil {
					fmt.Println("Write Error", write_err)
				}
				
				return
			}
			
			go handler(room, message)
		}
	}
}