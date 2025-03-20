package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
)

type ChatRequest struct {
	SessionID string `json:"session_id"`
	Prompt    string `json:"prompt"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	body, _ := io.ReadAll(r.Body)
	fmt.Println("🔍 Raw Request Body:", string(body))

	var req ChatRequest
	err := json.Unmarshal(body, &req)
	if err != nil {
		fmt.Println("❌ Error decoding JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("✅ Received Prompt:", req.Prompt)

	// Retrieve chat history for context
	history, err := GetChatHistory(req.SessionID)
	if err != nil {
		fmt.Println("❌ Error fetching chat history:", err)
	}

	// Build context
	var context string
	for _, entry := range history {
		context += entry + "\n"
	}
	context += "User: " + req.Prompt

	// Run Ollama with context
	cmd := exec.Command("ollama", "run", "deepseek-r1:8b", context)
	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &outBuffer

	err = cmd.Run()
	if err != nil {
		fmt.Println("❌ Error running DeepSeek R1:", err)
		http.Error(w, fmt.Sprintf("Error running DeepSeek: %s", err), http.StatusInternalServerError)
		return
	}

	cleanedOutput := CleanANSI(outBuffer.String()) // Remove escape codes

	// Save chat history
	SaveChatHistory(req.SessionID, req.Prompt, cleanedOutput)

	fmt.Println("🤖 AI Response:", cleanedOutput)

	json.NewEncoder(w).Encode(ChatResponse{Response: cleanedOutput})
}
