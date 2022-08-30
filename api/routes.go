package api

import (
	"net/http"

	"github.com/salmanrf/svelte-go-quiz-server/api/quiz"
	"github.com/salmanrf/svelte-go-quiz-server/api/room"
)

func corsMiddleware(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
}

func SetupHandler() {
	// ? Ws endpoints
	http.HandleFunc("/ws/room", func(w http.ResponseWriter, r *http.Request) {
		corsMiddleware(w, r)

		if r.Method == http.MethodOptions {
			w.Write([]byte("OK"))
			return
		}
		
		room.JoinRoom(w, r)
	})

	// http.HandleFunc("/ws/room", func (w http.ResponseWriter, r *http.Request) {
	// 	corsMiddleware(w, r)

	// 	if r.Method == http.MethodOptions {
	// 		w.Write([]byte("OK"))
	// 		return
	// 	}
		
	// 	if r.Method == http.MethodGet {
	// 		quiz.FindQuizzes(w, r)
	// 		return
	// 	}
		
	// 	if r.Method == http.MethodPost {
	// 		quiz.CreateQuiz(w, r)
	// 		return
	// 	}

	// 	http.Error(w, "Not Found", http.StatusNotFound)
	// })
	
	// ? Quizzes handler
	http.HandleFunc("/api/quiz", func (w http.ResponseWriter, r *http.Request) {
		corsMiddleware(w, r)

		if r.Method == http.MethodOptions {
			w.Write([]byte("OK"))
			return
		}
		
		if r.Method == http.MethodGet {
			quiz.FindQuizzes(w, r)
			return
		}
		
		if r.Method == http.MethodPost {
			quiz.CreateQuiz(w, r)
			return
		}

		http.Error(w, "Not Found", http.StatusNotFound)
	})

	// ? Rooms handler
	http.HandleFunc("/api/room", func (w http.ResponseWriter, r * http.Request) {
		corsMiddleware(w, r)

		if r.Method == http.MethodOptions {
			w.Write([]byte("OK"))
			return
		}
		
		if r.Method == http.MethodPost {
			room.CreateRoom(w, r)
			return
		}

		if r.Method == http.MethodGet {
			room.FindRooms(w, r)
			return
		}
		
		http.Error(w, "Not Found", http.StatusNotFound)
	})
}