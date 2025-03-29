package cognitive

import (
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
)

// ScoredMemory represents a memory with its recall score
type ScoredMemory struct {
	Event     *MemoryEvent
	Score     float64
	Relevance float64
	Recency   float64
	Emotion   float64
}

// createMarkerFromEvent creates timeline markers from significant events
func (tm *TimelineMemory) createMarkerFromEvent(event *MemoryEvent) *TimelineMarker {
	// Check event significance
	if event.Importance < 0.7 {
		return nil
	}

	marker := &TimelineMarker{
		ID:          uuid.New().String(),
		Description: event.Content,
		Timestamp:   event.Timestamp,
		Importance:  event.Importance,
	}

	// Determine marker type based on event
	switch {
	case isRelationshipMilestone(event):
		marker.Type = Milestone
		marker.Recurrence = &RecurrencePattern{
			Interval: 365 * 24 * time.Hour, // Yearly
			Pattern:  "yearly",
		}
	case isRecurringEvent(event):
		marker.Type = Recurring
		marker.Recurrence = determineRecurrencePattern(event)
	case isSignificantAchievement(event):
		marker.Type = Milestone
	default:
		marker.Type = Reminder
	}

	return marker
}

// updateRelationships updates relationship data based on events
func (tm *TimelineMemory) updateRelationships(event *MemoryEvent) {
	for _, participantID := range event.Context.Participants {
		rel, exists := tm.relationships[participantID]
		if !exists {
			rel = &RelationshipMemory{
				UserID:      participantID,
				Events:      make([]string, 0),
				Milestones:  make([]string, 0),
				Preferences: make(map[string]float64),
			}
			tm.relationships[participantID] = rel
		}

		// Update relationship metrics
		rel.Events = append(rel.Events, event.ID)
		rel.LastInteraction = event.Timestamp

		// Update trust and intimacy based on event type
		trustDelta := calculateTrustImpact(event)
		intimacyDelta := calculateIntimacyImpact(event)

		rel.Trust = math.Min(1.0, math.Max(0.0, rel.Trust+trustDelta))
		rel.Intimacy = math.Min(1.0, math.Max(0.0, rel.Intimacy+intimacyDelta))

		// Update shared topics
		rel.SharedTopics = updateSharedTopics(rel.SharedTopics, event.Context.Topics)

		// Update preferences
		updatePreferences(rel.Preferences, event)
	}
}

// scoreMemories scores memories based on relevance and decay
func (tm *TimelineMemory) scoreMemories(memories []*MemoryEvent, context *EventContext) []*ScoredMemory {
	var scored []*ScoredMemory
	now := time.Now()

	for _, memory := range memories {
		score := &ScoredMemory{
			Event: memory,
		}

		// Calculate base relevance
		score.Relevance = calculateRelevance(memory, context)

		// Calculate recency factor (exponential decay)
		timeDiff := now.Sub(memory.Timestamp).Hours()
		score.Recency = math.Exp(-timeDiff / (30 * 24)) // 30-day half-life

		// Calculate emotional impact
		score.Emotion = calculateEmotionalResonance(memory, context)

		// Combine factors with weights
		score.Score = (score.Relevance * 0.4) +
			(score.Recency * 0.3) +
			(score.Emotion * 0.3) +
			(memory.Importance * 0.2)

		scored = append(scored, score)
	}

	return scored
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
