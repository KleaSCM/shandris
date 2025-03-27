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
