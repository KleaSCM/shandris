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
