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
	if r.Method == http.MethodOptions { // Handle CORS preflight requests
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow frontend
	w.Header().Set("Content-Type", "application/json") // Set response type

	// Read and print request
	body, _ := io.ReadAll(r.Body)
	fmt.Println("🔍 Raw Request Body:", string(body))

	// Parse JSON
	var req ChatRequest
	err := json.Unmarshal(body, &req)
	if err != nil {
		fmt.Println("❌ Error decoding JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("✅ Received Prompt:", req.Prompt)

	// Run Ollama DeepSeek R1 via CLI and log output
	cmd := exec.Command("ollama", "run", "deepseek-r1:8b", req.Prompt) // Explicit model version

	fmt.Println("🔧 Running Ollama Command:", cmd.String()) // Log the exact command

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("❌ Error running DeepSeek R1:", err)
		http.Error(w, fmt.Sprintf("Error running DeepSeek: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("🤖 DeepSeek R1 Response:", string(output)) // Log response

	// Return response
	json.NewEncoder(w).Encode(ChatResponse{Response: string(output)})
}

func main() {
	http.HandleFunc("/api/chat", chatHandler)
	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
