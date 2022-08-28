package main

import (
	"flag"
	"log"
	"net/http"

	api "github.com/salmanrf/svelte-go-quiz-server/api/quiz"
	"github.com/salmanrf/svelte-go-quiz-server/api/room"
	"github.com/salmanrf/svelte-go-quiz-server/quiz"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	hub := quiz.NewHub()

	go hub.Run()

	http.HandleFunc("/api/quiz", func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.Write([]byte("OK"))

			return
		}

		if r.Method == http.MethodGet {
			api.FindQuizzes(w, r)
			return
		}
		
		if r.Method == http.MethodPost {
			api.CreateQuiz(w, r)
			return
		}

		http.Error(w, "Not Found", http.StatusNotFound)
	})

	http.HandleFunc("/api/room", func (w http.ResponseWriter, r * http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

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

	// http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	// 	quiz.ServeWs(hub, w, r)
	// })
	
	log.Println("Starting server at ", *addr)
	err := http.ListenAndServe(*addr, nil)
	
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}