package server

import (
	"fmt"
	"log"
	"net/http"
)

func StartServer() {
	InitDB()

	http.HandleFunc("/api/chat", ChatHandler)

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
