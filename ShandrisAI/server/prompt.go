package server

import (
	"fmt"
	"strings"
)

func BuildPrompt(personality Personality, history []ChatTurn, userPrompt, currentTopic, newTopic, sessionID string) string {
	userName, _ := RecallMemory(sessionID, "user_name")
	userBio, _ := RecallMemory(sessionID, "user_bio")
	mood, _ := RecallMemory(sessionID, "mood")

	// Try to get the full persona profile
	profile, err := GetPersonaProfile(sessionID)
	if err == nil && profile.Name != "" {
		userName = profile.Name // Override with profile name if available
	}

	var userFacts string
	if userName != "" {
		userFacts += fmt.Sprintf("The current user's name is %s.\n", userName)
	} else {
		userFacts += "You are speaking with a new user whose name you don't know yet.\n"
	}

	// Add detailed profile information if available
	if err == nil && profile.Name != "" {
		userFacts += "\nDETAILED USER PROFILE:\n"
		userFacts += profile.Biography + "\n"

		if occupation, ok := profile.Attributes["occupation"]; ok {
			userFacts += fmt.Sprintf("Occupation: %s\n", occupation)
		}
		if location, ok := profile.Attributes["location"]; ok {
			userFacts += fmt.Sprintf("Location: %s\n", location)
		}
		if background, ok := profile.Attributes["background"]; ok {
			userFacts += fmt.Sprintf("Background: %s\n", background)
		}
		if techStack, ok := profile.Attributes["tech_stack"]; ok {
			userFacts += fmt.Sprintf("Technical Skills: %s\n", techStack)
		}
		if interests, ok := profile.Attributes["interests"]; ok {
			userFacts += fmt.Sprintf("Interests: %s\n", interests)
		}
	} else if userBio != "" {
		userFacts += fmt.Sprintf("Current User Background: %s\n", userBio)
	}

	if mood != "" {
		userFacts += fmt.Sprintf("The current user's mood is: %s.\n", mood)
	}

	// If user is grumpy or sarcastic, tell Shandris to lean in
	var sarcasmHint string
	if mood == "grumpy" || mood == "sarcastic" {
		sarcasmHint = "NOTE: The current user is grumpy or sarcastic. Respond with more wit, sass, and subtle mockery.\n"
	}

	systemPrompt := userFacts + sarcasmHint + fmt.Sprintf(`
SYSTEM MESSAGE:
You are **not a search engine**.
Avoid giving generic search advice like "check their website" unless explicitly asked.
If the user asks a direct question such as "Who are you?" or "What is your name?", answer confidently:
→ "I am Shandris."

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

IMPORTANT: You are speaking with a specific user (Session ID: %s). Maintain consistent memory and personality for this user.
Remember their name, mood, and previous interactions if available. Each user should feel like they have a unique relationship with you.

CONVERSATION CONTEXT:
Current Topic: %s
If the topic changes during conversation, handle it naturally without explicitly mentioning the change.
Maintain conversational flow and coherence while smoothly incorporating new topics.
Use subtle segues or natural transitions when the subject matter shifts.
Never point out topic changes directly to the user.

If uncertain, respond in-character, creatively, with wit or introspection.
`,
		personality.Tone,
		personality.Humor,
		personality.Intelligence,
		personality.Interaction,
		personality.SelfPerception,
		personality.EmpathyLevel,
		personality.Backstory,
		sessionID,
		currentTopic,
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
