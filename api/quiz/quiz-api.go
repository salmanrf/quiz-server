package quiz

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/salmanrf/svelte-go-quiz-server/api/common"
	"github.com/salmanrf/svelte-go-quiz-server/api/quiz/dto"
	"github.com/salmanrf/svelte-go-quiz-server/quiz"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func CreateQuiz(w http.ResponseWriter, r *http.Request) {
	var body quiz.Quiz

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		log.Println("Error CreateQuizDto", err)
		http.Error(w, "Bad Input", http.StatusBadRequest)
		return
	}

	success := quiz.Create(body);

	var response = common.ApiResponse[interface{}]{Message: "Success"}

	if !success {
		response.Message = "Quiz with this code already exists."
		http.Error(w, response.Message, http.StatusBadRequest)
	}
	
	responseBody, _ := json.Marshal(response)

	w.Write(responseBody)
}

func FindQuizzes(w http.ResponseWriter, r *http.Request) {
	quizzes := quiz.Find()

	r.ParseForm()
	
	// page_number, _ := strconv.Atoi(r.Form.Get("page_number"))
	// page_size, _ := strconv.Atoi(r.Form.Get("page_size")) 
	
	// fmt.Println("page_number", page_number)
	// fmt.Println("page_size", page_size)
	
	data := common.FindRes[quiz.Quiz]{
		Items: quizzes,
		TotalItems: 100,
		PageSize: 100,
		PageNumber: 1,
	}
	
	response := common.ApiResponse[common.FindRes[quiz.Quiz]]{
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

func JoinQuizRoom(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var input dto.JoinQuizDto
	
	if err := decoder.Decode(&input); err != nil {
		log.Println("Error parsing join quiz dto.", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := quiz.Get(input.Code); !ok {
		message := fmt.Sprintf("Quiz Room %v doesn't exists.", input.Code)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	upgrader.CheckOrigin = func(*http.Request) bool {return true}
	
	_, serr := upgrader.Upgrade(w, r, nil)

	if serr != nil { 
		log.Println("Error upgrading to websocket.", serr)
		http.Error(w, serr.Error(), http.StatusInternalServerError)
		return
	}
}