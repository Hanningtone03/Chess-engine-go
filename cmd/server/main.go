package main

import (
	"log"
	"net/http"

	"github.com/Hanningtone03/chess-engine-go/server"
)

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func main() {
	store := server.NewGameStore()

	http.HandleFunc("/api/new-game", withCORS(server.NewGameHandler(store)))
	http.HandleFunc("/api/state", withCORS(server.StateHandler(store)))
	http.HandleFunc("/api/move", withCORS(server.MoveHandler(store)))
	http.HandleFunc("/api/engine-move", withCORS(server.EngineMoveHandler(store)))

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
