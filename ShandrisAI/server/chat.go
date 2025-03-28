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
	defer LogOperation("ChatHandler", map[string]interface{}{
		"method": r.Method,
		"path":   r.URL.Path,
	})(nil)

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
		LogError(err, "Failed to parse request body")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate session ID
	if req.SessionID == "" {
		LogError(fmt.Errorf("missing session ID"), "Request validation")
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	DebugLogger.Printf("📥 Received chat request - SessionID: %s, Prompt: %s", req.SessionID, req.Prompt)

	// Handle mood clearing separately
	if detectMoodClear(req.Prompt) {
		SaveMemory(req.SessionID, "mood", "")
		fmt.Println("🧹 Cleared user mood for session:", req.SessionID)
		json.NewEncoder(w).Encode(ChatResponse{Response: "Got it. Mood deleted. I'll stop pretending you're grumpy, even if your typing says otherwise. 😏"})
		return
	}

	// Handle memory-based inferences
	if profile, isProfileCreation := ExtractPersonaProfile(req.Prompt); isProfileCreation {
		err := SavePersonaProfile(req.SessionID, profile)
		if err != nil {
			LogError(err, "Failed to save persona profile")
		} else {
			InfoLogger.Printf("✅ Saved persona profile for session: %s", req.SessionID)
			json.NewEncoder(w).Encode(ChatResponse{Response: "I've stored your personal information permanently. I'll remember these details about you, even across different conversations. Is there anything specific from your background you'd like to discuss?"})
			return
		}
	}

	if name := extractName(req.Prompt); name != "" {
		// Check if this name matches an existing profile
		if existingSessionID, err := FindUserByName(name); err == nil {
			DebugLogger.Printf("🔍 Found existing profile for name %s (SessionID: %s)", name, existingSessionID)
			// If current session doesn't have a profile, use the existing one
			if !HasExistingProfile(req.SessionID) {
				profile, err := GetPersonaProfile(existingSessionID)
				if err == nil {
					SavePersonaProfile(req.SessionID, profile)
					InfoLogger.Printf("🔄 Linked existing profile for %s to new session: %s", name, req.SessionID)
				}
			}
		}

		SaveMemory(req.SessionID, "user_name", name)
		InfoLogger.Printf("🧠 Saved user name: %s for session: %s", name, req.SessionID)
	}
	if mood := extractMood(req.Prompt); mood != "" {
		SaveMemory(req.SessionID, "mood", mood)
		InfoLogger.Printf("🧠 Saved user mood: %s for session: %s", mood, req.SessionID)
	}

	// Topic tracking logic
	newTopic := ClassifyPrompt(req.Prompt)
	currentTopic := GetCurrentTopic(req.SessionID)

	DebugLogger.Printf("📊 Topic Analysis - Current: %s, New: %s", currentTopic, newTopic)

	if currentTopic == "uncategorized" && newTopic != "uncategorized" {
		SetCurrentTopic(req.SessionID, newTopic)
		currentTopic = newTopic
		InfoLogger.Printf("🔥 Topic set to: %s for session: %s", newTopic, req.SessionID)
	} else if newTopic != currentTopic && newTopic != "uncategorized" {
		SetCurrentTopic(req.SessionID, newTopic)
		currentTopic = newTopic
		InfoLogger.Printf("🔄 Topic naturally transitioned from %s to %s for session: %s", currentTopic, newTopic, req.SessionID)
	}

	// Fetch persona and context
	personality, err := GetPersonality(db)
	if err != nil {
		LogError(err, "Failed to fetch personality")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	history, err := GetChatHistoryByTopic(req.SessionID, currentTopic)
	if err != nil {
		LogError(err, "Failed to fetch chat history")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	DebugLogger.Printf("📚 Retrieved %d historical chat turns for topic %s", len(history), currentTopic)

	context := BuildPrompt(personality, history, req.Prompt, currentTopic, newTopic, req.SessionID)
	DebugLogger.Printf("🎯 Built context for model (length: %d characters)", len(context))

	fullModelOutput, err := RunDeepSeek(context)
	if err != nil {
		LogError(err, "Failed to get model response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cleanedOutput := stripChainOfThought(fullModelOutput)
	LogChatOperation("Saving chat history", req.SessionID, req.Prompt, currentTopic)
	SaveChatHistory(req.SessionID, req.Prompt, cleanedOutput, newTopic)

	InfoLogger.Printf("💬 Chat response generated - Length: %d characters", len(cleanedOutput))
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
