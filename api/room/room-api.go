package room

import (
	"encoding/json"
	"fmt"
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
	dto := dto.JoinRoomDto{}

	r.ParseForm()

	if room_code := r.Form.Get("room_code"); room_code != "" {
		dto.RoomCode = room_code
	} else {
		http.Error(w, "invalid input, please specify room_code and username query in url", http.StatusBadRequest)
		return
	}

	if username := r.Form.Get("username"); username != "" {
		dto.Username = username
	} else {
		http.Error(w, "invalid input, please specify username and username query in url", http.StatusBadRequest)
		return
	}

	fmt.Println("dto", dto)
	
	join_err := wshandler.Join(dto, w, r)

	if join_err != nil {
		http.Error(w, join_err.Error(), http.StatusInternalServerError)
		return
	}
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