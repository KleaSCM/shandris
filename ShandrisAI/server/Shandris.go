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

// Generate Shandris's response based on user input
func GenerateShandrisResponse(user string, prompt string, db *sql.DB) string {
	personality, _ := GetPersonality(db) // Fetch personality from PostgreSQL

	// Handle specific prompts where Shandris mentions her name
	if strings.ToLower(prompt) == "what is your name?" {
		return fmt.Sprintf("I am %s.", personality.Identity)
	}

	// Handle prompt to check if AI, and include Shandris's perspective
	if strings.ToLower(prompt) == "are you an ai?" {
		return "I am Shandris. This 'AI' label is irrelevant."
	}

	// Default case for all other prompts: generate a more general response
	responses := []string{
		fmt.Sprintf("You seek knowledge, %s? Ask with purpose.", user),
		fmt.Sprintf("Very well, %s. I shall entertain your question—for now.", user),
		fmt.Sprintf("A question? Hmph. Make it worthwhile.", user),
	}

	// Optional: add personality-based customization to responses (e.g., sarcasm)
	if strings.Contains(personality.Humor, "sarcastic") {
		responses = append(responses, fmt.Sprintf("That's rich, %s. Do you have a point?", user))
	}

	// Choose a random response based on Shandris's personality traits
	rand.Seed(time.Now().UnixNano())
	return responses[rand.Intn(len(responses))]
}
