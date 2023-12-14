package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"last-time-visit-go/event"
)

type GeoLocation struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Flag    string `json:"flag"`
}

func Visit(handlerEvent event.HandlerEvent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			var geo GeoLocation

			err := json.NewDecoder(r.Body).Decode(&geo)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			msg := event.Message{
				Name: "update",
				Data: geo,
			}

			handlerEvent.Broadcast(msg)

			w.WriteHeader(http.StatusOK)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		clientID := r.RemoteAddr + fmt.Sprintf("%d", time.Now().Unix())

		client, err := handlerEvent.Subscribe(clientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println("New client subscribed:", clientID)

		defer handlerEvent.Unsubscribe(clientID)

		for {
			select {
			case <-r.Context().Done():
				fmt.Println("client Unsubscribe:", clientID)
				return

			case msg := <-client.Receive():
				client.Send(msg, w)

				flusher.Flush()

				fmt.Println("Client message send!")
			}
		}
	}
}
