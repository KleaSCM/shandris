package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

type ChatRequest struct {
	Prompt string `json:"prompt"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read raw request body and print it for debugging
	body, _ := io.ReadAll(r.Body)
	fmt.Println("🔍 Raw Request Body:", string(body)) // Debug log

	// Try parsing JSON
	var req ChatRequest
	err := json.Unmarshal(body, &req)
	if err != nil {
		fmt.Println("❌ Error decoding JSON:", err) // Debug log
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("✅ Received Prompt:", req.Prompt) // Debug log

	// Run DeepSeek R1 via CLI
	cmd := exec.Command("ollama", "run", "deepseek-r1", req.Prompt)
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error running DeepSeek: %s", err), http.StatusInternalServerError)
		return
	}

	resp := ChatResponse{Response: string(output)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/chat", chatHandler)
	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
