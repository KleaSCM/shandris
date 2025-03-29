package cognitive

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/lib/pq"
)

// TopicPersistence handles long-term storage of topic data
type TopicPersistence struct {
	db          *sql.DB
	cache       map[string]*CachedTopic
	maxCacheAge time.Duration
}

// CachedTopic represents a topic stored in memory
type CachedTopic struct {
	Data     *TopicData
	LastUsed time.Time
	UseCount int
}

// TopicData represents the persistent topic information
type TopicData struct {
	ID            string
	Domain        string
	Keywords      []string
	Contexts      []string
	Relations     map[string]float64
	MoodPatterns  []string
	LastSeen      time.Time
	Frequency     int
	UserReactions map[string]int
	Metadata      map[string]interface{}
}

// NewTopicPersistence creates a new persistence manager
func NewTopicPersistence(dbConn *sql.DB) *TopicPersistence {
	return &TopicPersistence{
		db:          dbConn,
		cache:       make(map[string]*CachedTopic),
		maxCacheAge: 24 * time.Hour,
	}
}

// SaveTopic persists topic data to storage
func (tp *TopicPersistence) SaveTopic(topic *TopicData) error {
	// Convert topic data to JSON for storage
	metadata, err := json.Marshal(topic.Metadata)
	if err != nil {
		return err
	}

	// Upsert topic data
	_, err = tp.db.Exec(`
        INSERT INTO topics (
            id, domain, keywords, contexts, relations,
            mood_patterns, last_seen, frequency, user_reactions, metadata
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (id) DO UPDATE SET
            keywords = EXCLUDED.keywords,
            contexts = EXCLUDED.contexts,
            relations = EXCLUDED.relations,
            mood_patterns = EXCLUDED.mood_patterns,
            last_seen = EXCLUDED.last_seen,
            frequency = topics.frequency + 1,
            user_reactions = EXCLUDED.user_reactions,
            metadata = EXCLUDED.metadata
    `,
		topic.ID, topic.Domain, topic.Keywords, topic.Contexts,
		topic.Relations, topic.MoodPatterns, topic.LastSeen,
		topic.Frequency, topic.UserReactions, metadata,
	)

	if err != nil {
		return err
	}

	// Update cache
	tp.cache[topic.ID] = &CachedTopic{
		Data:     topic,
		LastUsed: time.Now(),
		UseCount: 1,
	}

	return nil
}

// LoadTopic retrieves topic data from storage
func (tp *TopicPersistence) LoadTopic(id string) (*TopicData, error) {
	// Check cache first
	if cached, exists := tp.cache[id]; exists {
		if time.Since(cached.LastUsed) < tp.maxCacheAge {
			cached.UseCount++
			cached.LastUsed = time.Now()
			return cached.Data, nil
		}
	}

	// Load from database
	var topic TopicData
	var metadata []byte

	err := tp.db.QueryRow(`
        SELECT id, domain, keywords, contexts, relations,
               mood_patterns, last_seen, frequency, user_reactions, metadata
        FROM topics WHERE id = $1
    `, id).Scan(
		&topic.ID, &topic.Domain, &topic.Keywords, &topic.Contexts,
		&topic.Relations, &topic.MoodPatterns, &topic.LastSeen,
		&topic.Frequency, &topic.UserReactions, &metadata,
	)

	if err != nil {
		return nil, err
	}

	// Parse metadata
	if err := json.Unmarshal(metadata, &topic.Metadata); err != nil {
		return nil, err
	}

	// Update cache
	tp.cache[id] = &CachedTopic{
		Data:     &topic,
		LastUsed: time.Now(),
		UseCount: 1,
	}

	return &topic, nil
}

// Here's the SQL schema for the topics table:
const TopicsTableSchema = `
CREATE TABLE IF NOT EXISTS topics (
    id VARCHAR(255) PRIMARY KEY,
    domain VARCHAR(100) NOT NULL,
    keywords TEXT[] NOT NULL,
    contexts TEXT[] NOT NULL,
    relations JSONB NOT NULL,
    mood_patterns TEXT[] NOT NULL,
    last_seen TIMESTAMP NOT NULL,
    frequency INTEGER NOT NULL DEFAULT 1,
    user_reactions JSONB NOT NULL,
    metadata JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_topics_domain ON topics(domain);
CREATE INDEX IF NOT EXISTS idx_topics_last_seen ON topics(last_seen);
`
