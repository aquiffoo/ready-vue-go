package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/rs/cors"
)

var messages []Message
var users map[string]bool = make(map[string]bool)
var usersMu sync.Mutex

type Message struct {
	Message string `json:"message"`
	Sender string `json:"sender"`
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodPost {
		var user struct {
			Username string `json:"username"`
		}

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "invalid user", http.StatusBadRequest)
			return
		}

		usersMu.Lock()
		users[user.Username] = true
		usersMu.Unlock()
	}
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodPost {
		var msg Message
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(w, "invalid message", http.StatusBadRequest)
			return
		}

		usersMu.Lock()
		if _, ok := users[msg.Sender]; !ok {
			usersMu.Unlock()
			http.Error(w, "user not found", http.StatusBadRequest)
			return
		}
		usersMu.Unlock()

		messages = append(messages, msg)
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("message sent successfully"))
	}
}

func handleGetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func main() {
	http.HandleFunc("/users", handleCreateUser)
	http.HandleFunc("/message", handleMessage)
	http.HandleFunc("/messages", handleGetMessages)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8081"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}).Handler(http.DefaultServeMux)

	port := 8080

	fmt.Printf("server listening on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}
