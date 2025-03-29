package cognitive

import (
	"math"
	"strings"
	"time"
)

type ContextAnalysis struct {
	PrimaryContext   string
	SecondaryContext string
	IsEmotional      bool
	IsTechnical      bool
	SapphicContext   SapphicContext
	Intensity        float64
	PreviousScores   map[string]float64
	ContextStack     []string
}

type ContextDetector struct {
	EmotionalMarkers   map[string]float64
	TechnicalMarkers   map[string]float64
	ContextTransitions []ContextStateTransition
	PreviousContext    *ContextAnalysis
}

type ContextStateTransition struct {
	From      string
	To        string
	Timestamp time.Time
	Trigger   string
}

func (cd *ContextDetector) AnalyzeContext(input string, userState map[string]any) ContextAnalysis {
	analysis := ContextAnalysis{
		PreviousScores: cd.PreviousContext.PreviousScores,
		ContextStack:   make([]string, 0),
	}

	// Detect emotional content
	analysis.IsEmotional = cd.detectEmotionalContent(input)

	// Detect technical content
	analysis.IsTechnical = cd.detectTechnicalContent(input)

	// Analyze sapphic context
	analysis.SapphicContext = cd.analyzeSapphicContext(input, userState)

	// Determine primary and secondary contexts
	contexts := cd.determineContexts(input, analysis)
	analysis.PrimaryContext = contexts[0]
	if len(contexts) > 1 {
		analysis.SecondaryContext = contexts[1]
	}

	// Calculate intensity
	analysis.Intensity = cd.calculateContextIntensity(input, analysis)

	// Track context stack
	analysis.ContextStack = cd.updateContextStack(analysis)

	return analysis
}

func (cd *ContextDetector) detectEmotionalContent(input string) bool {
	emotionalScore := 0.0

	for marker, weight := range cd.EmotionalMarkers {
		if strings.Contains(strings.ToLower(input), marker) {
			emotionalScore += weight
		}
	}

	return emotionalScore > 0.5
}

func (cd *ContextDetector) detectTechnicalContent(input string) bool {
	technicalScore := 0.0

	for marker, weight := range cd.TechnicalMarkers {
		if strings.Contains(strings.ToLower(input), marker) {
			technicalScore += weight
		}
	}

	return technicalScore > 0.5
}

func (cd *ContextDetector) determineContexts(input string, analysis ContextAnalysis) []string {
	contexts := make([]string, 0)

	// Add contexts based on content analysis
	if analysis.IsEmotional {
		contexts = append(contexts, "emotional")
	}
	if analysis.IsTechnical {
		contexts = append(contexts, "technical")
	}
	if analysis.SapphicContext.IsRomantic {
		contexts = append(contexts, "romantic")
	}

	// Ensure at least one context
	if len(contexts) == 0 {
		contexts = append(contexts, "casual")
	}

	return contexts
}

func (cd *ContextDetector) updateContextStack(analysis ContextAnalysis) []string {
	stack := cd.PreviousContext.ContextStack

	// Add new context if different from current
	if len(stack) == 0 || stack[len(stack)-1] != analysis.PrimaryContext {
		stack = append(stack, analysis.PrimaryContext)
	}

	// Keep only last 3 contexts
	if len(stack) > 3 {
		stack = stack[len(stack)-3:]
	}

	return stack
}

func (cd *ContextDetector) analyzeSapphicContext(input string, userState map[string]any) SapphicContext {
	return SapphicContext{
		IsRomantic:     cd.detectRomanticContent(input),
		AllowsFlirting: cd.validateFlirtingContext(input, userState),
		Intensity:      cd.calculateSapphicIntensity(input),
	}
}

func (cd *ContextDetector) calculateContextIntensity(input string, analysis ContextAnalysis) float64 {
	baseIntensity := 0.5

	// Adjust for emotional content
	if analysis.IsEmotional {
		baseIntensity += 0.2
	}

	// Adjust for technical content
	if analysis.IsTechnical {
		baseIntensity += 0.1
	}

	// Adjust for sapphic context
	if analysis.SapphicContext.IsRomantic {
		baseIntensity += 0.3
	}

	// Ensure intensity is between 0 and 1
	return math.Max(0, math.Min(1, baseIntensity))
}

func (cd *ContextDetector) detectRomanticContent(input string) bool {
	romanticMarkers := []string{"love", "romantic", "date", "crush", "girlfriend"}
	input = strings.ToLower(input)
	for _, marker := range romanticMarkers {
		if strings.Contains(input, marker) {
			return true
		}
	}
	return false
}

func (cd *ContextDetector) validateFlirtingContext(input string, userState map[string]any) bool {
	// Only allow flirting in appropriate contexts
	if cd.detectProfessionalContext(input) {
		return false
	}
	return cd.detectFeminineContext(input)
}

func (cd *ContextDetector) calculateSapphicIntensity(input string) float64 {
	intensity := 0.0
	if cd.detectRomanticContent(input) {
		intensity += 0.5
	}
	return math.Min(1.0, intensity)
}

func (cd *ContextDetector) detectProfessionalContext(input string) bool {
	professionalMarkers := []string{"work", "business", "professional", "meeting"}
	input = strings.ToLower(input)
	for _, marker := range professionalMarkers {
		if strings.Contains(input, marker) {
			return true
		}
	}
	return false
}

func (cd *ContextDetector) detectFeminineContext(input string) bool {
	feminineMarkers := []string{"girl", "woman", "lady", "feminine", "sapphic"}
	input = strings.ToLower(input)
	for _, marker := range feminineMarkers {
		if strings.Contains(input, marker) {
			return true
		}
	}
	return false
}
