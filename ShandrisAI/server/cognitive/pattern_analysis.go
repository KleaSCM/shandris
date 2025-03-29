package cognitive

import (
	"sort"
	"time"
)

// PatternAnalysisEngine handles sophisticated pattern detection
type PatternAnalysisEngine struct {
	topicAnalyzer  *TopicAnalyzer
	persistence    *TopicPersistence
	moodIntegrator *TopicMoodIntegrator
	patterns       map[string]AnalysisPattern
	contextHistory []ContextSnapshot
	maxHistorySize int
}

type AnalysisPattern struct {
	ID         string
	Type       PatternType
	Triggers   []PatternTrigger
	Conditions []PatternCondition
	Weight     float64
	Decay      float64
	MinMatches int
	MaxGap     time.Duration
}

type PatternType string

const (
	BehavioralPattern PatternType = "behavioral"
	EmotionalPattern  PatternType = "emotional"
	TopicalPattern    PatternType = "topical"
	ContextualPattern PatternType = "contextual"
	HybridPattern     PatternType = "hybrid"
)

type PatternTrigger struct {
	Type       string
	Conditions map[string]interface{}
	Weight     float64
	Cooldown   time.Duration
}

type PatternCondition struct {
	Field     string
	Operator  string
	Value     interface{}
	Threshold float64
}

type ContextSnapshot struct {
	Timestamp    time.Time
	Topics       []string
	Mood         *MoodState
	UserContext  map[string]interface{}
	Interactions []Interaction
}

type Interaction struct {
	Type      string
	Value     interface{}
	Intensity float64
	Timestamp time.Time
}

type PatternResult struct {
	Pattern    *AnalysisPattern
	Confidence float64
	Context    *ContextSnapshot
}

func (pa *PatternAnalysisEngine) updateContextHistory(context *ContextSnapshot) {
	pa.contextHistory = append(pa.contextHistory, *context)
	if len(pa.contextHistory) > pa.maxHistorySize {
		pa.contextHistory = pa.contextHistory[1:]
	}
}

func (pa *PatternAnalysisEngine) AnalyzePatterns(context *ContextSnapshot) []PatternResult {
	var results []PatternResult

	// Update context history
	pa.updateContextHistory(context)

	// Analyze each pattern type
	results = append(results, pa.analyzeBehavioralPatterns()...)
	results = append(results, pa.analyzeEmotionalPatterns()...)
	results = append(results, pa.analyzeTopicalPatterns()...)
	results = append(results, pa.analyzeContextualPatterns()...)
	results = append(results, pa.analyzeHybridPatterns()...)

	// Sort and filter results
	return pa.filterAndRankPatterns(results)
}

func (pa *PatternAnalysisEngine) analyzeBehavioralPatterns() []PatternResult {
	var results []PatternResult

	// Analyze interaction sequences
	sequences := pa.detectInteractionSequences()
	for _, seq := range sequences {
		if pattern := pa.matchBehavioralPattern(seq); pattern != nil {
			results = append(results, PatternResult{
				Pattern:    pattern,
				Confidence: pa.calculatePatternConfidence(pattern, seq),
				Context:    pa.extractPatternContext(seq),
			})
		}
	}

	// Analyze timing patterns
	timingPatterns := pa.detectTimingPatterns()
	for _, timing := range timingPatterns {
		results = append(results, pa.analyzeTiming(timing))
	}

	return results
}

func (pa *PatternAnalysisEngine) analyzeEmotionalPatterns() []PatternResult {
	var results []PatternResult

	// Analyze mood transitions
	transitions := pa.detectMoodTransitions()
	for _, transition := range transitions {
		results = append(results, pa.analyzeMoodTransition(transition))
	}

	// Analyze emotional triggers
	triggers := pa.detectEmotionalTriggers()
	for _, trigger := range triggers {
		results = append(results, pa.analyzeEmotionalTrigger(trigger))
	}

	// Analyze emotional resonance
	resonance := pa.calculateEmotionalResonance()
	results = append(results, pa.analyzeEmotionalResonance(resonance))

	return results
}

func (pa *PatternAnalysisEngine) analyzeTopicalPatterns() []PatternResult {
	var results []PatternResult
	// Analyze topic relationships and patterns using topicAnalyzer
	if pa.topicAnalyzer != nil {
		patterns := pa.topicAnalyzer.AnalyzeTopics(pa.contextHistory)
		for _, pattern := range patterns {
			results = append(results, PatternResult{
				Pattern:    &pattern,
				Confidence: 0.8, // Default confidence for topical patterns
				Context:    &pa.contextHistory[len(pa.contextHistory)-1],
			})
		}
	}
	return results
}

func (pa *PatternAnalysisEngine) analyzeContextualPatterns() []PatternResult {
	var results []PatternResult
	// Analyze patterns based on user context and environment
	if len(pa.contextHistory) > 0 {
		currentContext := pa.contextHistory[len(pa.contextHistory)-1]
		for _, pattern := range pa.patterns {
			if pattern.Type == ContextualPattern {
				results = append(results, PatternResult{
					Pattern:    &pattern,
					Confidence: 0.7, // Default confidence for contextual patterns
					Context:    &currentContext,
				})
			}
		}
	}
	return results
}

func (pa *PatternAnalysisEngine) analyzeHybridPatterns() []PatternResult {
	var results []PatternResult
	// Analyze patterns that combine multiple pattern types
	if len(pa.contextHistory) > 0 {
		currentContext := pa.contextHistory[len(pa.contextHistory)-1]
		for _, pattern := range pa.patterns {
			if pattern.Type == HybridPattern {
				results = append(results, PatternResult{
					Pattern:    &pattern,
					Confidence: 0.75, // Default confidence for hybrid patterns
					Context:    &currentContext,
				})
			}
		}
	}
	return results
}

func (pa *PatternAnalysisEngine) filterAndRankPatterns(results []PatternResult) []PatternResult {
	// Sort results by confidence
	sort.Slice(results, func(i, j int) bool {
		return results[i].Confidence > results[j].Confidence
	})

	// Filter out low confidence results
	filtered := make([]PatternResult, 0)
	for _, result := range results {
		if result.Confidence >= 0.5 { // Minimum confidence threshold
			filtered = append(filtered, result)
		}
	}
	return filtered
}

func (pa *PatternAnalysisEngine) detectInteractionSequences() [][]Interaction {
	var sequences [][]Interaction
	if len(pa.contextHistory) == 0 {
		return sequences
	}

	currentContext := pa.contextHistory[len(pa.contextHistory)-1]
	if len(currentContext.Interactions) == 0 {
		return sequences
	}

	// Group interactions by type and time proximity
	currentSeq := []Interaction{currentContext.Interactions[0]}
	for i := 1; i < len(currentContext.Interactions); i++ {
		interaction := currentContext.Interactions[i]
		if interaction.Timestamp.Sub(currentSeq[len(currentSeq)-1].Timestamp) <= time.Minute {
			currentSeq = append(currentSeq, interaction)
		} else {
			sequences = append(sequences, currentSeq)
			currentSeq = []Interaction{interaction}
		}
	}
	sequences = append(sequences, currentSeq)
	return sequences
}

func (pa *PatternAnalysisEngine) matchBehavioralPattern(sequence []Interaction) *AnalysisPattern {
	for _, pattern := range pa.patterns {
		if pattern.Type != BehavioralPattern {
			continue
		}

		// Check if sequence matches pattern conditions
		matches := 0
		for _, condition := range pattern.Conditions {
			if pa.checkCondition(sequence, condition) {
				matches++
			}
		}

		if matches >= pattern.MinMatches {
			return &pattern
		}
	}
	return nil
}

func (pa *PatternAnalysisEngine) checkCondition(sequence []Interaction, condition PatternCondition) bool {
	// Basic condition checking implementation
	for _, interaction := range sequence {
		if interaction.Type == condition.Field {
			return true
		}
	}
	return false
}

func (pa *PatternAnalysisEngine) calculatePatternConfidence(pattern *AnalysisPattern, sequence []Interaction) float64 {
	matches := 0
	for _, condition := range pattern.Conditions {
		if pa.checkCondition(sequence, condition) {
			matches++
		}
	}

	// Calculate confidence based on matches and pattern weight
	confidence := float64(matches) / float64(len(pattern.Conditions)) * pattern.Weight
	if confidence > 1.0 {
		confidence = 1.0
	}
	return confidence
}

func (pa *PatternAnalysisEngine) extractPatternContext(sequence []Interaction) *ContextSnapshot {
	if len(pa.contextHistory) == 0 {
		return nil
	}

	// Find the context containing the sequence's timestamp
	for _, context := range pa.contextHistory {
		if len(context.Interactions) > 0 {
			firstInteraction := context.Interactions[0]
			if firstInteraction.Timestamp.Equal(sequence[0].Timestamp) {
				return &context
			}
		}
	}
	return &pa.contextHistory[len(pa.contextHistory)-1]
}

func (pa *PatternAnalysisEngine) detectTimingPatterns() []time.Duration {
	var patterns []time.Duration
	if len(pa.contextHistory) < 2 {
		return patterns
	}

	// Calculate time gaps between consecutive contexts
	for i := 1; i < len(pa.contextHistory); i++ {
		gap := pa.contextHistory[i].Timestamp.Sub(pa.contextHistory[i-1].Timestamp)
		patterns = append(patterns, gap)
	}
	return patterns
}

func (pa *PatternAnalysisEngine) analyzeTiming(gap time.Duration) PatternResult {
	// Find matching timing pattern
	for _, pattern := range pa.patterns {
		if pattern.Type == BehavioralPattern && pattern.MaxGap > 0 {
			if gap <= pattern.MaxGap {
				return PatternResult{
					Pattern:    &pattern,
					Confidence: 0.8, // Default confidence for timing patterns
					Context:    &pa.contextHistory[len(pa.contextHistory)-1],
				}
			}
		}
	}
	return PatternResult{}
}

type MoodTransition struct {
	From      string
	To        string
	Timestamp time.Time
}

func (pa *PatternAnalysisEngine) detectMoodTransitions() []MoodTransition {
	var transitions []MoodTransition
	if len(pa.contextHistory) < 2 {
		return transitions
	}

	// Detect mood changes between consecutive contexts
	for i := 1; i < len(pa.contextHistory); i++ {
		prevMood := pa.contextHistory[i-1].Mood
		currMood := pa.contextHistory[i].Mood
		if prevMood != nil && currMood != nil && prevMood.Primary != currMood.Primary {
			transitions = append(transitions, MoodTransition{
				From:      prevMood.Primary,
				To:        currMood.Primary,
				Timestamp: pa.contextHistory[i].Timestamp,
			})
		}
	}
	return transitions
}

func (pa *PatternAnalysisEngine) analyzeMoodTransition(transition MoodTransition) PatternResult {
	// Find matching emotional pattern based on transition
	for _, pattern := range pa.patterns {
		if pattern.Type == EmotionalPattern {
			// Calculate confidence based on transition
			confidence := 0.7 // Base confidence
			if transition.From == transition.To {
				confidence *= 0.8 // Lower confidence for same mood transitions
			}
			return PatternResult{
				Pattern:    &pattern,
				Confidence: confidence,
				Context:    &pa.contextHistory[len(pa.contextHistory)-1],
			}
		}
	}
	return PatternResult{}
}

type EmotionalTrigger struct {
	Type      string
	Intensity float64
	Timestamp time.Time
}

func (pa *PatternAnalysisEngine) detectEmotionalTriggers() []EmotionalTrigger {
	var triggers []EmotionalTrigger
	if len(pa.contextHistory) == 0 {
		return triggers
	}

	currentContext := pa.contextHistory[len(pa.contextHistory)-1]
	for _, interaction := range currentContext.Interactions {
		if interaction.Type == "emotion" {
			triggers = append(triggers, EmotionalTrigger{
				Type:      interaction.Type,
				Intensity: interaction.Intensity,
				Timestamp: interaction.Timestamp,
			})
		}
	}
	return triggers
}

func (pa *PatternAnalysisEngine) analyzeEmotionalTrigger(trigger EmotionalTrigger) PatternResult {
	// Find matching emotional pattern
	for _, pattern := range pa.patterns {
		if pattern.Type == EmotionalPattern {
			return PatternResult{
				Pattern:    &pattern,
				Confidence: trigger.Intensity, // Use trigger intensity as confidence
				Context:    &pa.contextHistory[len(pa.contextHistory)-1],
			}
		}
	}
	return PatternResult{}
}

func (pa *PatternAnalysisEngine) calculateEmotionalResonance() float64 {
	if len(pa.contextHistory) == 0 {
		return 0.0
	}

	// Calculate average emotional intensity from recent interactions
	currentContext := pa.contextHistory[len(pa.contextHistory)-1]
	var totalIntensity float64
	var count int

	for _, interaction := range currentContext.Interactions {
		if interaction.Type == "emotion" {
			totalIntensity += interaction.Intensity
			count++
		}
	}

	if count == 0 {
		return 0.0
	}
	return totalIntensity / float64(count)
}

func (pa *PatternAnalysisEngine) analyzeEmotionalResonance(resonance float64) PatternResult {
	// Find matching emotional pattern
	for _, pattern := range pa.patterns {
		if pattern.Type == EmotionalPattern {
			return PatternResult{
				Pattern:    &pattern,
				Confidence: resonance, // Use calculated resonance as confidence
				Context:    &pa.contextHistory[len(pa.contextHistory)-1],
			}
		}
	}
	return PatternResult{}
}

// Enhanced persistence features
type TopicPersistenceEnhanced struct {
	*TopicPersistence
	relationshipGraph map[string]map[string]float64
	patternHistory    []PatternOccurrence
	contextCache      *ContextCache
}

type PatternOccurrence struct {
	PatternID string
	Timestamp time.Time
	Context   *ContextSnapshot
	Strength  float64
	Duration  time.Duration
}

type ContextCache struct {
	recent     []ContextSnapshot
	indexed    map[string][]int
	maxSize    int
	timeWindow time.Duration
}

func (tpe *TopicPersistenceEnhanced) SavePatternOccurrence(occurrence PatternOccurrence) error {
	// Store in database
	_, err := tpe.db.Exec(`
        INSERT INTO pattern_occurrences (
            pattern_id, timestamp, context_data, strength, duration
        ) VALUES ($1, $2, $3, $4, $5)
    `, occurrence.PatternID, occurrence.Timestamp, occurrence.Context,
		occurrence.Strength, occurrence.Duration)

	if err != nil {
		return err
	}

	// Update pattern history
	tpe.patternHistory = append(tpe.patternHistory, occurrence)

	// Trim history if needed
	if len(tpe.patternHistory) > 1000 {
		tpe.patternHistory = tpe.patternHistory[1:]
	}

	return nil
}

func (tpe *TopicPersistenceEnhanced) UpdateRelationshipGraph(from, to string, strength float64) {
	if _, exists := tpe.relationshipGraph[from]; !exists {
		tpe.relationshipGraph[from] = make(map[string]float64)
	}

	// Update bidirectional relationship
	tpe.relationshipGraph[from][to] = strength

	if _, exists := tpe.relationshipGraph[to]; !exists {
		tpe.relationshipGraph[to] = make(map[string]float64)
	}
	tpe.relationshipGraph[to][from] = strength
}

func (tpe *TopicPersistenceEnhanced) GetRelatedTopics(topic string, minStrength float64) []string {
	var related []string

	if relationships, exists := tpe.relationshipGraph[topic]; exists {
		for relatedTopic, strength := range relationships {
			if strength >= minStrength {
				related = append(related, relatedTopic)
			}
		}
	}

	return related
}

func (tpe *TopicPersistenceEnhanced) UpdateContextCache(context ContextSnapshot) {
	if tpe.contextCache == nil {
		tpe.contextCache = &ContextCache{
			recent:     make([]ContextSnapshot, 0),
			indexed:    make(map[string][]int),
			maxSize:    1000,
			timeWindow: 24 * time.Hour,
		}
	}

	// Add to recent contexts
	tpe.contextCache.recent = append(tpe.contextCache.recent, context)
	if len(tpe.contextCache.recent) > tpe.contextCache.maxSize {
		tpe.contextCache.recent = tpe.contextCache.recent[1:]
	}

	// Index by timestamp
	timestamp := context.Timestamp.Format(time.RFC3339)
	tpe.contextCache.indexed[timestamp] = append(tpe.contextCache.indexed[timestamp], len(tpe.contextCache.recent)-1)
}

// SQL schema for enhanced persistence
const EnhancedPersistenceSchema = `
CREATE TABLE IF NOT EXISTS pattern_occurrences (
    id SERIAL PRIMARY KEY,
    pattern_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    context_data JSONB NOT NULL,
    strength FLOAT NOT NULL,
    duration INTERVAL NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_pattern_occurrences_pattern_id ON pattern_occurrences(pattern_id);
CREATE INDEX IF NOT EXISTS idx_pattern_occurrences_timestamp ON pattern_occurrences(timestamp);

CREATE TABLE IF NOT EXISTS topic_relationships (
    from_topic VARCHAR(255) NOT NULL,
    to_topic VARCHAR(255) NOT NULL,
    strength FLOAT NOT NULL,
    last_updated TIMESTAMP NOT NULL,
    metadata JSONB,
    PRIMARY KEY (from_topic, to_topic)
);

CREATE INDEX IF NOT EXISTS idx_topic_relationships_strength ON topic_relationships(strength);
`
