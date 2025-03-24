package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	var req ChatRequest
	err := json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//Topic classification
	newTopic := ClassifyPrompt(req.Prompt)
	currentTopic := GetCurrentTopic(req.SessionID)

	if currentTopic == "uncategorized" && newTopic != "uncategorized" {
		SetCurrentTopic(req.SessionID, newTopic)
		currentTopic = newTopic
		fmt.Println("🔥 Topic set to:", newTopic)

	}

	//Context confirmation
	if newTopic != currentTopic && IsConfirmation(req.Prompt) {
		SetCurrentTopic(req.SessionID, newTopic)
		currentTopic = newTopic
	}

	//Load personality and topic-matched memory
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

	//Build prompt and run DeepSeek
	context := BuildPrompt(personality, history, req.Prompt, currentTopic, newTopic)
	rawResponse, err := RunDeepSeek(context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Save new message
	SaveChatHistory(req.SessionID, req.Prompt, rawResponse, newTopic)

	//Respond
	json.NewEncoder(w).Encode(ChatResponse{Response: rawResponse})
}
