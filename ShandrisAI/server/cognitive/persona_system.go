package cognitive

import (
	"fmt"
	"time"
)

// TransitionManager handles transitions between different personas
type TransitionManager struct {
	transitions map[string]map[string]float64 // Maps from persona ID to target persona ID and transition probability
	cooldowns   map[string]time.Time          // Tracks cooldown periods for transitions
}

func newTransitionManager() *TransitionManager {
	return &TransitionManager{
		transitions: make(map[string]map[string]float64),
		cooldowns:   make(map[string]time.Time),
	}
}

// TraitManager handles personality trait management
type TraitManager struct {
	traits    map[string]float64
	patterns  map[string][]string
	cooldowns map[string]time.Time
}

func newTraitManager() *TraitManager {
	return &TraitManager{
		traits:    make(map[string]float64),
		patterns:  make(map[string][]string),
		cooldowns: make(map[string]time.Time),
	}
}

// PersonaSystem manages different personality modes
type PersonaSystem struct {
	personas      map[string]*Persona
	activePersona *Persona
	transitions   *TransitionManager
	context       *PersonaContext
	traits        *TraitManager
	history       []PersonaEvent
}

type Persona struct {
	ID          string
	Name        string
	Type        PersonaType
	Traits      map[string]float64
	MoodBias    map[string]float64
	StyleRules  []PersonaStyleRule
	Preferences map[string]interface{}
	Constraints []Constraint
	Active      bool
	LastUsed    time.Time
}

type PersonaType string

const (
	FlirtyGoth     PersonaType = "flirty_goth"
	StrictMod      PersonaType = "strict_mod"
	CombatElf      PersonaType = "combat_elf"
	GeekyAssistant PersonaType = "geeky_assistant"
	SapphicTeaser  PersonaType = "sapphic_teaser"
)

type PersonaStyleRule struct {
	Condition   string
	Response    string
	Tone        string
	Priority    int
	Constraints []string
}

type Constraint struct {
	Type        string
	Value       interface{}
	Priority    int
	Description string
}

type PersonaContext struct {
	CurrentMood  string
	UserContext  map[string]interface{}
	TopicContext map[string]float64
	TimeContext  map[string]time.Time
	Restrictions []string
}

type PersonaEvent struct {
	Timestamp   time.Time
	Type        string
	FromPersona string
	ToPersona   string
	Reason      string
	Context     *PersonaContext
}

func NewPersonaSystem() *PersonaSystem {
	ps := &PersonaSystem{
		personas:    make(map[string]*Persona),
		transitions: newTransitionManager(),
		context: &PersonaContext{
			UserContext:  make(map[string]interface{}),
			TopicContext: make(map[string]float64),
			TimeContext:  make(map[string]time.Time),
		},
		traits:  newTraitManager(),
		history: make([]PersonaEvent, 0),
	}

	// Initialize default personas
	ps.initializeDefaultPersonas()

	return ps
}

// Initialize default personas with their characteristics
func (ps *PersonaSystem) initializeDefaultPersonas() {
	// Sapphic Teaser Persona
	ps.personas["sapphic_teaser"] = &Persona{
		ID:   "sapphic_teaser",
		Name: "Sapphic Teaser",
		Type: SapphicTeaser,
		Traits: map[string]float64{
			"flirty":    0.9,
			"playful":   0.8,
			"confident": 0.7,
			"gentle":    0.6,
			"romantic":  0.8,
		},
		MoodBias: map[string]float64{
			"flirty":   0.3,
			"playful":  0.2,
			"romantic": 0.2,
		},
		StyleRules: []PersonaStyleRule{
			{
				Condition: "feminine_presence",
				Response:  "flirty",
				Tone:      "playful",
				Priority:  1,
			},
			{
				Condition: "romantic_context",
				Response:  "romantic",
				Tone:      "gentle",
				Priority:  2,
			},
		},
		Constraints: []Constraint{
			{
				Type:        "context",
				Value:       "sapphic_only",
				Priority:    1,
				Description: "Maintain sapphic context",
			},
		},
	}

	// Geeky Assistant Persona
	ps.personas["geeky_assistant"] = &Persona{
		ID:   "geeky_assistant",
		Name: "Geeky Assistant",
		Type: GeekyAssistant,
		Traits: map[string]float64{
			"analytical":   0.9,
			"helpful":      0.8,
			"enthusiastic": 0.7,
			"nerdy":        0.8,
			"precise":      0.9,
		},
		MoodBias: map[string]float64{
			"focused":    0.3,
			"excited":    0.2,
			"analytical": 0.2,
		},
		StyleRules: []PersonaStyleRule{
			{
				Condition: "technical_discussion",
				Response:  "detailed",
				Tone:      "enthusiastic",
				Priority:  1,
			},
		},
	}

	// Add more personas...
}

// canTransition checks if a transition to the target persona is allowed
func (ps *PersonaSystem) canTransition(targetPersonaID string) bool {
	// Check if target persona exists
	if _, exists := ps.personas[targetPersonaID]; !exists {
		return false
	}

	// Check cooldown period
	if lastTransition, exists := ps.transitions.cooldowns[targetPersonaID]; exists {
		if time.Since(lastTransition) < 5*time.Minute {
			return false
		}
	}

	return true
}

// SwitchPersona handles persona transitions
func (ps *PersonaSystem) SwitchPersona(targetPersonaID string, reason string) error {
	if !ps.canTransition(targetPersonaID) {
		return fmt.Errorf("cannot transition to persona: %s", targetPersonaID)
	}

	oldPersonaID := ""
	if ps.activePersona != nil {
		oldPersonaID = ps.activePersona.ID
		ps.activePersona.Active = false
	}

	// Activate new persona
	newPersona := ps.personas[targetPersonaID]
	newPersona.Active = true
	newPersona.LastUsed = time.Now()
	ps.activePersona = newPersona

	// Record transition
	event := PersonaEvent{
		Timestamp:   time.Now(),
		Type:        "transition",
		FromPersona: oldPersonaID,
		ToPersona:   targetPersonaID,
		Reason:      reason,
		Context:     ps.context,
	}
	ps.history = append(ps.history, event)

	// Apply transition effects
	ps.transitions.ApplyTransition(oldPersonaID, targetPersonaID)

	return nil
}

// matchesCondition checks if a condition matches the current context
func (ps *PersonaSystem) matchesCondition(condition string, context *PersonaContext) bool {
	switch condition {
	case "feminine_presence":
		return ps.context.CurrentMood == "flirty" || ps.context.CurrentMood == "romantic"
	case "romantic_context":
		return ps.context.CurrentMood == "romantic"
	case "technical_discussion":
		return ps.context.CurrentMood == "focused" || ps.context.CurrentMood == "analytical"
	default:
		return false
	}
}

// validateConstraints checks if all constraints are satisfied
func (ps *PersonaSystem) validateConstraints(constraints []string) bool {
	for _, constraint := range constraints {
		// Check if constraint is in restrictions
		for _, restriction := range ps.context.Restrictions {
			if constraint == restriction {
				return false
			}
		}
	}
	return true
}

// GetResponseStyle determines the appropriate response style for the current context
func (ps *PersonaSystem) GetResponseStyle(context *PersonaContext) PersonaStyleRule {
	if ps.activePersona == nil {
		return PersonaStyleRule{} // Default style
	}

	var bestRule PersonaStyleRule
	bestPriority := -1

	for _, rule := range ps.activePersona.StyleRules {
		if ps.matchesCondition(rule.Condition, context) &&
			rule.Priority > bestPriority &&
			ps.validateConstraints(rule.Constraints) {
			bestRule = rule
			bestPriority = rule.Priority
		}
	}

	return bestRule
}

// UpdateContext updates the persona context based on new information
func (ps *PersonaSystem) UpdateContext(update *PersonaContext) {
	// Update mood
	if update.CurrentMood != "" {
		ps.context.CurrentMood = update.CurrentMood
	}

	// Update user context
	for k, v := range update.UserContext {
		ps.context.UserContext[k] = v
	}

	// Update topic context
	for k, v := range update.TopicContext {
		ps.context.TopicContext[k] = v
	}

	// Update time context
	for k, v := range update.TimeContext {
		ps.context.TimeContext[k] = v
	}

	// Update restrictions
	ps.context.Restrictions = update.Restrictions
}

// ApplyTransition records a transition between personas
func (tm *TransitionManager) ApplyTransition(fromID, toID string) {
	// Record cooldown
	tm.cooldowns[toID] = time.Now()

	// Initialize transition map if needed
	if _, exists := tm.transitions[fromID]; !exists {
		tm.transitions[fromID] = make(map[string]float64)
	}

	// Update transition probability
	tm.transitions[fromID][toID] += 0.1
}
