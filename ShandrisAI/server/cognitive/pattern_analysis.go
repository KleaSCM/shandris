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

// Add more sophisticated mood patterns
func initializeAdvancedMoodPatterns() map[string]MoodPattern {
	return map[string]MoodPattern{
		"sapphic_flirty": {
			Keywords: []string{
				"cute", "pretty", "beautiful", "gorgeous", "attractive",
				"flirt", "tease", "wink", "blush", "smile",
				"gay", "lesbian", "sapphic", "wlw", "queer",
			},
			Sentiment:    0.8,
			MoodShift:    "flirty",
			Intensity:    0.7,
			Decay:        0.1,
			Requirements: []string{"feminine_presence", "romantic_context"},
			Exclusions:   []string{"professional_context", "serious_discussion"},
		},
		"tech_passionate": {
			Keywords: []string{
				"code", "programming", "algorithm", "development", "software",
				"excited", "fascinating", "amazing", "innovative", "elegant",
			},
			Sentiment:    0.6,
			MoodShift:    "enthusiastic",
			Intensity:    0.6,
			Decay:        0.15,
			Requirements: []string{"technical_context"},
		},
		"protective_caring": {
			Keywords: []string{
				"protect", "care", "support", "help", "comfort",
				"safe", "secure", "gentle", "kind", "warm",
			},
			Sentiment:    0.7,
			MoodShift:    "nurturing",
			Intensity:    0.5,
			Decay:        0.08,
			Requirements: []string{"emotional_context", "support_needed"},
		},
		"playful_teasing": {
			Keywords: []string{
				"tease", "joke", "play", "laugh", "giggle",
				"silly", "fun", "witty", "clever", "sassy",
			},
			Sentiment:    0.75,
			MoodShift:    "playful",
			Intensity:    0.6,
			Decay:        0.12,
			Requirements: []string{"comfortable_context", "positive_rapport"},
		},
		// Add more sophisticated patterns...
	}
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
