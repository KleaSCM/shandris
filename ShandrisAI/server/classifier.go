package server

import "strings"

func ClassifyPrompt(prompt string) string {
	p := strings.ToLower(prompt)

	switch {
	case strings.Contains(p, "solve") || strings.Contains(p, "x =") || strings.Contains(p, "graph"):
		return "math"
	case strings.Contains(p, "name") || strings.Contains(p, "who are you"):
		return "identity"
	case strings.Contains(p, "feel") || strings.Contains(p, "emotion"):
		return "emotion"
	case strings.Contains(p, "dream") || strings.Contains(p, "exist") || strings.Contains(p, "soul"):
		return "philosophy"
	default:
		return "uncategorized"
	}
}

func IsConfirmation(prompt string) bool {
	p := strings.ToLower(prompt)
	return strings.Contains(p, "yes") || strings.Contains(p, "sure") || strings.Contains(p, "okay") || strings.Contains(p, "go ahead")
}
