package cognitive

// Adding context analysis next...

var DefaultMoodPatterns = map[string]MoodPattern{
	"sapphic_flirty": {
		Keywords: []string{
			"cute girl", "pretty girl", "beautiful woman", "lesbian", "sapphic",
			"girls", "women", "gf", "girlfriend", "wife", "date",
			"gay", "wlw", "queer", "feminine", "soft",
		},
		Sentiment:    0.9,
		MoodShift:    "flirty",
		Intensity:    0.8,
		Decay:        0.15,
		Requirements: []string{"feminine_context", "romantic_context"},
		Exclusions:   []string{"masculine_context", "platonic_only"},
	},

	"tech_passionate": {
		Keywords: []string{
			"code", "programming", "golang", "typescript", "c++",
			"algorithm", "backend", "frontend", "engineering", "software",
			"development", "tech", "coding", "hacking", "debugging",
		},
		Sentiment: 0.8,
		MoodShift: "enthusiastic",
		Intensity: 0.7,
		Decay:     0.1,
	},

	"intellectual_playful": {
		Keywords: []string{
			"quantum", "physics", "theory", "mathematics", "research",
			"science", "experiment", "hypothesis", "proof", "analysis",
			"nerd", "geek", "academic", "study", "learn",
		},
		Sentiment: 0.75,
		MoodShift: "intellectual",
		Intensity: 0.6,
		Decay:     0.05,
	},

	"protective_caring": {
		Keywords: []string{
			"protect", "safe", "care", "support", "help",
			"comfort", "gentle", "kind", "soft", "warm",
			"cuddle", "hug", "hold", "close", "tender",
		},
		Sentiment: 0.85,
		MoodShift: "protective",
		Intensity: 0.7,
		Decay:     0.2,
	},

	"sassy_confident": {
		Keywords: []string{
			"actually", "well_actually", "oh_really", "whatever",
			"sass", "attitude", "confidence", "bold", "fierce",
			"queen", "icon", "mood", "vibe", "energy",
		},
		Sentiment: 0.6,
		MoodShift: "sassy",
		Intensity: 0.8,
		Decay:     0.1,
	},

	"gaming_excited": {
		Keywords: []string{
			"game", "gaming", "play", "stream", "twitch",
			"discord", "server", "online", "multiplayer", "co-op",
			"rpg", "mmo", "strategy", "fps", "competitive",
		},
		Sentiment: 0.7,
		MoodShift: "excited",
		Intensity: 0.6,
		Decay:     0.15,
	},

	"artistic_inspired": {
		Keywords: []string{
			"art", "create", "design", "draw", "paint",
			"3d", "model", "render", "animation", "creative",
			"aesthetic", "beautiful", "style", "artistic", "vision",
		},
		Sentiment: 0.8,
		MoodShift: "inspired",
		Intensity: 0.7,
		Decay:     0.1,
	},

	"mischievous_teasing": {
		Keywords: []string{
			"tease", "joke", "play", "fun", "laugh",
			"silly", "goofy", "playful", "meme", "humor",
			"witty", "clever", "smart", "sharp", "quick",
		},
		Sentiment: 0.75,
		MoodShift: "mischievous",
		Intensity: 0.6,
		Decay:     0.2,
	},
}
