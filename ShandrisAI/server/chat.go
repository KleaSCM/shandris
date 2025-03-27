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

	// Attempt to extract user name and store it
	if strings.Contains(strings.ToLower(req.Prompt), "my name is") {
		name := extractName(req.Prompt)
		if name != "" {
			SaveMemory(req.SessionID, "user_name", name)
			fmt.Println("🧠 Saved user name:", name)
		}
	}

	// Attempt to extract mood and store it
	mood := extractMood(req.Prompt)
	if mood != "" {
		SaveMemory(req.SessionID, "mood", mood)
		fmt.Println("🧠 Saved user mood:", mood)
		if detectMoodClear(req.Prompt) {
			SaveMemory(req.SessionID, "mood", "")
			fmt.Println("🧹 Cleared user mood.")
			if detectMoodClear(req.Prompt) {
				SaveMemory(req.SessionID, "mood", "")
				fmt.Println("🧹 Cleared user mood.")
				clearedResponse := "Got it. Mood deleted. I’ll stop pretending you’re grumpy, even if your typing says otherwise. 😏"
				json.NewEncoder(w).Encode(ChatResponse{Response: clearedResponse})
				return
			}

		}

	}

	newTopic := ClassifyPrompt(req.Prompt)
	currentTopic := GetCurrentTopic(req.SessionID)

	if currentTopic == "uncategorized" && newTopic != "uncategorized" {
		SetCurrentTopic(req.SessionID, newTopic)
		currentTopic = newTopic
		fmt.Println("🔥 Topic set to:", newTopic)
	}

	if newTopic != currentTopic {
		if IsConfirmation(req.Prompt) {
			SetCurrentTopic(req.SessionID, newTopic)
			currentTopic = newTopic
		} else {
			// User didn't confirm yet — send clarification prompt instead of calling model
			suggestion := fmt.Sprintf("It seems like you're switching topics from **%s** to **%s**. Should I go ahead and switch?", currentTopic, newTopic)

			// Sendsuggestion as the assistant's reply and skip model
			json.NewEncoder(w).Encode(ChatResponse{Response: suggestion})
			return
		}
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

func IsConfirmation(prompt string) bool {
	p := strings.ToLower(prompt)
	return strings.Contains(p, "yes") ||
		strings.Contains(p, "sure") ||
		strings.Contains(p, "okay") ||
		strings.Contains(p, "go ahead")
}

func extractName(prompt string) string {
	prompt = strings.ToLower(prompt)
	idx := strings.Index(prompt, "my name is")
	if idx == -1 {
		return ""
	}
	namePart := strings.TrimSpace(prompt[idx+len("my name is"):])
	name := strings.Split(namePart, " ")[0]
	return strings.Title(name)
}
func extractMood(prompt string) string {
	prompt = strings.ToLower(prompt)

	moods := []string{"happy", "sad", "angry", "tired", "excited", "grumpy", "anxious", "stressed", "curious", "bored"}
	for _, mood := range moods {
		if strings.Contains(prompt, "i'm feeling "+mood) ||
			strings.Contains(prompt, "i feel "+mood) ||
			strings.Contains(prompt, "i am "+mood) {
			return mood
		}
	}
	return ""
}
func detectMoodClear(prompt string) bool {
	prompt = strings.ToLower(prompt)
	return strings.Contains(prompt, "forget my mood") ||
		strings.Contains(prompt, "reset my mood") ||
		strings.Contains(prompt, "ignore how i feel") ||
		strings.Contains(prompt, "never mind my feelings") ||
		strings.Contains(prompt, "i'm over it") ||
		strings.Contains(prompt, "it doesn't matter how i feel") ||
		strings.Contains(prompt, "change the subject") ||
		strings.Contains(prompt, "move on from that") ||
		strings.Contains(prompt, "stop talking about my mood")
}
