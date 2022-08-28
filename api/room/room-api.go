package room

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/salmanrf/svelte-go-quiz-server/api/room/dto"
	"github.com/salmanrf/svelte-go-quiz-server/quiz"
	"github.com/salmanrf/svelte-go-quiz-server/wshandler"
)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateRoomDto

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		log.Println("Error CreateRoomDto", err)
		http.Error(w, "Bad Input", http.StatusBadRequest)
	}

	_, exists := quiz.Get(body.QuizCode)

	if !exists {
		http.Error(w, "Quiz doesn't exist.", http.StatusNotFound)
		return
	}

	wshandler.Create(body, w, r)
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	var body dto.JoinRoomDto

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		log.Println("Error JoinRoomDto", err)
		http.Error(w, "Bad Input", http.StatusBadRequest)
	}

	wshandler.Join(body, w, r)
}