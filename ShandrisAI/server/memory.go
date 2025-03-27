package server

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SaveMemory stores a key-value pair for a session in long_term_memory.
func SaveMemory(sessionID, key, value string) {
	_, err := db.Exec(`
		INSERT INTO long_term_memory (session_id, key, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (session_id, key) DO UPDATE SET value = EXCLUDED.value
	`, sessionID, key, value)
	if err != nil {
		fmt.Println("❌ Error saving memory:", err)
	}
}

// RecallMemory retrieves a value for a key in a session's memory.
func RecallMemory(sessionID, key string) (string, error) {
	var value string
	err := db.QueryRow(`
		SELECT value FROM long_term_memory
		WHERE session_id = $1 AND key = $2
	`, sessionID, key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

// SaveTraits stores the JSON traits blob for a session.
func SaveTraits(sessionID string, traits map[string]string) {
	blob, err := json.Marshal(traits)
	if err != nil {
		fmt.Println("❌ Error serializing traits:", err)
		return
	}
	_, err = db.Exec(`
		INSERT INTO persona_memory (session_id, traits)
		VALUES ($1, $2)
		ON CONFLICT (session_id) DO UPDATE SET traits = EXCLUDED.traits
	`, sessionID, blob)
	if err != nil {
		fmt.Println("❌ Error saving traits:", err)
	}
}

// RecallTraits returns the full trait map for a session.
func RecallTraits(sessionID string) (map[string]string, error) {
	var blob []byte
	traits := make(map[string]string)

	err := db.QueryRow(`
		SELECT traits FROM persona_memory WHERE session_id = $1
	`, sessionID).Scan(&blob)
	if err != nil {
		return traits, err
	}

	err = json.Unmarshal(blob, &traits)
	return traits, err
}

// extractName attempts to parse "my name is X" from the prompt.
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

// extractMood parses mood expressions from a prompt.
func extractMood(prompt string) string {
	prompt = strings.ToLower(prompt)

	moods := []string{
		"happy", "sad", "angry", "tired", "excited",
		"grumpy", "anxious", "stressed", "curious", "bored",
	}

	for _, mood := range moods {
		if strings.Contains(prompt, "i'm feeling "+mood) ||
			strings.Contains(prompt, "i feel "+mood) ||
			strings.Contains(prompt, "i am "+mood) {
			return mood
		}
	}
	return ""
}

// detectMoodClear returns true if the user is trying to erase mood memory.
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
