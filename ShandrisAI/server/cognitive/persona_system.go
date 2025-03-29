package cognitive

import (
	"fmt"
	"time"
)

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
	StyleRules  []StyleRule
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

type StyleRule struct {
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
		StyleRules: []StyleRule{
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
		StyleRules: []StyleRule{
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

// GetResponseStyle determines the appropriate response style for the current context
func (ps *PersonaSystem) GetResponseStyle(context *PersonaContext) StyleRule {
	if ps.activePersona == nil {
		return StyleRule{} // Default style
	}

	var bestRule StyleRule
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
