package server

import (
	"encoding/json"
	"fmt"
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
