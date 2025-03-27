package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// stripChainOfThought removes any <think>...</think> blocks (chain-of-thought).
func stripChainOfThought(resp string) string {
	// Remove <think>...</think> block (non-greedy match across newlines)
	re := regexp.MustCompile(`(?s)<think>.*?</think>\s*`)
	cleaned := re.ReplaceAllString(resp, "")

	// Remove any 'User:' and 'Assistant:' labels and separator lines
	reMeta := regexp.MustCompile(`(?m)^(User:|Assistant:|---)+\s*`)
	cleaned = reMeta.ReplaceAllString(cleaned, "")

	return strings.TrimSpace(cleaned)
}

type ChatRequest struct {
	SessionID string `json:"session_id"`
	Prompt    string `json:"prompt"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Set CORS and JSON content-type headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// Parse request JSON
	body, _ := io.ReadAll(r.Body)
	var req ChatRequest
	err := json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	newTopic := ClassifyPrompt(req.Prompt)
	currentTopic := GetCurrentTopic(req.SessionID)

	if currentTopic == "uncategorized" && newTopic != "uncategorized" {
		SetCurrentTopic(req.SessionID, newTopic)
		currentTopic = newTopic
		fmt.Println("🔥 Topic set to:", newTopic)
	}

	if newTopic != currentTopic && IsConfirmation(req.Prompt) {
		SetCurrentTopic(req.SessionID, newTopic)
		currentTopic = newTopic
	}

	personality, err := GetPersonality(db)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	history, err := GetChatHistoryByTopic(req.SessionID, currentTopic)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	context := BuildPrompt(personality, history, req.Prompt, currentTopic, newTopic)
	fullModelOutput, err := RunDeepSeek(context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cleanedOutput := stripChainOfThought(fullModelOutput)

	SaveChatHistory(req.SessionID, req.Prompt, cleanedOutput, newTopic)

	json.NewEncoder(w).Encode(ChatResponse{Response: cleanedOutput})
}

func IsConfirmation(prompt string) bool {
	p := strings.ToLower(prompt)
	return strings.Contains(p, "yes") ||
		strings.Contains(p, "sure") ||
		strings.Contains(p, "okay") ||
		strings.Contains(p, "go ahead")
}
