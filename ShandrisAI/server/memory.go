package server

import (
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
