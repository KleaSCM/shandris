package cognitive

import (
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
)

// TopicRelationshipManager handles sophisticated topic relationships
type TopicRelationshipManager struct {
	relationships   map[string]*TopicRelationship
	patterns        map[string]*RelationshipPattern
	contextAnalyzer *RelationshipContextAnalyzer
	strengthCalc    *RelationshipStrengthCalculator
	evolution       *RelationshipEvolution
}

type TopicRelationship struct {
	ID            string
	FromTopic     string
	ToTopic       string
	Type          RelationType
	Strength      float64
	Bidirectional bool
	Context       *RelationshipContext
	History       []RelationshipEvent
	Properties    map[string]interface{}
	LastActive    time.Time
}

type RelationType string

const (
	Direct       RelationType = "direct"
	Hierarchical RelationType = "hierarchical"
	Associative  RelationType = "associative"
	Causal       RelationType = "causal"
	Temporal     RelationType = "temporal"
	Semantic     RelationType = "semantic"
	Emotional    RelationType = "emotional"
)

type RelationshipContext struct {
	Domain        string
	UserContext   map[string]float64
	EmotionalTone map[string]float64
	Frequency     int
	Significance  float64
	SharedTraits  []string
}

type RelationshipPattern struct {
	Type       RelationType
	Triggers   []string
	Conditions map[string]float64
	Evolution  *EvolutionRule
	Strength   float64
}

type EvolutionRule struct {
	GrowthRate   float64
	DecayRate    float64
	Threshold    float64
	Dependencies []string
}

type RelationshipEvent struct {
	Timestamp time.Time
	Type      string
	Impact    float64
	Context   map[string]interface{}
}

func NewTopicRelationshipManager() *TopicRelationshipManager {
	return &TopicRelationshipManager{
		relationships:   make(map[string]*TopicRelationship),
		patterns:        initializeRelationshipPatterns(),
		contextAnalyzer: newRelationshipContextAnalyzer(),
		strengthCalc:    newRelationshipStrengthCalculator(),
		evolution:       newRelationshipEvolution(),
	}
}

// CreateRelationship establishes a new topic relationship
func (trm *TopicRelationshipManager) CreateRelationship(from, to string, relType RelationType) *TopicRelationship {
	rel := &TopicRelationship{
		ID:            uuid.New().String(),
		FromTopic:     from,
		ToTopic:       to,
		Type:          relType,
		Strength:      0.1, // Initial strength
		Bidirectional: isRelationshipBidirectional(relType),
		Context: &RelationshipContext{
			EmotionalTone: make(map[string]float64),
			UserContext:   make(map[string]float64),
		},
		Properties: make(map[string]interface{}),
		History:    make([]RelationshipEvent, 0),
		LastActive: time.Now(),
	}

	trm.relationships[rel.ID] = rel
	return rel
}

// UpdateRelationship processes new interactions and updates relationship properties
func (trm *TopicRelationshipManager) UpdateRelationship(rel *TopicRelationship, context *InteractionContext) {
	// Analyze context impact
	contextImpact := trm.contextAnalyzer.AnalyzeContext(context)

	// Update relationship strength
	newStrength := trm.strengthCalc.CalculateStrength(rel, contextImpact)

	// Record event
	event := RelationshipEvent{
		Timestamp: time.Now(),
		Type:      "interaction",
		Impact:    newStrength - rel.Strength,
		Context:   context.ToMap(),
	}

	// Update relationship
	rel.Strength = newStrength
	rel.LastActive = time.Now()
	rel.History = append(rel.History, event)

	// Update context
	trm.updateRelationshipContext(rel, context, contextImpact)

	// Apply evolution rules
	trm.evolution.EvolveRelationship(rel)
}

// FindRelatedTopics finds topics related to the given topic with advanced filtering
func (trm *TopicRelationshipManager) FindRelatedTopics(topicID string, filters *RelationshipFilters) []TopicRelationship {
	var related []TopicRelationship

	for _, rel := range trm.relationships {
		if (rel.FromTopic == topicID || (rel.Bidirectional && rel.ToTopic == topicID)) &&
			trm.matchesFilters(rel, filters) {
			related = append(related, *rel)
		}
	}

	// Sort by relevance
	trm.sortByRelevance(related, filters)

	return related
}

// AnalyzeTopicCluster analyzes relationships within a topic cluster
func (trm *TopicRelationshipManager) AnalyzeTopicCluster(topics []string) *ClusterAnalysis {
	analysis := &ClusterAnalysis{
		Topics:        topics,
		Relationships: make(map[string]float64),
		Centrality:    make(map[string]float64),
		Clusters:      make(map[string][]string),
	}

	// Calculate relationship strengths within cluster
	for i, topic1 := range topics {
		for j, topic2 := range topics {
			if i != j {
				strength := trm.calculateClusterStrength(topic1, topic2)
				analysis.Relationships[topic1+":"+topic2] = strength
			}
		}
	}

	// Calculate topic centrality
	for _, topic := range topics {
		centrality := trm.calculateTopicCentrality(topic, analysis.Relationships)
		analysis.Centrality[topic] = centrality
	}

	// Identify sub-clusters
	analysis.Clusters = trm.identifySubClusters(analysis)

	return analysis
}

// Helper functions

func (trm *TopicRelationshipManager) updateRelationshipContext(
	rel *TopicRelationship,
	context *InteractionContext,
	impact *ContextImpact,
) {
	// Update emotional tone
	for emotion, value := range context.EmotionalTone {
		current := rel.Context.EmotionalTone[emotion]
		rel.Context.EmotionalTone[emotion] = (current*0.7 + value*0.3) // Weighted average
	}

	// Update user context
	for ctx, value := range context.UserContext {
		rel.Context.UserContext[ctx] = value
	}

	// Update frequency and significance
	rel.Context.Frequency++
	rel.Context.Significance = calculateSignificance(rel.Context)

	// Update shared traits
	rel.Context.SharedTraits = updateSharedTraits(rel.Context.SharedTraits, context.Traits)
}

func (trm *TopicRelationshipManager) matchesFilters(rel *TopicRelationship, filters *RelationshipFilters) bool {
	if filters == nil {
		return true
	}

	// Check type filter
	if filters.Type != "" && rel.Type != filters.Type {
		return false
	}

	// Check strength threshold
	if rel.Strength < filters.MinStrength {
		return false
	}

	// Check context requirements
	if !trm.contextAnalyzer.MatchesRequirements(rel.Context, filters.ContextRequirements) {
		return false
	}

	return true
}

func (trm *TopicRelationshipManager) sortByRelevance(relationships []TopicRelationship, filters *RelationshipFilters) {
	sort.Slice(relationships, func(i, j int) bool {
		scoreI := trm.calculateRelevanceScore(&relationships[i], filters)
		scoreJ := trm.calculateRelevanceScore(&relationships[j], filters)
		return scoreI > scoreJ
	})
}

func (trm *TopicRelationshipManager) calculateRelevanceScore(rel *TopicRelationship, filters *RelationshipFilters) float64 {
	score := rel.Strength

	// Apply recency factor
	timeSinceActive := time.Since(rel.LastActive).Hours()
	recencyFactor := math.Exp(-timeSinceActive / (24 * 30)) // 30-day half-life
	score *= (0.7 + 0.3*recencyFactor)

	// Apply context relevance
	if filters != nil && filters.ContextRequirements != nil {
		contextRelevance := trm.contextAnalyzer.CalculateRelevance(rel.Context, filters.ContextRequirements)
		score *= (0.8 + 0.2*contextRelevance)
	}

	return score
}
