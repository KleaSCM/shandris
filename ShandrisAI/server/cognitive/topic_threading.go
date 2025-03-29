package cognitive

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TopicNode represents a single topic in the conversation
type TopicNode struct {
	ID         string
	Name       string
	Keywords   []string
	Context    string
	StartTime  time.Time
	LastActive time.Time
	Confidence float64
	ParentID   string
	Children   []string
	Related    map[string]float64 // Related topics and their relevance scores
	Attributes map[string]string  // Domain-specific attributes
}

// TopicThread represents an active conversation thread
type TopicThread struct {
	ID          string
	ActiveNodes []string
	MainTopic   string
	StartTime   time.Time
	LastActive  time.Time
	Depth       int
	Context     *EmotionalContext
}

// TopicManager handles topic threading and context switching
type TopicManager struct {
	ActiveThreads   map[string]*TopicThread
	TopicGraph      map[string]*TopicNode
	DomainRules     map[string]DomainRule
	TransitionRules []TransitionRule

	// Configuration
	MaxThreadDepth int
	MinConfidence  float64
	DecayRate      float64
}

// DomainRule defines how to handle specific conversation domains
type DomainRule struct {
	Domain      string
	Keywords    []string
	Validators  []func(string) bool
	Priority    int
	Transitions map[string]float64 // Allowed transitions to other domains
}

// TransitionRule defines when and how topics can transition
type TransitionRule struct {
	FromDomain string
	ToDomain   string
	Conditions []string
	Weight     float64
	MinContext float64
}

// NewTopicManager creates a new topic threading system
func NewTopicManager() *TopicManager {
	tm := &TopicManager{
		ActiveThreads:  make(map[string]*TopicThread),
		TopicGraph:     make(map[string]*TopicNode),
		DomainRules:    initializeDomainRules(),
		MaxThreadDepth: 5,
		MinConfidence:  0.6,
		DecayRate:      0.1,
	}

	tm.initializeTransitionRules()
	return tm
}

// ProcessInput analyzes new input and updates topic threads
func (tm *TopicManager) ProcessInput(input string, context *EmotionalContext) error {
	// Detect topics in input
	topics := tm.detectTopics(input)

	// Update or create threads
	for _, topic := range topics {
		if thread := tm.findRelevantThread(topic.Domain); thread != nil {
			tm.updateThread(thread, topic.Domain, context)
		} else {
			tm.createNewThread(topic.Domain, context)
		}
	}

	// Prune old or inactive threads
	tm.pruneThreads()

	return nil
}

// detectTopics identifies potential topics in the input
func (tm *TopicManager) detectTopics(input string) []TopicDetection {
	var detections []TopicDetection

	// Normalize input
	normalizedInput := strings.ToLower(input)

	// Check each domain's keywords and rules
	for domain, rule := range tm.DomainRules {
		confidence := 0.0
		matches := make(map[string]int)

		// Check for keyword matches
		for _, keyword := range rule.Keywords {
			if count := strings.Count(normalizedInput, keyword); count > 0 {
				matches[keyword] = count
				confidence += float64(count) * 0.2 // Base confidence boost per match
			}
		}

		// Apply domain-specific validators
		for _, validator := range rule.Validators {
			if validator(normalizedInput) {
				confidence += 0.3 // Additional confidence for validated patterns
			}
		}

		// If we have sufficient confidence, create a detection
		if confidence >= tm.MinConfidence {
			detection := TopicDetection{
				Domain:     domain,
				Keywords:   matches,
				Confidence: confidence,
				Priority:   rule.Priority,
			}
			detections = append(detections, detection)
		}
	}

	return detections
}

// findRelevantThread finds the most relevant existing thread for a topic
func (tm *TopicManager) findRelevantThread(topic string) *TopicThread {
	var bestMatch *TopicThread
	var highestRelevance float64

	for _, thread := range tm.ActiveThreads {
		relevance := tm.calculateThreadRelevance(thread, topic)
		if relevance > highestRelevance && relevance >= tm.MinConfidence {
			highestRelevance = relevance
			bestMatch = thread
		}
	}

	return bestMatch
}

// updateThread updates an existing thread with new information
func (tm *TopicManager) updateThread(thread *TopicThread, topic string, context *EmotionalContext) {
	now := time.Now()

	// Update timing
	thread.LastActive = now

	// Add new topic node if it doesn't exist
	if _, exists := tm.TopicGraph[topic]; !exists {
		newNode := &TopicNode{
			ID:         uuid.New().String(),
			Name:       topic,
			StartTime:  now,
			LastActive: now,
			Related:    make(map[string]float64),
			Attributes: make(map[string]string),
		}
		tm.TopicGraph[topic] = newNode
	}

	// Update relationships
	currentNode := tm.TopicGraph[topic]
	for _, existingTopic := range thread.ActiveNodes {
		if existingNode, exists := tm.TopicGraph[existingTopic]; exists {
			// Update bidirectional relationship scores
			currentNode.Related[existingTopic] = tm.calculateRelationshipScore(topic, existingTopic)
			existingNode.Related[topic] = currentNode.Related[existingTopic]
		}
	}

	// Update emotional context
	if context != nil {
		thread.Context = context
		currentNode.Attributes["mood"] = context.PrimaryMood
		currentNode.Attributes["intensity"] = fmt.Sprintf("%.2f", context.Intensity)
	}

	// Add to active nodes if not present
	if !contains(thread.ActiveNodes, topic) {
		thread.ActiveNodes = append(thread.ActiveNodes, topic)
		// Maintain max depth
		if len(thread.ActiveNodes) > tm.MaxThreadDepth {
			thread.ActiveNodes = thread.ActiveNodes[1:]
		}
	}
}

// createNewThread starts a new conversation thread
func (tm *TopicManager) createNewThread(topic string, context *EmotionalContext) *TopicThread {
	thread := &TopicThread{
		ID:          uuid.New().String(),
		ActiveNodes: []string{topic},
		MainTopic:   topic,
		StartTime:   time.Now(),
		LastActive:  time.Now(),
		Depth:       1,
		Context:     context,
	}

	tm.ActiveThreads[thread.ID] = thread
	return thread
}

// pruneThreads removes old or inactive threads
func (tm *TopicManager) pruneThreads() {
	now := time.Now()
	for id, thread := range tm.ActiveThreads {
		// Remove threads inactive for more than 30 minutes
		if now.Sub(thread.LastActive) > 30*time.Minute {
			// Store important information before removing
			tm.archiveThread(thread)
			delete(tm.ActiveThreads, id)
		}
	}
}

// Helper functions

func (tm *TopicManager) calculateThreadRelevance(thread *TopicThread, topic string) float64 {
	relevance := 0.0

	// Check direct topic matches
	for _, activeTopic := range thread.ActiveNodes {
		if activeTopic == topic {
			relevance += 1.0
		}

		// Check related topics
		if node, exists := tm.TopicGraph[activeTopic]; exists {
			if relScore, hasRel := node.Related[topic]; hasRel {
				relevance += relScore
			}
		}
	}

	// Apply time decay
	timeSinceActive := time.Since(thread.LastActive).Minutes()
	decay := math.Exp(-tm.DecayRate * timeSinceActive)

	return relevance * decay
}

func (tm *TopicManager) calculateRelationshipScore(topic1, topic2 string) float64 {
	// Get domain rules for both topics
	domain1 := tm.getDomain(topic1)
	domain2 := tm.getDomain(topic2)

	if domain1 == "" || domain2 == "" {
		return 0.0
	}

	// Check transition rules
	if rule, exists := tm.DomainRules[domain1]; exists {
		if score, allowed := rule.Transitions[domain2]; allowed {
			return score
		}
	}

	return 0.0
}

func (tm *TopicManager) archiveThread(thread *TopicThread) {
	// TODO: Implement persistence logic here
	// This will be implemented when we add persistence
}

// TopicDetection represents a detected topic with confidence
type TopicDetection struct {
	Domain     string
	Keywords   map[string]int
	Confidence float64
	Priority   int
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// initializeDomainRules sets up the initial domain rules
func initializeDomainRules() map[string]DomainRule {
	return map[string]DomainRule{
		"tech": {
			Domain:   "tech",
			Keywords: []string{"coding", "programming", "software", "computer", "algorithm"},
			Priority: 3,
			Transitions: map[string]float64{
				"science": 0.8,
				"gaming":  0.7,
				"sapphic": 0.5,
			},
		},
		"sapphic": {
			Domain:   "sapphic",
			Keywords: []string{"girls", "lesbian", "queer", "dating", "cute"},
			Priority: 4,
			Transitions: map[string]float64{
				"gaming": 0.6,
				"tech":   0.5,
				"memes":  0.7,
			},
		},
		"gaming": {
			Domain:   "gaming",
			Keywords: []string{"game", "play", "steam", "console", "rpg"},
			Priority: 2,
			Transitions: map[string]float64{
				"tech":    0.6,
				"sapphic": 0.6,
				"memes":   0.8,
			},
		},
		// Add more domains as needed
	}
}

func (tm *TopicManager) initializeTransitionRules() {
	tm.TransitionRules = []TransitionRule{
		{
			FromDomain: "tech",
			ToDomain:   "science",
			Weight:     0.8,
			MinContext: 0.6,
		},
		{
			FromDomain: "sapphic",
			ToDomain:   "gaming",
			Weight:     0.6,
			MinContext: 0.5,
		},
		{
			FromDomain: "gaming",
			ToDomain:   "tech",
			Weight:     0.6,
			MinContext: 0.5,
		},
	}
}
