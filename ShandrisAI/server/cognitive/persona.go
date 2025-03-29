package cognitive

import "time"

// PersonaManager handles different personality modes
type PersonaManager struct {
	CurrentPersona Persona
	Available      map[string]Persona
	Transitions    []PersonaTransition
}

type Persona struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Traits     map[string]float64 `json:"traits"`    // base personality traits
	MoodBias   map[string]float64 `json:"mood_bias"` // mood tendencies
	StyleRules []StyleRule        `json:"style_rules"`
	Active     bool               `json:"active"`
}

type StyleRule struct {
	Trigger    string   `json:"trigger"`
	Response   string   `json:"response"`
	Conditions []string `json:"conditions"`
	Weight     float64  `json:"weight"`
}

type PersonaTransition struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Trigger   string    `json:"trigger"`
	Timestamp time.Time `json:"timestamp"`
}
