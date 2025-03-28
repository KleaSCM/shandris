package server

import "strings"

func ClassifyPrompt(prompt string) string {
	p := strings.ToLower(prompt)

	// General conversation markers
	if strings.Contains(p, "hello") || strings.Contains(p, "hi ") || strings.Contains(p, "hey") ||
		strings.Contains(p, "good morning") || strings.Contains(p, "good evening") || strings.Contains(p, "good afternoon") {
		return "greeting"
	}

	// Personal/emotional content
	if strings.Contains(p, "feel") || strings.Contains(p, "emotion") || strings.Contains(p, "sad") ||
		strings.Contains(p, "happy") || strings.Contains(p, "angry") || strings.Contains(p, "tired") ||
		strings.Contains(p, "love") || strings.Contains(p, "hate") || strings.Contains(p, "miss") {
		return "emotional"
	}

	// Identity and personal questions
	if strings.Contains(p, "who are you") || strings.Contains(p, "what are you") ||
		strings.Contains(p, "your name") || strings.Contains(p, "about you") ||
		strings.Contains(p, "tell me about yourself") {
		return "identity"
	}

	// Philosophy and deep thoughts
	if strings.Contains(p, "think") || strings.Contains(p, "believe") || strings.Contains(p, "opinion") ||
		strings.Contains(p, "meaning") || strings.Contains(p, "purpose") || strings.Contains(p, "life") ||
		strings.Contains(p, "death") || strings.Contains(p, "exist") || strings.Contains(p, "consciousness") {
		return "philosophy"
	}

	// Knowledge and learning
	if strings.Contains(p, "know") || strings.Contains(p, "learn") || strings.Contains(p, "teach") ||
		strings.Contains(p, "explain") || strings.Contains(p, "how does") || strings.Contains(p, "what is") ||
		strings.Contains(p, "why does") {
		return "knowledge"
	}

	// Casual conversation
	if strings.Contains(p, "like") || strings.Contains(p, "fun") || strings.Contains(p, "interesting") ||
		strings.Contains(p, "cool") || strings.Contains(p, "awesome") || strings.Contains(p, "nice") {
		return "casual"
	}

	// Keep the current topic if we can't clearly classify
	return "uncategorized"
}
