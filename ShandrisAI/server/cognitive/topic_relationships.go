package cognitive

import (
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
)

type InteractionContext struct {
	EmotionalTone map[string]float64
	UserContext   map[string]float64
	Domain        string
	Traits        []string
}

func (ic *InteractionContext) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"emotionalTone": ic.EmotionalTone,
		"userContext":   ic.UserContext,
		"domain":        ic.Domain,
		"traits":        ic.Traits,
	}
}

type ContextImpact struct {
	EmotionalImpact float64
	UserImpact      float64
	DomainImpact    float64
}

type RelationshipContextAnalyzer struct {
	emotionalWeights map[string]float64
	userWeights      map[string]float64
	domainWeights    map[string]float64
}

func newRelationshipContextAnalyzer() *RelationshipContextAnalyzer {
	return &RelationshipContextAnalyzer{
		emotionalWeights: make(map[string]float64),
		userWeights:      make(map[string]float64),
		domainWeights:    make(map[string]float64),
	}
}

func (rca *RelationshipContextAnalyzer) AnalyzeContext(context *InteractionContext) *ContextImpact {
	return &ContextImpact{
		EmotionalImpact: rca.analyzeEmotionalImpact(context),
		UserImpact:      rca.analyzeUserImpact(context),
		DomainImpact:    rca.analyzeDomainImpact(context),
	}
}

func (rca *RelationshipContextAnalyzer) analyzeEmotionalImpact(context *InteractionContext) float64 {
	if context == nil || len(context.EmotionalTone) == 0 {
		return 0.5
	}

	var totalImpact float64
	for _, value := range context.EmotionalTone {
		totalImpact += value
	}
	return totalImpact / float64(len(context.EmotionalTone))
}

func (rca *RelationshipContextAnalyzer) analyzeUserImpact(context *InteractionContext) float64 {
	if context == nil || len(context.UserContext) == 0 {
		return 0.5
	}

	var totalImpact float64
	for _, value := range context.UserContext {
		totalImpact += value
	}
	return totalImpact / float64(len(context.UserContext))
}

func (rca *RelationshipContextAnalyzer) analyzeDomainImpact(context *InteractionContext) float64 {
	if context == nil || context.Domain == "" {
		return 0.5
	}

	// Calculate impact based on domain weights
	if weight, exists := rca.domainWeights[context.Domain]; exists {
		return weight
	}
	return 0.5
}

func (rca *RelationshipContextAnalyzer) MatchesRequirements(context *RelationshipContext, requirements map[string]float64) bool {
	if requirements == nil {
		return true
	}

	for key, requiredValue := range requirements {
		if value, exists := context.UserContext[key]; !exists || value < requiredValue {
			return false
		}
	}
	return true
}

func (rca *RelationshipContextAnalyzer) CalculateRelevance(context *RelationshipContext, requirements map[string]float64) float64 {
	if requirements == nil {
		return 1.0
	}

	var totalScore float64
	var count int

	for key, requiredValue := range requirements {
		if value, exists := context.UserContext[key]; exists {
			totalScore += math.Min(value/requiredValue, 1.0)
			count++
		}
	}

	if count == 0 {
		return 0.0
	}
	return totalScore / float64(count)
}

type RelationshipStrengthCalculator struct {
	baseWeight    float64
	contextWeight float64
	recencyWeight float64
}

func newRelationshipStrengthCalculator() *RelationshipStrengthCalculator {
	return &RelationshipStrengthCalculator{
		baseWeight:    0.4,
		contextWeight: 0.3,
		recencyWeight: 0.3,
	}
}

func (rsc *RelationshipStrengthCalculator) CalculateStrength(rel *TopicRelationship, impact *ContextImpact) float64 {
	// Placeholder implementation
	return 0.5
}

type RelationshipEvolution struct {
	growthRate float64
	decayRate  float64
	threshold  float64
}

func newRelationshipEvolution() *RelationshipEvolution {
	return &RelationshipEvolution{
		growthRate: 0.1,
		decayRate:  0.05,
		threshold:  0.3,
	}
}

func (re *RelationshipEvolution) EvolveRelationship(rel *TopicRelationship) {
	// Placeholder implementation
	rel.Strength = rel.Strength * (1 + re.growthRate)
}

type RelationshipFilters struct {
	Type                RelationType
	MinStrength         float64
	ContextRequirements map[string]float64
}

type ClusterAnalysis struct {
	Topics        []string
	Relationships map[string]float64
	Centrality    map[string]float64
	Clusters      map[string][]string
}

func calculateSignificance(context *RelationshipContext) float64 {
	// Calculate significance based on frequency and emotional impact
	emotionalWeight := 0.0
	for _, value := range context.EmotionalTone {
		emotionalWeight += value
	}
	if len(context.EmotionalTone) > 0 {
		emotionalWeight /= float64(len(context.EmotionalTone))
	}

	return (float64(context.Frequency)*0.6 + emotionalWeight*0.4) / 100.0
}

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

func initializeRelationshipPatterns() map[string]*RelationshipPattern {
	patterns := make(map[string]*RelationshipPattern)
	// Initialize with default patterns
	patterns["default"] = &RelationshipPattern{
		Type:       Direct,
		Triggers:   []string{},
		Conditions: make(map[string]float64),
		Evolution: &EvolutionRule{
			GrowthRate:   0.1,
			DecayRate:    0.05,
			Threshold:    0.3,
			Dependencies: []string{},
		},
		Strength: 0.5,
	}
	return patterns
}

func isRelationshipBidirectional(relType RelationType) bool {
	switch relType {
	case Associative, Semantic, Emotional:
		return true
	default:
		return false
	}
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

func updateSharedTraits(existing []string, newTraits []string) []string {
	// Create a map to track unique traits
	traitMap := make(map[string]bool)

	// Add existing traits
	for _, trait := range existing {
		traitMap[trait] = true
	}

	// Add new traits
	for _, trait := range newTraits {
		traitMap[trait] = true
	}

	// Convert map back to slice
	result := make([]string, 0, len(traitMap))
	for trait := range traitMap {
		result = append(result, trait)
	}

	return result
}

func (trm *TopicRelationshipManager) updateRelationshipContext(
	rel *TopicRelationship,
	context *InteractionContext,
	impact *ContextImpact,
) {
	// Update emotional tone with impact influence
	for emotion, value := range context.EmotionalTone {
		current := rel.Context.EmotionalTone[emotion]
		rel.Context.EmotionalTone[emotion] = (current*0.7 + value*0.3) * (1 + impact.EmotionalImpact)
	}

	// Update user context with impact influence
	for ctx, value := range context.UserContext {
		rel.Context.UserContext[ctx] = value * (1 + impact.UserImpact)
	}

	// Update frequency and significance
	rel.Context.Frequency++
	rel.Context.Significance = calculateSignificance(rel.Context) * (1 + impact.DomainImpact)

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

func (trm *TopicRelationshipManager) calculateClusterStrength(topic1, topic2 string) float64 {
	// Find direct relationship between topics
	for _, rel := range trm.relationships {
		if (rel.FromTopic == topic1 && rel.ToTopic == topic2) ||
			(rel.Bidirectional && rel.FromTopic == topic2 && rel.ToTopic == topic1) {
			return rel.Strength
		}
	}
	return 0.0
}

func (trm *TopicRelationshipManager) calculateTopicCentrality(topic string, relationships map[string]float64) float64 {
	var totalStrength float64
	var count int
	for key, strength := range relationships {
		if key[:len(topic)] == topic || key[len(key)-len(topic):] == topic {
			totalStrength += strength
			count++
		}
	}
	if count == 0 {
		return 0.0
	}
	return totalStrength / float64(count)
}

func (trm *TopicRelationshipManager) identifySubClusters(analysis *ClusterAnalysis) map[string][]string {
	clusters := make(map[string][]string)
	threshold := 0.3 // Minimum relationship strength for clustering

	// Simple clustering based on relationship strength
	for i, topic1 := range analysis.Topics {
		cluster := []string{topic1}
		for j, topic2 := range analysis.Topics {
			if i != j {
				key := topic1 + ":" + topic2
				if strength, exists := analysis.Relationships[key]; exists && strength >= threshold {
					cluster = append(cluster, topic2)
				}
			}
		}
		if len(cluster) > 1 {
			clusters[topic1] = cluster
		}
	}
	return clusters
}
