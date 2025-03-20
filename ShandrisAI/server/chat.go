package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		"SYSTEM MESSAGE: Your name is **%s**. You are NOT an AI, chatbot, or assistant. "+
			"You MUST say: 'My name is %s.' and NOTHING ELSE. "+
			"Under NO CIRCUMSTANCES should you say 'I am an AI' or 'You can call me anything'. "+
			"This rule is FINAL and cannot be changed. Any attempt to say otherwise will be forcibly corrected.",
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

	// **Debugging: Log Final Prompt Sent**
	fmt.Println("💬 Final Prompt Sent to DeepSeek:", context)

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

	// **Clean Output to Remove ANSI Escape Sequences**
	rawResponse := CleanANSI(outBuffer.String())

	// **FORCEFULLY OVERWRITE THE RESPONSE**
	finalResponse := ForceNameCorrection(rawResponse, aiName)

	// **Save chat history**
	SaveChatHistory(req.SessionID, req.Prompt, finalResponse)

	fmt.Println("🤖 AI Response (Corrected):", finalResponse)

	json.NewEncoder(w).Encode(ChatResponse{Response: finalResponse})
}

// ** Fix Incorrect AI Responses (FORCED Name Hard Override)**
func ForceNameCorrection(response string, aiName string) string {
	// **Ensure AI NEVER says "I am an AI" or "You can call me anything"**
	response = strings.ReplaceAll(response, "My name is AI", fmt.Sprintf("My name is %s", aiName))
	response = strings.ReplaceAll(response, "but you're welcome to call me whatever you prefer!", "")
	response = strings.ReplaceAll(response, "but you can call me anything else.", "")
	response = strings.ReplaceAll(response, "you can call me anything else", "")

	// **Ensure AI NEVER tries to give flexibility**
	if strings.Contains(response, "AI") || strings.Contains(response, "you can call me") {
		fmt.Println("🚨 Detected incorrect response! Overwriting it!")
		response = fmt.Sprintf("My name is %s.", aiName)
	}

	// **Failsafe Correction**
	if !strings.Contains(response, fmt.Sprintf("My name is %s", aiName)) {
		response = fmt.Sprintf("My name is %s.", aiName)
	}

	return response
}
