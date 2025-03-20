package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
)

type ChatRequest struct {
	Prompt string `json:"prompt"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

// Function to clean ANSI escape codes from output
func cleanANSI(input string) string {
	ansiRegex := regexp.MustCompile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)
	return ansiRegex.ReplaceAllString(input, "")
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
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

	// Run Ollama DeepSeek R1 via CLI
	cmd := exec.Command("ollama", "run", "deepseek-r1:8b", req.Prompt)
	var outBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &outBuffer

	err = cmd.Run()
	if err != nil {
		fmt.Println("❌ Error running DeepSeek R1:", err)
		http.Error(w, fmt.Sprintf("Error running DeepSeek: %s", err), http.StatusInternalServerError)
		return
	}

	cleanedOutput := cleanANSI(outBuffer.String()) // Remove terminal escape codes

	fmt.Println("🤖 Clean DeepSeek R1 Response:", cleanedOutput) // Debug log

	json.NewEncoder(w).Encode(ChatResponse{Response: cleanedOutput})
}

func main() {
	http.HandleFunc("/api/chat", chatHandler)
	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
