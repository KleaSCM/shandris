package server

import "strings"

func ClassifyPrompt(prompt string) string {
	p := strings.ToLower(prompt)

	switch {
	// Math
	case strings.Contains(p, "solve") || strings.Contains(p, "x =") ||
		strings.Contains(p, "graph") || strings.Contains(p, "equation") ||
		strings.Contains(p, "derivative") || strings.Contains(p, "integral"):
		return "math"

	// Identity
	case strings.Contains(p, "name") || strings.Contains(p, "who are you") ||
		strings.Contains(p, "what is your name"):
		return "identity"

	// Emotion
	case strings.Contains(p, "feel") || strings.Contains(p, "emotion"):
		return "emotion"

	// Philosophy
	case strings.Contains(p, "dream") || strings.Contains(p, "exist") ||
		strings.Contains(p, "soul") || strings.Contains(p, "believe"):
		return "philosophy"

	// Business-ish or quantity prompts
	case strings.Contains(p, "how many") || strings.Contains(p, "how much"):
		return "business"

	default:
		return "uncategorized"
	}
}
