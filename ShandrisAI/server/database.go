package server

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Connect to PostgreSQL
func InitDB() {
	var err error
	connStr := "postgres://postgres:yourpassword@localhost/shandris_ai?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("❌ Database connection error:", err)
		panic(err)
	}

	fmt.Println("✅ Connected to PostgreSQL!")
}

// Save chat history
func SaveChatHistory(sessionID, userMessage, aiResponse string) {
	_, err := db.Exec("INSERT INTO chat_history (session_id, user_message, ai_response) VALUES ($1, $2, $3)", sessionID, userMessage, aiResponse)
	if err != nil {
		fmt.Println("❌ Error saving chat history:", err)
	}
}

// Retrieve chat history
func GetChatHistory(sessionID string) ([]string, error) {
	rows, err := db.Query("SELECT ai_response FROM chat_history WHERE session_id = $1 ORDER BY timestamp ASC", sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []string
	for rows.Next() {
		var response string
		if err := rows.Scan(&response); err != nil {
			return nil, err
		}
		history = append(history, response)
	}
	return history, nil
}
