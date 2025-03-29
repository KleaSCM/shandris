package cognitive

import (
	"time"
)

// AdvancedMoodPattern extends the basic MoodPattern with more sophisticated controls
type AdvancedMoodPattern struct {
	Base         MoodPattern
	Transitions  map[string]MoodTransitionRule
	Combinations []CombinedPattern
	Context      ContextualBehavior
	Modifiers    []MoodModifier
	Cooldowns    map[string]time.Duration
}

type MoodTransitionRule struct {
	ToMood      string
	Conditions  []string
	Probability float64
	MinDuration time.Duration
	MaxDuration time.Duration
	Smoothing   float64
}

type CombinedPattern struct {
	Patterns []string
	Weight   float64
	Synergy  float64
	Duration time.Duration
}

type ContextualBehavior struct {
	Triggers     map[string]float64
	Inhibitors   map[string]float64
	Requirements map[string]float64
	TimeFactors  map[string]TimeInfluence
}

type TimeInfluence struct {
	Peak        time.Duration
	Decay       float64
	Periodicity time.Duration
}

type MoodModifier struct {
	Type      string
	Strength  float64
	Duration  time.Duration
	Stackable bool
}

// Initialize advanced mood patterns
func initializeAdvancedMoodPatternsExtended() map[string]AdvancedMoodPattern {
	return map[string]AdvancedMoodPattern{
		"sapphic_romantic": {
			Base: MoodPattern{
				Keywords: []string{
					"romantic", "intimate", "tender", "gentle", "soft",
					"loving", "affectionate", "warm", "close", "sweet",
				},
				Sentiment:    0.9,
				MoodShift:    "romantic",
				Intensity:    0.8,
				Decay:        0.05,
				Requirements: []string{"feminine_presence", "intimate_context"},
			},
			Transitions: map[string]MoodTransitionRule{
				"playful": {
					Conditions:  []string{"positive_response", "mutual_comfort"},
					Probability: 0.7,
					MinDuration: 2 * time.Minute,
					Smoothing:   0.3,
				},
				"protective": {
					Conditions:  []string{"vulnerability_expressed", "trust_established"},
					Probability: 0.8,
					MinDuration: 5 * time.Minute,
					Smoothing:   0.4,
				},
			},
			Combinations: []CombinedPattern{
				{
					Patterns: []string{"tender_care", "gentle_flirt"},
					Weight:   0.9,
					Synergy:  1.2,
					Duration: 10 * time.Minute,
				},
			},
			Context: ContextualBehavior{
				Triggers: map[string]float64{
					"shared_interests":   0.8,
					"emotional_openness": 0.9,
					"mutual_attraction":  1.0,
				},
				Inhibitors: map[string]float64{
					"professional_context": 0.8,
					"public_setting":       0.6,
					"negative_mood":        0.7,
				},
			},
		},
		"tech_mentor": {
			Base: MoodPattern{
				Keywords: []string{
					"explain", "teach", "guide", "help", "understand",
					"learn", "explore", "discover", "solve", "debug",
				},
				Sentiment: 0.7,
				MoodShift: "supportive",
				Intensity: 0.6,
				Decay:     0.1,
			},
			Transitions: map[string]MoodTransitionRule{
				"excited": {
					Conditions:  []string{"breakthrough", "understanding_achieved"},
					Probability: 0.8,
					MinDuration: 1 * time.Minute,
				},
				"encouraging": {
					Conditions:  []string{"struggle_detected", "persistence_shown"},
					Probability: 0.9,
					MinDuration: 3 * time.Minute,
				},
			},
			Context: ContextualBehavior{
				Triggers: map[string]float64{
					"technical_question":   0.9,
					"learning_opportunity": 0.8,
					"problem_solving":      0.7,
				},
			},
		},
		"playful_sass": {
			Base: MoodPattern{
				Keywords: []string{
					"tease", "sass", "witty", "clever", "playful",
					"banter", "joke", "quip", "flirt", "charm",
				},
				Sentiment: 0.8,
				MoodShift: "sassy",
				Intensity: 0.7,
				Decay:     0.15,
			},
			Transitions: map[string]MoodTransitionRule{
				"flirty": {
					Conditions:  []string{"sapphic_context", "mutual_interest"},
					Probability: 0.7,
					MinDuration: 2 * time.Minute,
				},
			},
			Modifiers: []MoodModifier{
				{
					Type:      "confidence_boost",
					Strength:  0.3,
					Duration:  5 * time.Minute,
					Stackable: true,
				},
			},
		},
		// Add more patterns as needed...
	}
}
