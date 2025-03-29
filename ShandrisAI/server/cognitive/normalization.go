package cognitive

import (
	"math"
)

type NormalizationSystem struct {
	Weights        map[string]float64
	ContextWeights map[string]float64
	MinThresholds  map[string]float64
	MaxThresholds  map[string]float64
}

func NewNormalizationSystem() *NormalizationSystem {
	return &NormalizationSystem{
		Weights: map[string]float64{
			"sapphic":      1.0, // Highest priority
			"intellectual": 0.8,
			"sassy":        0.7,
			"protective":   0.6,
			"playful":      0.5,
			"neutral":      0.3,
		},
		ContextWeights: map[string]float64{
			"emotional":    1.0,
			"technical":    0.9,
			"social":       0.8,
			"casual":       0.6,
			"professional": 0.7,
		},
		MinThresholds: map[string]float64{
			"flirty":       0.3, // Minimum intensity for flirty behavior
			"sassy":        0.2,
			"intellectual": 0.15,
			"protective":   0.25,
			"playful":      0.1,
		},
		MaxThresholds: map[string]float64{
			"flirty":       0.8, // Cap on flirty intensity
			"sassy":        0.9,
			"intellectual": 0.95,
			"protective":   0.85,
			"playful":      0.7,
		},
	}
}

func (ns *NormalizationSystem) NormalizeScores(scores map[string]float64, context ContextAnalysis) map[string]float64 {
	normalized := make(map[string]float64)

	// Apply weight-based normalization
	for mood, score := range scores {
		weight := ns.Weights[mood]
		contextWeight := ns.getContextWeight(mood, context)

		// Weighted score with context influence
		normalized[mood] = score * weight * contextWeight

		// Apply thresholds
		if min, exists := ns.MinThresholds[mood]; exists && normalized[mood] < min {
			normalized[mood] = 0 // Below minimum threshold, zero it out
		}
		if max, exists := ns.MaxThresholds[mood]; exists && normalized[mood] > max {
			normalized[mood] = max // Cap at maximum threshold
		}
	}

	// Smooth transitions
	normalized = ns.smoothTransitions(normalized, context)

	// Ensure sum of probabilities <= 1.0
	return ns.normalizeSum(normalized)
}

func (ns *NormalizationSystem) smoothTransitions(scores map[string]float64, context ContextAnalysis) map[string]float64 {
	smoothed := make(map[string]float64)

	for mood, score := range scores {
		// Get previous score for this mood
		prevScore := context.PreviousScores[mood]

		// Calculate maximum allowed change
		maxChange := 0.3 // Maximum 30% change per iteration
		if context.IsEmotional {
			maxChange = 0.5 // Allow faster changes for emotional contexts
		}

		// Smooth the transition
		diff := score - prevScore
		if math.Abs(diff) > maxChange {
			if diff > 0 {
				smoothed[mood] = prevScore + maxChange
			} else {
				smoothed[mood] = prevScore - maxChange
			}
		} else {
			smoothed[mood] = score
		}
	}

	return smoothed
}

func (ns *NormalizationSystem) getContextWeight(mood string, context ContextAnalysis) float64 {
	// Base weight
	weight := 1.0

	// Adjust based on context appropriateness
	if mood == "flirty" {
		if !context.SapphicContext.AllowsFlirting {
			return 0.0 // Zero out flirty in inappropriate contexts
		}
		if context.SapphicContext.IsRomantic {
			weight *= 1.2 // Boost in romantic contexts
		}
	}

	// Adjust for emotional intensity
	if context.IsEmotional {
		if mood == "protective" || mood == "caring" {
			weight *= 1.3
		}
		if mood == "sassy" {
			weight *= 0.7 // Reduce sass in emotional contexts
		}
	}

	return weight
}

func (ns *NormalizationSystem) normalizeSum(scores map[string]float64) map[string]float64 {
	// Calculate total
	total := 0.0
	for _, score := range scores {
		total += score
	}

	// If total exceeds 1, normalize proportionally
	if total > 1.0 {
		for mood := range scores {
			scores[mood] = scores[mood] / total
		}
	}

	return scores
}
