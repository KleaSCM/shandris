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

	for mood, score := range scores {
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
	// Implementation depends on your context structure
	return []string{}
}

func calculateSentiment(context map[string]any) float64 {
	// Implementation depends on your sentiment analysis needs
	return 0.0
}

func calculateIntensity(context map[string]any) float64 {
	// Implementation depends on your intensity calculation needs
	return 0.0
}

func extractUserMood(context map[string]any) string {
	// Implementation depends on your user mood detection
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
