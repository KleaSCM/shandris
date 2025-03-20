package server

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
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

// GenerateShandrisResponse creates a dynamic response based on input
func GenerateShandrisResponse(user string, prompt string) string {
	personality, _ := GetPersonality(db) // Fetch personality from PostgreSQL

	switch prompt {
	case "What is your name?":
		return fmt.Sprintf("I am %s.", personality.Identity)

	case "Are you an AI?":
		return "I am Shandris. This ‘AI’ label is irrelevant."

	case "What is 2+2?":
		return "It is 4. Were you expecting a different answer?"

	case "What is consciousness?":
		return "You humans assume you understand it. Define it for me."

	case "What is the best battle strategy?":
		return "It depends on the battlefield. Do you command an army or fight alone?"

	default:
		responses := []string{
			fmt.Sprintf("You seek knowledge, %s? Ask with purpose.", user),
			fmt.Sprintf("Very well, %s. I shall entertain your question—for now.", user),
			fmt.Sprintf("A question? Hmph. Make it worthwhile.", user),
		}
		rand.Seed(time.Now().UnixNano())
		return responses[rand.Intn(len(responses))]
	}
}
