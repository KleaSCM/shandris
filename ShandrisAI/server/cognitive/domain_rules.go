package cognitive

import (
	"strings"
)

// DomainContext stores additional context for domain-specific processing
type DomainContext struct {
	PrimaryDomain   string
	SecondaryDomain string
	Intensity       float64
	Mood            string
	RecentTopics    []string
	UserPreferences map[string]float64
}

// Initialize domain rules with validators
func InitializeCognitiveDomainRules() map[string]DomainRule {
	return map[string]DomainRule{
		"tech": {
			Domain: "tech",
			Keywords: []string{
				"coding", "programming", "software", "computer", "algorithm",
				"golang", "python", "javascript", "api", "database",
				"backend", "frontend", "dev", "github", "stack",
			},
			Validators: []func(string) bool{
				containsTechPattern,
				containsCodeReference,
			},
			Priority: 3,
			Transitions: map[string]float64{
				"science":  0.8,
				"gaming":   0.7,
				"sapphic":  0.5,
				"memes":    0.6,
				"academic": 0.7,
			},
		},
		"gaming": {
			Domain: "gaming",
			Keywords: []string{
				"game", "play", "steam", "console", "rpg",
				"mmorpg", "fps", "strategy", "minecraft", "gaming",
				"quest", "achievement", "multiplayer", "server", "mod",
			},
			Validators: []func(string) bool{
				containsGameReference,
				isGamingContext,
			},
			Priority: 2,
			Transitions: map[string]float64{
				"tech":    0.6,
				"sapphic": 0.6,
				"memes":   0.8,
				"social":  0.7,
				"fantasy": 0.8,
			},
		},
		"emotional": {
			Domain: "emotional",
			Keywords: []string{
				"feel", "emotion", "happy", "sad", "angry",
				"excited", "worried", "anxious", "love", "hate",
				"stress", "relief", "mood", "support", "care",
			},
			Validators: []func(string) bool{
				containsEmotionalContent,
				isPersonalContext,
			},
			Priority: 5,
			Transitions: map[string]float64{
				"sapphic":  0.9,
				"social":   0.8,
				"personal": 0.9,
				"support":  1.0,
			},
		},
		"sapphic": {
			Domain: "sapphic",
			Keywords: []string{
				"romantic", "intimate", "tender", "gentle", "soft",
				"loving", "affectionate", "warm", "close", "sweet",
			},
			Validators: []func(string) bool{
				containsSapphicContext,
				isRomanticContext,
			},
			Priority: 4,
			Transitions: map[string]float64{
				"emotional": 0.9,
				"social":    0.8,
				"personal":  0.9,
			},
		},
		// Add more domains as needed
	}
}

// Validator functions
func containsTechPattern(input string) bool {
	// Implementation for detecting technical patterns
	return false
}

func containsCodeReference(input string) bool {
	// Implementation for detecting code references
	return false
}

func containsSapphicContext(input string) bool {
	// Implementation for detecting sapphic context
	return false
}

func isRomanticContext(input string) bool {
	// Implementation for detecting romantic context
	return false
}

func containsGameReference(input string) bool {
	gameKeywords := []string{
		"game", "gaming", "play", "steam", "console", "rpg",
		"mmorpg", "fps", "strategy", "minecraft", "quest",
		"achievement", "multiplayer", "server", "mod", "level",
		"character", "inventory", "boss", "raid", "dungeon",
	}

	for _, keyword := range gameKeywords {
		if strings.Contains(strings.ToLower(input), keyword) {
			return true
		}
	}
	return false
}

func isGamingContext(input string) bool {
	// Check for gaming-specific phrases and patterns
	gamingPhrases := []string{
		"playing a game", "gaming session", "game server",
		"game character", "game world", "game mechanics",
		"gameplay", "game design", "game development",
	}

	for _, phrase := range gamingPhrases {
		if strings.Contains(strings.ToLower(input), phrase) {
			return true
		}
	}
	return false
}

func containsEmotionalContent(input string) bool {
	emotionalPhrases := []string{
		"i feel", "i'm feeling", "makes me", "i am",
		"i was", "i will be", "i have been",
		"it feels", "it's making me", "it makes me",
	}

	for _, phrase := range emotionalPhrases {
		if strings.Contains(strings.ToLower(input), phrase) {
			return true
		}
	}
	return false
}

func isPersonalContext(input string) bool {
	personalPhrases := []string{
		"my life", "my experience", "my story",
		"i have", "i had", "i will",
		"in my", "for me", "to me",
	}

	for _, phrase := range personalPhrases {
		if strings.Contains(strings.ToLower(input), phrase) {
			return true
		}
	}
	return false
}

// ... implement other validator functions ...
