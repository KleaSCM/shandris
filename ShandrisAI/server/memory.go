package server

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PersonaProfile struct {
	Name       string            `json:"name"`
	Biography  string            `json:"biography"`
	Attributes map[string]string `json:"attributes"`
}

// SavePersonaProfile stores a complete persona profile for a user
func SavePersonaProfile(sessionID string, profile PersonaProfile) error {
	defer LogOperation("SavePersonaProfile", map[string]interface{}{
		"sessionID": sessionID,
		"name":      profile.Name,
	})(nil)

	LogProfileOperation("Saving persona profile", sessionID, profile)

	profileJSON, err := json.Marshal(profile)
	if err != nil {
		LogError(err, "Failed to serialize profile")
		return fmt.Errorf("error serializing profile: %w", err)
	}

	_, err = db.Exec(`
		INSERT INTO persona_profiles (session_id, profile_data)
		VALUES ($1, $2)
		ON CONFLICT (session_id) DO UPDATE SET profile_data = EXCLUDED.profile_data
	`, sessionID, profileJSON)

	if err != nil {
		LogError(err, "Failed to save profile to database")
		return fmt.Errorf("error saving profile: %w", err)
	}
	return nil
}

// GetPersonaProfile retrieves a user's complete persona profile
func GetPersonaProfile(sessionID string) (PersonaProfile, error) {
	defer LogOperation("GetPersonaProfile", map[string]interface{}{
		"sessionID": sessionID,
	})(nil)

	var profile PersonaProfile
	var profileJSON []byte

	err := db.QueryRow(`
		SELECT profile_data FROM persona_profiles 
		WHERE session_id = $1
	`, sessionID).Scan(&profileJSON)

	if err != nil {
		LogError(err, "Failed to retrieve profile from database")
		return profile, err
	}

	err = json.Unmarshal(profileJSON, &profile)
	if err != nil {
		LogError(err, "Failed to deserialize profile")
		return profile, err
	}

	LogProfileOperation("Retrieved persona profile", sessionID, profile)
	return profile, err
}

// ExtractPersonaProfile attempts to parse a profile creation request
func ExtractPersonaProfile(prompt string) (PersonaProfile, bool) {
	defer LogOperation("ExtractPersonaProfile", map[string]interface{}{
		"promptLength": len(prompt),
	})(nil)

	var profile PersonaProfile
	profile.Attributes = make(map[string]string)

	// Check if this is a profile creation request
	if !strings.Contains(strings.ToLower(prompt), "remember the following about me") &&
		!strings.Contains(strings.ToLower(prompt), "store this permanently") {
		return profile, false
	}

	DebugLogger.Printf("🔍 Detected profile creation request: %s", prompt)

	// Split into lines and process each attribute
	lines := strings.Split(prompt, "\n")
	var biography strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "-") {
			continue
		}

		// Remove the leading "-" and trim
		line = strings.TrimSpace(strings.TrimPrefix(line, "-"))

		// Extract name if present
		if strings.HasPrefix(strings.ToLower(line), "my name is") {
			profile.Name = strings.TrimSpace(strings.TrimPrefix(line, "My name is"))
			profile.Name = strings.TrimSpace(strings.TrimPrefix(profile.Name, "my name is"))
			DebugLogger.Printf("📝 Extracted name: %s", profile.Name)
			continue
		}

		// Add to biography
		if biography.Len() > 0 {
			biography.WriteString("\n")
		}
		biography.WriteString(line)

		// Extract key attributes
		lower := strings.ToLower(line)
		if strings.Contains(lower, "work") {
			profile.Attributes["occupation"] = line
			DebugLogger.Printf("📝 Extracted occupation: %s", line)
		}
		if strings.Contains(lower, "background") {
			profile.Attributes["background"] = line
			DebugLogger.Printf("📝 Extracted background: %s", line)
		}
		if strings.Contains(lower, "live") || strings.Contains(lower, "living") {
			profile.Attributes["location"] = line
			DebugLogger.Printf("📝 Extracted location: %s", line)
		}
		if strings.Contains(lower, "hobby") || strings.Contains(lower, "interest") ||
			strings.Contains(lower, "love") || strings.Contains(lower, "enjoy") {
			profile.Attributes["interests"] = line
			DebugLogger.Printf("📝 Extracted interests: %s", line)
		}
		if strings.Contains(lower, "language") || strings.Contains(lower, "programming") {
			profile.Attributes["tech_stack"] = line
			DebugLogger.Printf("📝 Extracted tech stack: %s", line)
		}
	}

	profile.Biography = biography.String()
	DebugLogger.Printf("📝 Compiled biography: %s", profile.Biography)
	return profile, true
}

// SaveMemory stores a key-value pair for a session in long_term_memory.
func SaveMemory(sessionID, key, value string) {
	defer LogOperation("SaveMemory", map[string]interface{}{
		"sessionID": sessionID,
		"key":       key,
	})(nil)

	LogMemoryOperation("Saving memory", sessionID, key, value)

	_, err := db.Exec(`
		INSERT INTO long_term_memory (session_id, key, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (session_id, key) DO UPDATE SET value = EXCLUDED.value
	`, sessionID, key, value)
	if err != nil {
		LogError(err, "Failed to save memory")
		fmt.Println("❌ Error saving memory:", err)
	}
}

// RecallMemory retrieves a value for a key in a session's memory.
func RecallMemory(sessionID, key string) (string, error) {
	defer LogOperation("RecallMemory", map[string]interface{}{
		"sessionID": sessionID,
		"key":       key,
	})(nil)

	var value string
	err := db.QueryRow(`
		SELECT value FROM long_term_memory
		WHERE session_id = $1 AND key = $2
	`, sessionID, key).Scan(&value)

	if err != nil {
		LogError(err, "Failed to recall memory")
		return "", err
	}

	LogMemoryOperation("Retrieved memory", sessionID, key, value)
	return value, nil
}

// SaveTraits stores the JSON traits blob for a session.
func SaveTraits(sessionID string, traits map[string]string) {
	blob, err := json.Marshal(traits)
	if err != nil {
		fmt.Println("❌ Error serializing traits:", err)
		return
	}
	_, err = db.Exec(`
		INSERT INTO persona_memory (session_id, traits)
		VALUES ($1, $2)
		ON CONFLICT (session_id) DO UPDATE SET traits = EXCLUDED.traits
	`, sessionID, blob)
	if err != nil {
		fmt.Println("❌ Error saving traits:", err)
	}
}

// RecallTraits returns the full trait map for a session.
func RecallTraits(sessionID string) (map[string]string, error) {
	var blob []byte
	traits := make(map[string]string)

	err := db.QueryRow(`
		SELECT traits FROM persona_memory WHERE session_id = $1
	`, sessionID).Scan(&blob)
	if err != nil {
		return traits, err
	}

	err = json.Unmarshal(blob, &traits)
	return traits, err
}

// extractName attempts to parse "my name is X" from the prompt.
func extractName(prompt string) string {
	prompt = strings.ToLower(prompt)
	idx := strings.Index(prompt, "my name is")
	if idx == -1 {
		return ""
	}
	namePart := strings.TrimSpace(prompt[idx+len("my name is"):])
	name := strings.Split(namePart, " ")[0]
	return strings.Title(name)
}

// extractMood parses mood expressions from a prompt.
func extractMood(prompt string) string {
	prompt = strings.ToLower(prompt)

	moods := []string{
		"happy", "sad", "angry", "tired", "excited",
		"grumpy", "anxious", "stressed", "curious", "bored",
	}

	for _, mood := range moods {
		if strings.Contains(prompt, "i'm feeling "+mood) ||
			strings.Contains(prompt, "i feel "+mood) ||
			strings.Contains(prompt, "i am "+mood) {
			return mood
		}
	}
	return ""
}

// detectMoodClear returns true if the user is trying to erase mood memory.
func detectMoodClear(prompt string) bool {
	prompt = strings.ToLower(prompt)
	return strings.Contains(prompt, "forget my mood") ||
		strings.Contains(prompt, "reset my mood") ||
		strings.Contains(prompt, "ignore how i feel") ||
		strings.Contains(prompt, "never mind my feelings") ||
		strings.Contains(prompt, "i'm over it") ||
		strings.Contains(prompt, "it doesn't matter how i feel") ||
		strings.Contains(prompt, "change the subject") ||
		strings.Contains(prompt, "move on from that") ||
		strings.Contains(prompt, "stop talking about my mood")
}

// FindUserByName attempts to find a session ID for a user with the given name
func FindUserByName(name string) (string, error) {
	var sessionID string
	err := db.QueryRow(`
		SELECT session_id 
		FROM persona_profiles 
		WHERE profile_data->>'name' ILIKE $1
		ORDER BY updated_at DESC 
		LIMIT 1
	`, name).Scan(&sessionID)

	if err != nil {
		return "", err
	}
	return sessionID, nil
}

// HasExistingProfile checks if a session ID has a stored profile
func HasExistingProfile(sessionID string) bool {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM persona_profiles 
			WHERE session_id = $1
		)
	`, sessionID).Scan(&exists)

	if err != nil {
		fmt.Println("❌ Error checking profile existence:", err)
		return false
	}
	return exists
}
