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
	connStr := "postgres://postgres:Hisako1086@localhost/shandris_ai?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("❌ Database connection error:", err)
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("❌ Database ping failed:", err)
		panic(err)
	}

	fmt.Println("✅ Connected to PostgreSQL!")
}

// Fetch Shandris' name from the database
func GetAIName() string {
	var aiName string
	err := db.QueryRow(`
		SELECT value 
		FROM system_memory 
		WHERE key = 'ai_name' 
		LIMIT 1
	`).Scan(&aiName)

	if err != nil {
		fmt.Println("❌ Error fetching AI name:", err)
		return "Shandris" // fallback value
	}

	fmt.Println("🧠 AI Name from DB:", aiName)
	return aiName
}

// Save chat history to PostgreSQL
func SaveChatHistory(sessionID, userMessage, aiResponse, topic string) {
	_, err := db.Exec(`
		INSERT INTO chat_history (session_id, user_message, ai_response, topic)
		VALUES ($1, $2, $3, $4)
	`, sessionID, userMessage, aiResponse, topic)

	if err != nil {
		fmt.Println("❌ Error saving chat history:", err)
	}
}

// ChatTurn represents a single user/assistant exchange
type ChatTurn struct {
	UserMessage string
	AIResponse  string
}

// Retrieve topic-specific chat history
func GetChatHistoryByTopic(sessionID, topic string) ([]ChatTurn, error) {
	rows, err := db.Query(`
		SELECT user_message, ai_response 
		FROM chat_history 
		WHERE session_id = $1 AND topic = $2 
		ORDER BY timestamp ASC
	`, sessionID, topic)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []ChatTurn
	for rows.Next() {
		var turn ChatTurn
		if err := rows.Scan(&turn.UserMessage, &turn.AIResponse); err != nil {
			return nil, err
		}
		history = append(history, turn)
	}

	return history, nil
}
