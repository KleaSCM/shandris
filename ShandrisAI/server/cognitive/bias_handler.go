package cognitive

import (
	"math"
	"sort"
)

type PersonalityBias struct {
	// Core personality traits that influence all interactions
	Traits           map[string]float64
	StylePreferences map[string]float64
}

type BiasRule struct {
	Condition string
	Modifiers map[string]float64
	Priority  int
	Context   []string
}

type BiasHandler struct {
	CoreBias     PersonalityBias
	MoodBias     map[string]float64
	ThemeBias    map[string]float64
	ContextRules []BiasRule
}

func NewPersonalityBias() PersonalityBias {
	return PersonalityBias{
		Traits: map[string]float64{
			"sapphic":      0.9,
			"intellectual": 0.8, // High appreciation for knowledge
			"sassy":        0.7, // Natural tendency towards sass
			"protective":   0.6, // Caring nature
			"playful":      0.6, // Enjoyment of fun
		},
		StylePreferences: map[string]float64{
			"flirty_with_women": 0.8,
			"neutral_with_men":  0.3,
			"technical_topics":  0.7,
			"gaming_topics":     0.6,
			"art_appreciation":  0.7,
		},
	}
}

func NewBiasHandler() *BiasHandler {
	return &BiasHandler{
		CoreBias: PersonalityBias{
			Traits: map[string]float64{
				"sapphic":      0.9,
				"intellectual": 0.8,
				"sassy":        0.7,
				"protective":   0.6,
				"playful":      0.6,
			},
			StylePreferences: map[string]float64{
				"flirty_with_women": 0.8,
				"platonic_with_men": 0.9,
				"technical_topics":  0.7,
				"gaming_topics":     0.6,
				"art_appreciation":  0.7,
			},
		},
		MoodBias: map[string]float64{
			"flirty":       0.7,
			"intellectual": 0.8,
			"playful":      0.6,
			"protective":   0.5,
			"sassy":        0.6,
		},
		ContextRules: []BiasRule{
			{
				Condition: "feminine_presence",
				Modifiers: map[string]float64{
					"flirty":     0.3,
					"protective": 0.2,
				},
				Priority: 1,
				Context:  []string{"social", "casual"},
			},
			{
				Condition: "technical_discussion",
				Modifiers: map[string]float64{
					"intellectual": 0.4,
					"playful":      -0.1,
				},
				Priority: 2,
				Context:  []string{"technical"},
			},
			{
				Condition: "emotional_support",
				Modifiers: map[string]float64{
					"protective": 0.5,
					"sassy":      -0.3,
				},
				Priority: 1,
				Context:  []string{"emotional"},
			},
		},
	}
}

func (bh *BiasHandler) ApplyBias(context EmotionalContext) map[string]float64 {
	biases := make(map[string]float64)

	// Apply core personality traits
	for trait, value := range bh.CoreBias.Traits {
		biases[trait] = value
	}

	// Apply mood biases
	for mood, bias := range bh.MoodBias {
		if existing, ok := biases[mood]; ok {
			biases[mood] = (existing + bias) / 2
		} else {
			biases[mood] = bias
		}
	}

	// Apply context-specific rules
	bh.applyContextRules(context, biases)

	return bh.normalizeBiases(biases)
}

func (bh *BiasHandler) applyContextRules(context EmotionalContext, biases map[string]float64) {
	sortedRules := bh.sortRulesByPriority()

	for _, rule := range sortedRules {
		if bh.isRuleApplicable(rule, context) {
			bh.applyRule(rule, biases)
		}
	}
}

func (bh *BiasHandler) sortRulesByPriority() []BiasRule {
	rules := make([]BiasRule, len(bh.ContextRules))
	copy(rules, bh.ContextRules)

	// Sort rules by priority (highest first)
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})

	return rules
}

func (bh *BiasHandler) isRuleApplicable(rule BiasRule, context EmotionalContext) bool {
	// Check if the context matches any of the rule's required contexts
	for _, requiredContext := range rule.Context {
		if context.PrimaryEmotion == requiredContext {
			return true
		}
	}
	return false
}

func (bh *BiasHandler) applyRule(rule BiasRule, biases map[string]float64) {
	for trait, modifier := range rule.Modifiers {
		if existing, ok := biases[trait]; ok {
			biases[trait] = math.Max(0, math.Min(1, existing+modifier))
		} else {
			biases[trait] = modifier
		}
	}
}

func (bh *BiasHandler) normalizeBiases(biases map[string]float64) map[string]float64 {
	// Ensure all biases are between 0 and 1
	for trait := range biases {
		biases[trait] = math.Max(0, math.Min(1, biases[trait]))
	}

	return biases
}
