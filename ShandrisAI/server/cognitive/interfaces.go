package cognitive

// MoodEngine interface
type MoodProcessor interface {
	UpdateMood(context map[string]any) error
	GetCurrentMood() MoodState
	ProjectMoodInfluence(input string) float64
	RecordMoodTransition(from, to MoodState, trigger string)
}

// TraitSystem interface
type TraitProcessor interface {
	InferTraits(input string) map[string]Trait
	UpdateTraitConfidence(traitName string, evidence string)
	GetDominantTraits() []Trait
	MatchTraitPatterns(input string) []TraitPattern
}

// TopicManager interface
type TopicProcessor interface {
	SwitchTopic(newTopic string, trigger string) error
	GetTopicContext() map[string]any
	TrackTopicFlow(input string) []string
	SuggestRelatedTopics() []string
}

// Timeline interface
type TimelineProcessor interface {
	RecordEvent(event TimelineEvent) error
	GetRelevantEvents(context map[string]any) []TimelineEvent
	SetMarker(marker TimeMarker) error
	CheckRecurringEvents() []TimeMarker
}

// PersonaManager interface
type PersonaProcessor interface {
	SwitchPersona(personaID string) error
	GetCurrentPersona() Persona
	SuggestPersonaTransition(context map[string]any) *Persona
	ApplyPersonaStyle(input string) string
}
