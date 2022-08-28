package room

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/salmanrf/svelte-go-quiz-server/api/common"
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

	room, err := wshandler.Create(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data common.ApiResponse[wshandler.Room]

	result, jsonerr := json.Marshal(data)

	if jsonerr != nil {
		http.Error(w, jsonerr.Error(), http.StatusInternalServerError)
		return
	}

	data.Data = *room

	w.Write(result)
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

func FindRooms(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	rooms := wshandler.Find()
	
	data := common.FindRes[wshandler.Room]{
		Items: rooms,
		PageSize: 100,
		PageNumber: 1,
		TotalItems: 100,
	}

	response := common.ApiResponse[common.FindRes[wshandler.Room]]{
		Message: "Success",
		Data: data,
	}

	responseBody, err := json.Marshal(response)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}