package cognitive

import "time"

// MoodPattern defines emotional triggers and their impacts
type MoodPattern struct {
	Keywords     []string // Trigger words/phrases
	Sentiment    float64  // -1.0 to 1.0
	MoodShift    string   // Target mood
	Intensity    float64  // Impact strength
	Decay        float64  // How fast this impact fades
	Requirements []string // Required context for this mood
	Exclusions   []string // Contexts that prevent this mood
}

// EmotionalContext captures conversation tone and impact
type EmotionalContext struct {
	PrimaryEmotion string
	Sentiment      float64
	Intensity      float64
	Keywords       []string
	UserMood       string
	PriorContext   *EmotionalContext
	SapphicContext SapphicContext
	IsEmotional    bool
	IsTechnical    bool
	Timestamp      time.Time
	RawInput       string
	ThemeScores    map[string]float64
	UserContext    map[string]any
	EmotionalTone  string
}

// TopicManager handles conversation flow and context
type TopicManager struct {
	CurrentTopic string
	TopicHistory []TopicTransition
	Connections  map[string][]string
	ActiveTopics []string
	TopicGraph   map[string]*TopicNode
	DomainRules  map[string]DomainRule
}

type TopicTransition struct {
	From      string
	To        string
	Timestamp time.Time
	Trigger   string
}

type TopicNode struct {
	ID         string
	Name       string
	Category   string
	Keywords   []string
	LastActive time.Time
	Frequency  int
	Confidence float64
	Relations  map[string]float64
	Context    map[string]interface{}
}

type DomainRule struct {
	Domain      string
	Keywords    []string
	Validators  []func(string) bool
	Priority    int
	Transitions map[string]float64
}
