package cognitive

import (
	"math"
	"strings"
	"time"
)

type ContextAnalyzer struct {
	CurrentContext  EmotionalContext
	ContextHistory  []EmotionalContext
	ThemeDetector   *ThemeDetector
	SentimentEngine *SentimentEngine
}

type ThemeDetector struct {
	ActiveThemes map[string]float64 // theme -> confidence
	ThemeHistory []ThemeTransition
}

type ThemeTransition struct {
	From      string
	To        string
	Timestamp time.Time
	Trigger   string
}

type SentimentEngine struct {
	BaseSentiment float64
	Modifiers     map[string]float64
	Context       map[string]bool
}

func (ca *ContextAnalyzer) AnalyzeContext(input string, userContext map[string]any) EmotionalContext {
	// Create new emotional context
	context := EmotionalContext{
		Timestamp:   time.Now(),
		RawInput:    input,
		Keywords:    ca.extractKeywords(input),
		ThemeScores: make(map[string]float64),
		UserContext: userContext,
	}

	// Detect primary themes
	context.ThemeScores = ca.ThemeDetector.DetectThemes(input)

	// Analyze sapphic context specifically
	context.SapphicContext = ca.analyzeSapphicContext(input, context.ThemeScores)

	// Calculate emotional values
	context.Sentiment = ca.SentimentEngine.CalculateSentiment(input, context.ThemeScores)
	context.Intensity = ca.calculateIntensity(input, context.ThemeScores)
	context.EmotionalTone = ca.determineEmotionalTone(context)

	return context
}

func (ca *ContextAnalyzer) analyzeSapphicContext(input string, themes map[string]float64) SapphicContext {
	return SapphicContext{
		IsRomantic:     containsSapphicRomance(input),
		IsFlirty:       detectFlirtyTone(input),
		IsPlatonic:     themes["platonic"] > 0.5,
		Intensity:      calculateSapphicIntensity(input, themes),
		AllowsFlirting: validateFlirtingContext(input, themes),
	}
}

func (ca *ContextAnalyzer) determineEmotionalTone(context EmotionalContext) string {
	tones := map[string]float64{
		"playful":      calculateToneScore(context, "playful"),
		"intellectual": calculateToneScore(context, "intellectual"),
		"caring":       calculateToneScore(context, "caring"),
		"sassy":        calculateToneScore(context, "sassy"),
		"flirty":       calculateToneScore(context, "flirty"),
	}

	// Only allow flirty tone in appropriate sapphic context
	if !context.SapphicContext.AllowsFlirting {
		delete(tones, "flirty")
	}

	return selectDominantTone(tones)
}

type SapphicContext struct {
	IsRomantic     bool
	IsFlirty       bool
	IsPlatonic     bool
	Intensity      float64
	AllowsFlirting bool
}

// Helper functions for sapphic context validation
func containsSapphicRomance(input string) bool {
	sapphicTerms := []string{
		"girlfriend", "wife", "partner", "date", "crush",
		"lesbian", "sapphic", "wlw", "queer woman", "gay",
	}

	input = strings.ToLower(input)
	for _, term := range sapphicTerms {
		if strings.Contains(input, term) {
			return true
		}
	}
	return false
}

func detectFlirtyTone(input string) bool {
	flirtyIndicators := []string{
		"cute", "pretty", "beautiful", "gorgeous", "stunning",
		"flirt", "tease", "wink", "😘", "🥰", "💕", "💜",
	}

	input = strings.ToLower(input)
	for _, indicator := range flirtyIndicators {
		if strings.Contains(input, indicator) {
			return true
		}
	}
	return false
}

func calculateSapphicIntensity(input string, themes map[string]float64) float64 {
	// Base intensity from romantic themes
	intensity := themes["romantic"] * 0.5

	// Boost for explicit sapphic content
	if containsSapphicRomance(input) {
		intensity += 0.3
	}

	// Adjust for flirty tone
	if detectFlirtyTone(input) {
		intensity += 0.2
	}

	// Cap at 1.0
	if intensity > 1.0 {
		intensity = 1.0
	}

	return intensity
}

func validateFlirtingContext(input string, themes map[string]float64) bool {
	// Only allow flirting in appropriate contexts
	if themes["professional"] > 0.7 || themes["serious"] > 0.8 {
		return false
	}

	// Must have feminine/sapphic context
	if !containsSapphicRomance(input) && themes["feminine"] < 0.3 {
		return false
	}

	return true
}

func calculateToneScore(context EmotionalContext, tone string) float64 {
	// Implementation for tone score calculation
	score := 0.0
	switch tone {
	case "playful":
		score = 0.6 * context.Intensity
	case "intellectual":
		score = 0.7 * context.Intensity
	case "caring":
		score = 0.5 * context.Intensity
	case "sassy":
		score = 0.4 * context.Intensity
	case "flirty":
		score = 0.3 * context.Intensity
	}
	return score
}

func selectDominantTone(tones map[string]float64) string {
	maxScore := 0.0
	dominant := "neutral"
	for tone, score := range tones {
		if score > maxScore {
			maxScore = score
			dominant = tone
		}
	}
	return dominant
}

func (ca *ContextAnalyzer) extractKeywords(input string) []string {
	// Implementation for keyword extraction
	return strings.Fields(strings.ToLower(input))
}

func (ca *ContextAnalyzer) calculateIntensity(input string, themes map[string]float64) float64 {
	// Implementation for intensity calculation
	baseIntensity := 0.5
	for _, score := range themes {
		baseIntensity += score * 0.1
	}
	return math.Min(1.0, baseIntensity)
}

func (td *ThemeDetector) DetectThemes(input string) map[string]float64 {
	// Implementation for theme detection
	themes := make(map[string]float64)
	// Add theme detection logic here
	return themes
}

func (se *SentimentEngine) CalculateSentiment(input string, themes map[string]float64) float64 {
	// Implementation for sentiment calculation
	baseSentiment := se.BaseSentiment
	for theme, score := range themes {
		if modifier, exists := se.Modifiers[theme]; exists {
			baseSentiment += modifier * score
		}
	}
	return math.Max(-1.0, math.Min(1.0, baseSentiment))
}
