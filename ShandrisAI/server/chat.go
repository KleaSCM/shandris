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
// '(?s)' makes '.' match across newlines. '.*?' is a non-greedy match.
func stripChainOfThought(resp string) string {
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	return re.ReplaceAllString(resp, "")
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

	// Classify the prompt into a new topic
	newTopic := ClassifyPrompt(req.Prompt)

	// Fetch or default to uncategorized for the current topic
	currentTopic := GetCurrentTopic(req.SessionID)

	// If current topic is uncategorized and newTopic isn't, set new topic
	if currentTopic == "uncategorized" && newTopic != "uncategorized" {
		SetCurrentTopic(req.SessionID, newTopic)
		currentTopic = newTopic
		fmt.Println("🔥 Topic set to:", newTopic)
	}

	// If user confirms a topic switch, update it
	if newTopic != currentTopic && IsConfirmation(req.Prompt) {
		SetCurrentTopic(req.SessionID, newTopic)
		currentTopic = newTopic
	}

	// Load personality / config from DB
	personality, err := GetPersonality(db)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Retrieve this session's prior chat history for the chosen topic
	history, err := GetChatHistoryByTopic(req.SessionID, currentTopic)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Build final prompt & run model
	context := BuildPrompt(personality, history, req.Prompt, currentTopic, newTopic)
	rawResponse, err := RunDeepSeek(context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Strip out any <think>...</think> reasoning text
	rawResponse = stripChainOfThought(rawResponse)

	// Save final exchange to DB
	SaveChatHistory(req.SessionID, req.Prompt, rawResponse, newTopic)

	// Respond to the user
	json.NewEncoder(w).Encode(ChatResponse{Response: rawResponse})
}

// IsConfirmation checks if prompt text is user saying "yes" etc.
func IsConfirmation(prompt string) bool {
	p := strings.ToLower(prompt)
	return strings.Contains(p, "yes") ||
		strings.Contains(p, "sure") ||
		strings.Contains(p, "okay") ||
		strings.Contains(p, "go ahead")
}
