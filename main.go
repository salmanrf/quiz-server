package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/salmanrf/svelte-go-quiz-server/api"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	// http.HandleFunc("ws/room", func (w http.ResponseWriter, r * http.Request) {
	// 	fmt.Println("request", r)
		
	// 	room.JoinRoom(w, r)
	// })
	
	api.SetupHandler()
	
	log.Println("Starting server at ", *addr)
	err := http.ListenAndServe(*addr, nil)
	
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}