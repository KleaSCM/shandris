package cognitive

import (
	"sort"
	"time"
)

// ScoredMemory represents a memory with its recall score
type ScoredMemory struct {
	Event     *MemoryEvent
	Score     float64
	Relevance float64
	Recency   float64
	Emotion   float64
}

// getTopMemories selects the most relevant memories
func (tm *TimelineMemory) getTopMemories(scored []*ScoredMemory, limit int) []*MemoryEvent {
	// Sort by score
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	// Select top memories
	result := make([]*MemoryEvent, 0, limit)
	for i := 0; i < limit && i < len(scored); i++ {
		result = append(result, scored[i].Event)
	}

	return result
}

// Helper functions for relationship updates
func calculateTrustImpact(event *MemoryEvent) float64 {
	impact := 0.0

	switch event.Type {
	case PersonalEvent:
		impact = 0.05
	case EmotionalEvent:
		impact = 0.1
	case RelationshipEvent:
		impact = 0.15
	}

	// Modify based on emotional context
	if emotion, exists := event.Emotions["trust"]; exists {
		impact *= (1 + emotion)
	}

	return impact
}

func calculateIntimacyImpact(event *MemoryEvent) float64 {
	impact := 0.0

	switch event.Type {
	case PersonalEvent:
		impact = 0.1
	case EmotionalEvent:
		impact = 0.15
	case RelationshipEvent:
		impact = 0.2
	}

	// Modify based on emotional context
	if emotion, exists := event.Emotions["intimacy"]; exists {
		impact *= (1 + emotion)
	}

	return impact
}

// Sophisticated memory recall algorithms
type MemoryRecall struct {
	contextMapper  *ContextMapper
	emotionMatcher *EmotionMatcher
	patternMatcher *PatternMatcher
}

type ContextMapper struct {
	topicWeights map[string]float64
	moodWeights  map[string]float64
	timeWeights  map[string]float64
}

type EmotionMatcher struct {
	emotionPatterns    map[string][]float64
	resonanceThreshold float64
}

type PatternMatcher struct {
	patterns     map[string]*RecallPattern
	associations map[string][]string
}

type RecallPattern struct {
	Triggers   []string
	Weights    map[string]float64
	TimeWindow time.Duration
}

func (mr *MemoryRecall) FindRelevantMemories(context *EventContext, events map[string]*MemoryEvent) []*MemoryEvent {
	var relevant []*MemoryEvent

	// Context matching
	contextScores := mr.contextMapper.MapContext(context)

	// Emotion matching
	emotionScores := mr.emotionMatcher.MatchEmotions(context, events)

	// Pattern matching
	patternScores := mr.patternMatcher.MatchPatterns(context, events)

	// Combine scores and filter memories
	for id, event := range events {
		score := mr.calculateCombinedScore(
			contextScores[id],
			emotionScores[id],
			patternScores[id],
		)

		if score > 0.5 { // Threshold for relevance
			relevant = append(relevant, event)
		}
	}

	return relevant
}

func (cm *ContextMapper) MapContext(context *EventContext) map[string]float64 {
	scores := make(map[string]float64)

	// Topic matching
	for _, topic := range context.Topics {
		if weight, exists := cm.topicWeights[topic]; exists {
			scores[topic] = weight
		}
	}

	// Mood matching
	if weight, exists := cm.moodWeights[context.Mood]; exists {
		scores["mood"] = weight
	}

	// Time context matching
	timeContext := getTimeContext(time.Now())
	if weight, exists := cm.timeWeights[timeContext]; exists {
		scores["time"] = weight
	}

	return scores
}

func (em *EmotionMatcher) MatchEmotions(context *EventContext, events map[string]*MemoryEvent) map[string]float64 {
	scores := make(map[string]float64)

	for id, event := range events {
		resonance := calculateEmotionalResonance(event, context)
		if resonance > em.resonanceThreshold {
			scores[id] = resonance
		}
	}

	return scores
}
