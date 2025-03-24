package server

import "fmt"

// GetCurrentTopic fetches the last known topic for the session
func GetCurrentTopic(sessionID string) string {
	var topic string
	err := db.QueryRow(`
		SELECT current_topic FROM session_context WHERE session_id = $1
	`, sessionID).Scan(&topic)

	if err != nil {
		// Default to uncategorized if not found or error
		return "uncategorized"
	}
	return topic
}

// SetCurrentTopic updates or inserts the current topic for the session
func SetCurrentTopic(sessionID, topic string) {
	_, err := db.Exec(`
		INSERT INTO session_context (session_id, current_topic)
		VALUES ($1, $2)
		ON CONFLICT (session_id) DO UPDATE SET current_topic = EXCLUDED.current_topic
	`, sessionID, topic)

	if err != nil {
		fmt.Println("❌ Error updating current topic:", err)
	}
}
