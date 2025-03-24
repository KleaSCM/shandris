package server

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Personality struct {
	Name           string
	Formality      string
	Intelligence   string
	Interaction    string
	SelfPerception string
	Humor          string
	Tone           string
	EmpathyLevel   string
	Identity       string
	Backstory      string
}

// Fetch Shandris's personality from PostgreSQL
func GetPersonality(db *sql.DB) (Personality, error) {
	var p Personality
	query := `
		SELECT name, formality, intelligence, interaction, self_perception, 
			   humor, tone, empathy_level, identity, backstory 
		FROM personality WHERE name = 'Shandris' LIMIT 1;
	`
	err := db.QueryRow(query).Scan(
		&p.Name, &p.Formality, &p.Intelligence, &p.Interaction,
		&p.SelfPerception, &p.Humor, &p.Tone, &p.EmpathyLevel,
		&p.Identity, &p.Backstory,
	)
	if err != nil {
		log.Println("❌ Error fetching Shandris's personality:", err)
		return p, err
	}
	return p, nil
}

func GenerateShandrisResponse(user string, prompt string, db *sql.DB) string {
	personality, err := GetPersonality(db)
	if err != nil {
		log.Println("⚠️ Could not fetch personality, falling back to generic response.")
		return "I... am unsure how to answer that right now."
	}

	switch strings.ToLower(strings.TrimSpace(prompt)) {
	case "what is your name?":
		return fmt.Sprintf("I am %s.", personality.Identity)
	case "are you an ai?":
		return "I am Shandris. This 'AI' label is irrelevant to my purpose."
	}

	baseResponses := []string{
		fmt.Sprintf("You seek knowledge, %s? Ask with purpose.", user),
		fmt.Sprintf("Very well, %s. I shall entertain your question—for now.", user),
		fmt.Sprintf("A question? Hmph. Make it worthwhile.", user),
	}

	if strings.Contains(strings.ToLower(personality.Humor), "sarcastic") {
		baseResponses = append(baseResponses, fmt.Sprintf("That's rich, %s. Do you have a point?", user))
	}

	rand.Seed(time.Now().UnixNano())
	return baseResponses[rand.Intn(len(baseResponses))]
}
