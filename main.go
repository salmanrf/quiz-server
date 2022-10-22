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

	api.SetupHandler()
	
	log.Println("Starting server at ", *addr)
	err := http.ListenAndServe(*addr, nil)
	
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}