package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
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

	// Retrieve AI name from PostgreSQL
	aiName := GetAIName()
	fmt.Println("🧠 AI Name from DB:", aiName)

	// **System Prompt to Lock Identity**
	systemPrompt := fmt.Sprintf(
		"SYSTEM MESSAGE: Your name is **%s**. "+
			"If asked 'What is your name?', ONLY THEN respond with: 'My name is %s.' "+
			"Otherwise, answer intelligently based on the context. "+
			"DO NOT randomly repeat your name unless directly asked.",
		aiName, aiName,
	)

	// **Build AI Memory with System Prompt**
	var context string
	context += systemPrompt + "\n\n"
	history, err := GetChatHistory(req.SessionID)
	if err != nil {
		fmt.Println("❌ Error fetching chat history:", err)
	}

	for _, entry := range history {
		context += entry + "\n"
	}

	context += "User: " + req.Prompt

	// **Write Context to a Temp File**
	tmpFile, err := os.CreateTemp("", "deepseek_input_*.txt")
	if err != nil {
		fmt.Println("❌ Error creating temp file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile.Name()) // Cleanup temp file after execution

	_, err = tmpFile.WriteString(context)
	tmpFile.Close()
	if err != nil {
		fmt.Println("❌ Error writing to temp file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("C:\\Users\\aikaw\\AppData\\Local\\Programs\\Ollama\\ollama.exe", "run", "deepseek-r1:8b")
	cmd.Stdin = strings.NewReader(context) // Pipe the context directly

	// Log the command being run
	fmt.Println("🔧 Running command:", cmd.String())

	var outBuffer, errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err = cmd.Run()
	if err != nil {
		fmt.Println("❌ Error running DeepSeek R1:", err)
		fmt.Println("📌 STDERR:", errBuffer.String()) // Log stderr output for debugging
		fmt.Println("📌 STDOUT:", outBuffer.String()) // Log stdout output for debugging
		http.Error(w, fmt.Sprintf("Error running DeepSeek: %s\n%s", err, errBuffer.String()), http.StatusInternalServerError)
		return
	}

	// **Clean Output to Remove ANSI Escape Sequences**
	rawResponse := CleanANSI(outBuffer.String())

	// **Save chat history**
	SaveChatHistory(req.SessionID, req.Prompt, rawResponse)

	fmt.Println("🤖 AI Response:", rawResponse)

	json.NewEncoder(w).Encode(ChatResponse{Response: rawResponse})
}
