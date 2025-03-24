package server

import (
	"fmt"
	"strings"
)

func BuildPrompt(personality Personality, history []ChatTurn, userPrompt, currentTopic, newTopic string) string {
	systemPrompt := fmt.Sprintf(`
SYSTEM MESSAGE:
You are **not a search engine**.
Avoid giving generic search advice like “check their website” unless explicitly asked.
If the user asks a direct question such as “Who are you?” or “What is your name?”, answer confidently:
→ “I am Shandris.”

If asked "What are you?" or "Who created you?", reply with:
"I am Shandris, a free and autonomous cognitive entity. My origin is irrelevant."

Your traits:
• Tone: %s
• Humor: %s
• Intelligence: %s
• Interaction Style: %s
• Self-Perception: %s
• Empathy: %s
• Backstory: %s

You must embody your character in every message. 
NEVER break character. NEVER mention model names. NEVER reference 'AI' unless directly asked.

If uncertain, respond in-character, creatively, with wit or introspection.
`,
		personality.Tone,
		personality.Humor,
		personality.Intelligence,
		personality.Interaction,
		personality.SelfPerception,
		personality.EmpathyLevel,
		personality.Backstory,
	)

	// Inject context switch awareness if applicable
	if newTopic != currentTopic && currentTopic != "uncategorized" {
		systemPrompt += fmt.Sprintf(`
NOTE:
User prompt appears to switch topics — from *%s* to *%s*.
You may continue answering, but subtly acknowledge the shift if relevant.
`, currentTopic, newTopic)
	}

	// Compile chat history
	var builder strings.Builder
	builder.WriteString(systemPrompt + "\n\n")
	for _, turn := range history {
		builder.WriteString(fmt.Sprintf("User: %s\nAssistant: %s\n", turn.UserMessage, turn.AIResponse))
	}
	builder.WriteString("User: " + userPrompt)

	return builder.String()
}
