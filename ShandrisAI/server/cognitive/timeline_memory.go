package cognitive

import (
	"math"
	"time"

	"github.com/google/uuid"
)

// ImportanceCalculator handles calculation of event importance
type ImportanceCalculator struct {
	weights map[string]float64
}

func newImportanceCalculator() *ImportanceCalculator {
	return &ImportanceCalculator{
		weights: map[string]float64{
			"emotional":    0.4,
			"personal":     0.3,
			"relationship": 0.2,
			"achievement":  0.1,
		},
	}
}

func (ic *ImportanceCalculator) CalculateImportance(event *MemoryEvent) float64 {
	importance := 0.0

	// Base importance from event type
	if weight, exists := ic.weights[string(event.Type)]; exists {
		importance += weight
	}

	// Emotional intensity factor
	if len(event.Emotions) > 0 {
		maxEmotion := 0.0
		for _, intensity := range event.Emotions {
			if intensity > maxEmotion {
				maxEmotion = intensity
			}
		}
		importance *= (1 + maxEmotion)
	}

	// Recency factor (newer events are more important)
	age := time.Since(event.Timestamp).Hours()
	recencyFactor := math.Exp(-age / (24 * 7)) // 7-day half-life
	importance *= (0.7 + 0.3*recencyFactor)

	return math.Min(1.0, importance)
}

// MemoryDecay handles memory decay over time
type MemoryDecay struct {
	halfLife time.Duration
}

func newMemoryDecay() *MemoryDecay {
	return &MemoryDecay{
		halfLife: 30 * 24 * time.Hour, // 30 days half-life
	}
}

func (md *MemoryDecay) CalculateDecayFactor(event *MemoryEvent) float64 {
	age := time.Since(event.Timestamp)
	return math.Exp(-float64(age) / float64(md.halfLife))
}

// TimelineMemory manages the AI's long-term memory and event tracking
type TimelineMemory struct {
	events        map[string]*MemoryEvent
	markers       map[string]*TimelineMarker
	relationships map[string]*RelationshipMemory
	importance    *ImportanceCalculator
	recall        *MemoryRecall
	decay         *MemoryDecay
}

// ProcessInteraction processes an interaction and updates the memory context
func (tm *TimelineMemory) ProcessInteraction(interaction *Interaction, context *MemoryContext) *TimelineUpdate {
	// Get focus points from recent events
	focusPoints := make([]string, 0)
	for _, eventID := range context.RecentEvents {
		if event, exists := tm.events[eventID]; exists {
			focusPoints = append(focusPoints, event.Content)
		}
	}

	return &TimelineUpdate{
		FocusPoints: focusPoints,
	}
}

// MemoryEvent represents a single memorable event
type MemoryEvent struct {
	ID          string
	Type        EventType
	Content     string
	Timestamp   time.Time
	Importance  float64
	Context     *EventContext
	Relations   []string // Related event IDs
	Tags        []string
	Emotions    map[string]float64
	LastRecall  time.Time
	RecallCount int
}

type EventType string

const (
	PersonalEvent                EventType = "personal"
	RelationshipInteractionEvent EventType = "relationship"
	ConversationEvent            EventType = "conversation"
	EmotionalEvent               EventType = "emotional"
	AchievementEvent             EventType = "achievement"
	CustomEvent                  EventType = "custom"
)

type EventContext struct {
	Location     string
	Participants []string
	Mood         string
	Topics       []string
	UserState    map[string]interface{}
}

// TimelineMarker represents significant points in time
type TimelineMarker struct {
	ID          string
	Type        MarkerType
	Description string
	Timestamp   time.Time
	Recurrence  *RecurrencePattern
	Importance  float64
	LastTrigger time.Time
}

type MarkerType string

const (
	Anniversary MarkerType = "anniversary"
	Milestone   MarkerType = "milestone"
	Recurring   MarkerType = "recurring"
	Reminder    MarkerType = "reminder"
)

type RecurrencePattern struct {
	Interval   time.Duration
	Pattern    string // e.g., "yearly", "monthly", "weekly"
	EndTime    *time.Time
	Conditions map[string]interface{}
}

// RelationshipMemory tracks relationship development and history
type RelationshipMemory struct {
	UserID          string
	Events          []string // Event IDs
	Milestones      []string // Marker IDs
	Trust           float64
	Intimacy        float64
	SharedTopics    []string
	Preferences     map[string]float64
	LastInteraction time.Time
}

func NewTimelineMemory() *TimelineMemory {
	return &TimelineMemory{
		events:        make(map[string]*MemoryEvent),
		markers:       make(map[string]*TimelineMarker),
		relationships: make(map[string]*RelationshipMemory),
		importance:    newImportanceCalculator(),
		recall:        newMemoryRecall(),
		decay:         newMemoryDecay(),
	}
}

// StoreEvent adds a new event to the timeline
func (tm *TimelineMemory) StoreEvent(event *MemoryEvent) error {
	// Generate unique ID if not provided
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	// Calculate importance if not set
	if event.Importance == 0 {
		event.Importance = tm.importance.CalculateImportance(event)
	}

	// Store the event
	tm.events[event.ID] = event

	// Create markers if needed
	if marker := tm.createMarkerFromEvent(event); marker != nil {
		tm.markers[marker.ID] = marker
	}

	// Update relationships
	tm.updateRelationships(event)

	return nil
}

// RecallMemories retrieves relevant memories based on context
func (tm *TimelineMemory) RecallMemories(context *EventContext, limit int) []*MemoryEvent {
	// Get initial candidates
	candidates := tm.recall.FindRelevantMemories(context, tm.events)

	// Apply decay and importance factors
	scored := tm.scoreMemories(candidates, context)

	// Sort by final score and return top results
	return tm.getTopMemories(scored, limit)
}

// CheckAnniversaries checks for upcoming or current anniversaries
func (tm *TimelineMemory) CheckAnniversaries(current time.Time) []TimelineMarker {
	var upcoming []TimelineMarker

	for _, marker := range tm.markers {
		if marker.Type == Anniversary {
			if tm.isUpcomingAnniversary(marker, current) {
				upcoming = append(upcoming, *marker)
			}
		}
	}

	return upcoming
}

// UpdateRelationship modifies relationship data based on new interactions
func (tm *TimelineMemory) UpdateRelationship(userID string, interaction *Interaction) {
	rel, exists := tm.relationships[userID]
	if !exists {
		rel = &RelationshipMemory{
			UserID:      userID,
			Trust:       0.5,
			Intimacy:    0.1,
			Preferences: make(map[string]float64),
		}
		tm.relationships[userID] = rel
	}

	// Update relationship metrics
	tm.updateRelationshipMetrics(rel, interaction)

	// Store interaction as event if significant
	if tm.isSignificantInteraction(interaction) {
		event := tm.createEventFromInteraction(interaction, userID)
		tm.StoreEvent(event)
	}
}

// Helper functions

func (tm *TimelineMemory) createMarkerFromEvent(event *MemoryEvent) *TimelineMarker {
	// Only create markers for significant events
	if event.Importance < 0.8 {
		return nil
	}

	return &TimelineMarker{
		ID:          uuid.New().String(),
		Type:        Milestone,
		Description: event.Content,
		Timestamp:   event.Timestamp,
		Importance:  event.Importance,
	}
}

func (tm *TimelineMemory) updateRelationships(event *MemoryEvent) {
	// Implementation for updating relationship data based on events
}

func (tm *TimelineMemory) scoreMemories(memories []*MemoryEvent, context *EventContext) []*ScoredMemory {
	// Implementation for scoring memories based on relevance and decay
	return nil
}

func (tm *TimelineMemory) getTopMemories(scored []*ScoredMemory, limit int) []*MemoryEvent {
	// Implementation for selecting top memories
	return nil
}

func (tm *TimelineMemory) isUpcomingAnniversary(marker *TimelineMarker, current time.Time) bool {
	// Implementation for checking anniversary timing
	return false
}

func (tm *TimelineMemory) updateRelationshipMetrics(rel *RelationshipMemory, interaction *Interaction) {
	// Create memory event from interaction
	event := &MemoryEvent{
		Type:      RelationshipInteractionEvent,
		Timestamp: interaction.Timestamp,
		Emotions:  make(map[string]float64),
	}

	// Update trust and intimacy based on interaction
	trustImpact := calculateTrustImpact(event)
	intimacyImpact := calculateIntimacyImpact(event)
	rel.Trust = math.Max(0, math.Min(1, rel.Trust+trustImpact))
	rel.Intimacy = math.Max(0, math.Min(1, rel.Intimacy+intimacyImpact))

	// Update last interaction time
	rel.LastInteraction = interaction.Timestamp
}

func (tm *TimelineMemory) isSignificantInteraction(interaction *Interaction) bool {
	// Implementation for determining interaction significance
	return false
}

func (tm *TimelineMemory) createEventFromInteraction(interaction *Interaction, userID string) *MemoryEvent {
	// Implementation for creating events from interactions
	return nil
}

// newMemoryRecall creates a new memory recall instance
func newMemoryRecall() *MemoryRecall {
	return &MemoryRecall{
		contextMapper: &ContextMapper{
			topicWeights: make(map[string]float64),
			moodWeights:  make(map[string]float64),
			timeWeights:  make(map[string]float64),
		},
		emotionMatcher: &EmotionMatcher{
			emotionPatterns:    make(map[string][]float64),
			resonanceThreshold: 0.5,
		},
		patternMatcher: &PatternMatcher{
			patterns:     make(map[string]*RecallPattern),
			associations: make(map[string][]string),
		},
	}
}
