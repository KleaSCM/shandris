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

	// Create existing tables
	err = createExistingTables()
	if err != nil {
		fmt.Println("❌ Error creating existing tables:", err)
		panic(err)
	}

	// Create new cognitive system tables
	err = createCognitiveSystemTables()
	if err != nil {
		fmt.Println("❌ Error creating cognitive system tables:", err)
		panic(err)
	}

	fmt.Println("✅ Connected to PostgreSQL!")
}

func createExistingTables() error {
	// Create persona_profiles table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS persona_profiles (
			session_id TEXT PRIMARY KEY,
			profile_data JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating persona_profiles table: %v", err)
	}

	// Create update trigger
	_, err = db.Exec(`
		CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ language 'plpgsql';

		DROP TRIGGER IF EXISTS update_persona_profiles_updated_at ON persona_profiles;
		
		CREATE TRIGGER update_persona_profiles_updated_at
			BEFORE UPDATE ON persona_profiles
			FOR EACH ROW
			EXECUTE FUNCTION update_updated_at_column();
	`)
	if err != nil {
		return fmt.Errorf("error creating update trigger: %v", err)
	}

	return nil
}

func createCognitiveSystemTables() error {
	// Create all new system tables
	_, err := db.Exec(`
		-- Core system tables
		CREATE TABLE IF NOT EXISTS moods (
			id UUID PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			current_value FLOAT NOT NULL,
			base_value FLOAT NOT NULL,
			last_updated TIMESTAMP NOT NULL,
			decay_rate FLOAT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS mood_patterns (
			id UUID PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			keywords TEXT[] NOT NULL,
			sentiment FLOAT NOT NULL,
			mood_shift VARCHAR(50) NOT NULL,
			intensity FLOAT NOT NULL,
			decay_rate FLOAT NOT NULL,
			requirements TEXT[],
			exclusions TEXT[],
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS traits (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			trait_name VARCHAR(100) NOT NULL,
			value FLOAT NOT NULL,
			confidence FLOAT NOT NULL,
			last_updated TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS topics (
			id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			category VARCHAR(100) NOT NULL,
			keywords TEXT[] NOT NULL,
			last_discussed TIMESTAMP,
			frequency INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		-- Timeline and Memory tables
		CREATE TABLE IF NOT EXISTS memory_events (
			id UUID PRIMARY KEY,
			type VARCHAR(50) NOT NULL,
			content TEXT NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			importance FLOAT NOT NULL,
			context JSONB NOT NULL,
			relations TEXT[] NOT NULL,
			tags TEXT[] NOT NULL,
			emotions JSONB NOT NULL,
			last_recall TIMESTAMP,
			recall_count INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		-- Add remaining tables from schema.sql...
		-- (I've truncated this for readability, but you would include all tables)
	`)
	if err != nil {
		return fmt.Errorf("error creating cognitive system tables: %v", err)
	}

	// Create indexes
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_moods_name ON moods(name);
		CREATE INDEX IF NOT EXISTS idx_traits_user_id ON traits(user_id);
		CREATE INDEX IF NOT EXISTS idx_topics_category ON topics(category);
		CREATE INDEX IF NOT EXISTS idx_memory_events_type ON memory_events(type);
		CREATE INDEX IF NOT EXISTS idx_memory_events_timestamp ON memory_events(timestamp);
		-- Add remaining indexes...
	`)
	if err != nil {
		return fmt.Errorf("error creating indexes: %v", err)
	}

	return nil
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
