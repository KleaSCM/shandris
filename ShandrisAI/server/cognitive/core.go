package cognitive

import (
	"time"
)

// CoreSystem represents Shandris's main cognitive architecture
type CoreSystem struct {
	MoodEngine     *MoodEngine
	TraitSystem    *TraitSystem
	TopicManager   *TopicManager
	MemoryTimeline *Timeline
	PersonaManager *PersonaManager
}

// MoodEngine handles emotional state and responses
type MoodEngine struct {
	CurrentState MoodState
	History      []MoodState
	Modifiers    map[string]float64
}

// ProcessInteraction processes an interaction and updates the emotional context
func (m *MoodEngine) ProcessInteraction(interaction *Interaction, context *EmotionalContext) *MoodUpdate {
	return &MoodUpdate{
		NewState: &m.CurrentState,
	}
}

type MoodState struct {
	Primary   string         `json:"primary"`
	Secondary string         `json:"secondary"`
	Intensity float64        `json:"intensity"`
	Timestamp time.Time      `json:"timestamp"`
	Context   map[string]any `json:"context"`
}

// TraitSystem handles personality inference and adaptation
type TraitSystem struct {
	Traits    map[string]Trait
	Patterns  []TraitPattern
	UpdatedAt time.Time
}

type Trait struct {
	Value      float64  `json:"value"`      // -1.0 to 1.0
	Confidence float64  `json:"confidence"` // 0.0 to 1.0
	Evidence   []string `json:"evidence"`
}

// TopicManager handles conversation flow and context
// type TopicManager struct {
//     CurrentTopic string
//     TopicHistory []TopicTransition
//     Connections  map[string][]string
// }

// Add this type definition before the TraitSystem struct
type TraitPattern struct {
	Name       string             `json:"name"`
	Indicators []string           `json:"indicators"`
	Weight     float64            `json:"weight"`
	Conditions map[string]bool    `json:"conditions"`
	Modifiers  map[string]float64 `json:"modifiers"`
	MinMatches int                `json:"min_matches"`
}
