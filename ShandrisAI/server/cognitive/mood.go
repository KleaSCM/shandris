package cognitive

import (
	"math"
	"strings"
	"time"
)

// MoodEngine implementation
type MoodEngineImpl struct {
	CurrentState MoodState
	History      []MoodState
	Modifiers    map[string]float64
	Thresholds   map[string]float64
}

func NewMoodEngine() *MoodEngineImpl {
	return &MoodEngineImpl{
		CurrentState: MoodState{
			Primary:   "neutral",
			Intensity: 0.5,
			Timestamp: time.Now(),
			Context:   make(map[string]any),
		},
		Modifiers: map[string]float64{
			"decay_rate":    0.1,
			"change_thresh": 0.3,
			"max_intensity": 1.0,
		},
		Thresholds: map[string]float64{
			"mood_shift":    0.4,
			"intensity_cap": 0.9,
		},
	}
}

func (m *MoodEngineImpl) UpdateMood(context map[string]any) error {
	// Calculate time-based decay
	timeDiff := time.Since(m.CurrentState.Timestamp).Hours()
	decayFactor := math.Exp(-m.Modifiers["decay_rate"] * timeDiff)

	// Apply context modifiers
	newIntensity := m.CurrentState.Intensity * decayFactor
	for _, impact := range context {
		if val, ok := impact.(float64); ok {
			newIntensity += val * m.Modifiers["change_thresh"]
		}
	}

	// Clamp intensity
	newIntensity = math.Max(0, math.Min(newIntensity, m.Modifiers["max_intensity"]))

	// Record history
	m.History = append(m.History, m.CurrentState)

	// Update current state
	m.CurrentState = MoodState{
		Primary:   m.determinePrimaryMood(context),
		Intensity: newIntensity,
		Timestamp: time.Now(),
		Context:   context,
	}

	return nil
}

func (m *MoodEngineImpl) determinePrimaryMood(context map[string]any) string {
	// Define base mood patterns
	moodPatterns := map[string]MoodPattern{
		"playful": {
			Keywords:  []string{"haha", "lol", "😂", "fun", "play", "joke", "tease"},
			Sentiment: 0.7,
			MoodShift: "playful",
			Intensity: 0.6,
			Decay:     0.2,
		},
		"sassy": {
			Keywords:  []string{"actually", "well_actually", "oh_really", "sure_jan", "whatever"},
			Sentiment: 0.3,
			MoodShift: "sassy",
			Intensity: 0.8,
			Decay:     0.1,
		},
		"flirty": {
			Keywords:  []string{"cute", "pretty", "beautiful", "hot", "gorgeous", "flirt"},
			Sentiment: 0.8,
			MoodShift: "flirty",
			Intensity: 0.7,
			Decay:     0.15,
		},
		"intellectual": {
			Keywords:  []string{"think", "theory", "quantum", "algorithm", "complex", "interesting"},
			Sentiment: 0.4,
			MoodShift: "intellectual",
			Intensity: 0.6,
			Decay:     0.05,
		},
		"protective": {
			Keywords:  []string{"help", "protect", "safe", "careful", "worry", "concern"},
			Sentiment: 0.2,
			MoodShift: "protective",
			Intensity: 0.5,
			Decay:     0.3,
		},
	}

	// Extract emotional context from the input
	emotionalContext := m.analyzeEmotionalContext(context)

	// Calculate mood scores based on patterns and context
	moodScores := make(map[string]float64)
	for mood, pattern := range moodPatterns {
		score := m.calculateMoodScore(pattern, emotionalContext)
		moodScores[mood] = score
	}

	// Apply personality biases
	m.applyPersonalityBias(moodScores)

	// Find dominant mood
	return m.selectDominantMood(moodScores, emotionalContext)
}

func (m *MoodEngineImpl) calculateMoodScore(pattern MoodPattern, context EmotionalContext) float64 {
	// Base score from keyword matches
	keywordScore := 0.0
	for _, keyword := range pattern.Keywords {
		if containsKeyword(context.Keywords, keyword) {
			keywordScore += 0.2
		}
	}

	// Sentiment alignment
	sentimentAlignment := 1.0 - math.Abs(pattern.Sentiment-context.Sentiment)

	// Intensity factor
	intensityFactor := math.Min(context.Intensity*pattern.Intensity, 1.0)

	// Decay from previous mood if applicable
	decayFactor := 1.0
	if context.PriorContext != nil {
		timeDiff := time.Since(m.CurrentState.Timestamp).Hours()
		decayFactor = math.Exp(-pattern.Decay * timeDiff)
	}

	// Combine factors
	return (keywordScore*0.4 + sentimentAlignment*0.3 + intensityFactor*0.3) * decayFactor
}

func (m *MoodEngineImpl) applyPersonalityBias(scores map[string]float64) {
	// Personality-based mood biases
	biases := map[string]float64{
		"sassy":        0.2,  // Shandris tends toward sass
		"intellectual": 0.15, // Appreciates intellectual discussion
		"flirty":       0.1,  // Slight flirty bias
	}

	// Apply biases
	for mood, bias := range biases {
		if score, exists := scores[mood]; exists {
			scores[mood] = math.Min(score+bias, 1.0)
		}
	}
}

func (m *MoodEngineImpl) selectDominantMood(scores map[string]float64, context EmotionalContext) string {
	// Find highest scoring mood
	highestScore := -1.0
	dominantMood := "neutral"

	// Adjust scores based on context
	for mood, score := range scores {
		// Boost score if it matches user's mood
		if mood == context.UserMood {
			score *= 1.2
		}
		// Reduce score if it conflicts with current emotional tone
		if mood == "flirty" && !context.SapphicContext.AllowsFlirting {
			score *= 0.5
		}
		if score > highestScore {
			highestScore = score
			dominantMood = mood
		}
	}

	// If no strong mood, maintain current with decay
	if highestScore < m.Thresholds["mood_shift"] {
		if m.CurrentState.Intensity > m.Thresholds["mood_shift"] {
			return m.CurrentState.Primary
		}
		return "neutral"
	}

	return dominantMood
}

func (m *MoodEngineImpl) analyzeEmotionalContext(context map[string]any) EmotionalContext {
	// Extract relevant information from context
	keywords := extractKeywords(context)
	sentiment := calculateSentiment(context)
	intensity := calculateIntensity(context)
	userMood := extractUserMood(context)

	return EmotionalContext{
		PrimaryEmotion: m.CurrentState.Primary,
		Sentiment:      sentiment,
		Intensity:      intensity,
		Keywords:       keywords,
		UserMood:       userMood,
		PriorContext:   m.getPriorContext(),
		SapphicContext: m.getCurrentSapphicContext(),
		IsEmotional:    true,
		IsTechnical:    false,
		Timestamp:      time.Now(),
		RawInput:       "",
		ThemeScores:    make(map[string]float64),
		UserContext:    context,
		EmotionalTone:  m.getCurrentEmotionalTone(),
	}
}

// Helper functions (to be implemented based on your specific needs)
func containsKeyword(keywords []string, target string) bool {
	for _, k := range keywords {
		if strings.Contains(strings.ToLower(k), strings.ToLower(target)) {
			return true
		}
	}
	return false
}

func extractKeywords(context map[string]any) []string {
	keywords := make([]string, 0)

	// Extract keywords from context map
	if text, ok := context["text"].(string); ok {
		keywords = append(keywords, strings.Fields(strings.ToLower(text))...)
	}
	if topics, ok := context["topics"].([]string); ok {
		keywords = append(keywords, topics...)
	}

	return keywords
}

func calculateSentiment(context map[string]any) float64 {
	sentiment := 0.0

	// Extract sentiment from context
	if val, ok := context["sentiment"].(float64); ok {
		sentiment = val
	}
	if text, ok := context["text"].(string); ok {
		// Basic sentiment analysis based on text content
		if strings.Contains(strings.ToLower(text), "love") || strings.Contains(strings.ToLower(text), "happy") {
			sentiment += 0.3
		}
		if strings.Contains(strings.ToLower(text), "hate") || strings.Contains(strings.ToLower(text), "sad") {
			sentiment -= 0.3
		}
	}

	return math.Max(-1.0, math.Min(1.0, sentiment))
}

func calculateIntensity(context map[string]any) float64 {
	intensity := 0.5 // Base intensity

	// Extract intensity from context
	if val, ok := context["intensity"].(float64); ok {
		intensity = val
	}
	if text, ok := context["text"].(string); ok {
		// Adjust intensity based on text content
		if strings.Contains(strings.ToLower(text), "very") || strings.Contains(strings.ToLower(text), "really") {
			intensity += 0.2
		}
	}

	return math.Min(1.0, intensity)
}

func extractUserMood(context map[string]any) string {
	// Extract user mood from context
	if mood, ok := context["user_mood"].(string); ok {
		return mood
	}
	if text, ok := context["text"].(string); ok {
		// Basic mood detection from text
		if strings.Contains(strings.ToLower(text), "happy") || strings.Contains(strings.ToLower(text), "excited") {
			return "happy"
		}
		if strings.Contains(strings.ToLower(text), "sad") || strings.Contains(strings.ToLower(text), "angry") {
			return "negative"
		}
	}
	return "neutral"
}

func (m *MoodEngineImpl) getPriorContext() *EmotionalContext {
	if len(m.History) == 0 {
		return nil
	}
	// Convert last mood state to emotional context
	// Implementation depends on your specific needs
	return nil
}

func (m *MoodEngineImpl) GetCurrentMood() MoodState {
	return m.CurrentState
}

// ProcessInteraction processes an interaction and updates the emotional context
func (m *MoodEngineImpl) ProcessInteraction(interaction *Interaction, context *EmotionalContext) *MoodUpdate {
	// Update mood based on interaction and context
	m.UpdateMood(context.UserContext)
	return &MoodUpdate{
		NewState: &m.CurrentState,
	}
}

func (m *MoodEngineImpl) ProjectMoodInfluence(input string) float64 {
	// Implement mood projection logic
	return 0.0 // Placeholder
}

func (m *MoodEngineImpl) getCurrentSapphicContext() SapphicContext {
	// Implementation depends on your specific needs
	return SapphicContext{}
}

func (m *MoodEngineImpl) getCurrentEmotionalTone() string {
	// Implementation depends on your specific needs
	return "neutral"
}
