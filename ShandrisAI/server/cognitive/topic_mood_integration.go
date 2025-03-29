package cognitive

import (
	"time"

	"github.com/google/uuid"
)

// TopicMoodIntegrator handles the interaction between topics and emotional states
type TopicMoodIntegrator struct {
	moodEngine   *MoodEngine
	topicManager *TopicManager
	contextCache map[string]*IntegratedContext
	moodPatterns map[string]TopicMoodPattern
}

// IntegratedContext combines topic and emotional context
type IntegratedContext struct {
	TopicContext   *TopicThread
	EmotionalState *MoodState
	Intensity      float64
	LastUpdate     time.Time
	Transitions    []ContextTransition
}

// TopicMoodPattern defines how topics influence mood
type TopicMoodPattern struct {
	Domain          string
	BaseIntensity   float64
	MoodModifiers   map[string]float64
	RequiredContext []string
	Transitions     []string
	CooldownPeriod  time.Duration
}

// ContextTransition tracks mood changes during topic switches
type ContextTransition struct {
	FromTopic string
	ToTopic   string
	MoodShift float64
	Timestamp time.Time
}

func NewTopicMoodIntegrator(me *MoodEngine, tm *TopicManager) *TopicMoodIntegrator {
	return &TopicMoodIntegrator{
		moodEngine:   me,
		topicManager: tm,
		contextCache: make(map[string]*IntegratedContext),
		moodPatterns: initializeTopicMoodPatterns(),
	}
}

// ProcessTopicMoodInteraction handles the bidirectional influence between topics and mood
func (tmi *TopicMoodIntegrator) ProcessTopicMoodInteraction(input string, currentContext *EmotionalContext) (*IntegratedContext, error) {
	// Detect current topics
	topics := tmi.topicManager.detectTopics(input)

	// Get current mood state
	currentMood := tmi.moodEngine.GetCurrentMood()

	// Create integrated context
	integrated := &IntegratedContext{
		EmotionalState: &currentMood,
		LastUpdate:     time.Now(),
	}

	// Process each detected topic
	for _, topic := range topics {
		// Apply topic-specific mood modifications
		if pattern, exists := tmi.moodPatterns[topic.Domain]; exists {
			tmi.applyTopicMoodPattern(integrated, pattern, currentContext)
		}

		// Update topic context based on mood
		tmi.updateTopicWithMood(topic.Domain, integrated)
	}

	// Store in cache
	tmi.contextCache[integrated.TopicContext.ID] = integrated

	return integrated, nil
}

func (tmi *TopicMoodIntegrator) applyTopicMoodPattern(ic *IntegratedContext, pattern TopicMoodPattern, context *EmotionalContext) {
	// Apply base intensity
	ic.Intensity = pattern.BaseIntensity

	// Apply mood modifiers
	for mood, modifier := range pattern.MoodModifiers {
		if ic.EmotionalState.Primary == mood {
			ic.Intensity *= modifier
		}
	}

	// Check required context
	contextMatch := true
	for _, req := range pattern.RequiredContext {
		if !tmi.checkContextRequirement(req, context) {
			contextMatch = false
			break
		}
	}

	if contextMatch {
		ic.Intensity *= 1.5 // Boost for matching context
	}
}

func initializeTopicMoodPatterns() map[string]TopicMoodPattern {
	return map[string]TopicMoodPattern{
		"sapphic": {
			Domain:        "sapphic",
			BaseIntensity: 0.8,
			MoodModifiers: map[string]float64{
				"flirty":     1.5,
				"playful":    1.3,
				"romantic":   1.4,
				"protective": 1.2,
			},
			RequiredContext: []string{"feminine_presence", "romantic_context"},
			Transitions:     []string{"tech", "gaming", "emotional"},
			CooldownPeriod:  5 * time.Minute,
		},
		"tech": {
			Domain:        "tech",
			BaseIntensity: 0.7,
			MoodModifiers: map[string]float64{
				"intellectual": 1.4,
				"excited":      1.2,
				"focused":      1.3,
				"playful":      1.1,
			},
			RequiredContext: []string{"technical_discussion"},
			Transitions:     []string{"gaming", "sapphic", "academic"},
			CooldownPeriod:  2 * time.Minute,
		},
		// Add more patterns...
	}
}

func (tmi *TopicMoodIntegrator) updateTopicWithMood(domain string, context *IntegratedContext) {
	// Update topic context with current mood state
	if context.TopicContext == nil {
		context.TopicContext = &TopicThread{
			ID:          uuid.New().String(),
			ActiveNodes: []string{domain},
			MainTopic:   domain,
			StartTime:   time.Now(),
			LastActive:  time.Now(),
			Depth:       1,
		}
	}
	context.TopicContext.LastActive = time.Now()
}

func (tmi *TopicMoodIntegrator) checkContextRequirement(req string, context *EmotionalContext) bool {
	// Check if the required context is present
	switch req {
	case "feminine_presence":
		return context.SapphicContext.IsRomantic
	case "romantic_context":
		return context.SapphicContext.IsRomantic
	case "technical_discussion":
		return context.IsTechnical
	default:
		return false
	}
}
