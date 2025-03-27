package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type ChatRequest struct {
	SessionID string `json:"session_id"`
	Prompt    string `json:"prompt"`
}

type ChatResponse struct {
	Response string `json:"response"`
}

// Removes any <think>...</think> blocks and metadata from the model output.
func stripChainOfThought(resp string) string {
	re := regexp.MustCompile(`(?s)<think>.*?</think>\s*`)
	cleaned := re.ReplaceAllString(resp, "")
	reMeta := regexp.MustCompile(`(?m)^(User:|Assistant:|---)+\s*`)
	return strings.TrimSpace(reMeta.ReplaceAllString(cleaned, ""))
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

	// Set CORS and content type headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	body, _ := io.ReadAll(r.Body)
	var req ChatRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Handle mood clearing separately
	if detectMoodClear(req.Prompt) {
		SaveMemory(req.SessionID, "mood", "")
		fmt.Println("🧹 Cleared user mood.")
		json.NewEncoder(w).Encode(ChatResponse{Response: "Got it. Mood deleted. I’ll stop pretending you’re grumpy, even if your typing says otherwise. 😏"})
		return
	}

	// Handle memory-based inferences
	if name := extractName(req.Prompt); name != "" {
		SaveMemory(req.SessionID, "user_name", name)
		fmt.Println("🧠 Saved user name:", name)
	}
	if mood := extractMood(req.Prompt); mood != "" {
		SaveMemory(req.SessionID, "mood", mood)
		fmt.Println("🧠 Saved user mood:", mood)
	}

	// Topic tracking logic
	newTopic := ClassifyPrompt(req.Prompt)
	currentTopic := GetCurrentTopic(req.SessionID)

	if currentTopic == "uncategorized" && newTopic != "uncategorized" {
		SetCurrentTopic(req.SessionID, newTopic)
		currentTopic = newTopic
		fmt.Println("🔥 Topic set to:", newTopic)
	} else if newTopic != currentTopic {
		if IsConfirmation(req.Prompt) {
			SetCurrentTopic(req.SessionID, newTopic)
			currentTopic = newTopic
		} else {
			suggestion := fmt.Sprintf("It seems like you're switching topics from **%s** to **%s**. Should I go ahead and switch?", currentTopic, newTopic)
			json.NewEncoder(w).Encode(ChatResponse{Response: suggestion})
			return
		}
	}

	// Fetch persona and context
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

	context := BuildPrompt(personality, history, req.Prompt, currentTopic, newTopic, req.SessionID)
	fullModelOutput, err := RunDeepSeek(context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cleanedOutput := stripChainOfThought(fullModelOutput)
	SaveChatHistory(req.SessionID, req.Prompt, cleanedOutput, newTopic)
	json.NewEncoder(w).Encode(ChatResponse{Response: cleanedOutput})
}

// Basic yes/no/okay prompt confirmation parser.
func IsConfirmation(prompt string) bool {
	p := strings.ToLower(prompt)
	return strings.Contains(p, "yes") ||
		strings.Contains(p, "sure") ||
		strings.Contains(p, "okay") ||
		strings.Contains(p, "go ahead")
}
