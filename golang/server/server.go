package server

import (
	"log"
	"net/http"

	"last-time-visit-go/event"
	"last-time-visit-go/handler"
)

func Start() {
	server := http.NewServeMux()

	handlerEvent := event.NewHandlerEvent()

	server.HandleFunc("/", handler.Home)
	server.HandleFunc("/visit", handler.Visit(handlerEvent))

	err := http.ListenAndServe(":8000", server)
	if err != nil {
		log.Fatal(err)
	}
}
