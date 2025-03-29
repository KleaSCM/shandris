package cognitive

import (
	"time"
)

// Timeline manages long-term memory and event tracking
type Timeline struct {
	Events    []TimelineEvent
	Markers   map[string]TimeMarker
	Relations map[string][]string // Maps events to related events
}

type TimelineEvent struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"` // conversation, milestone, mood_shift, etc
	Timestamp  time.Time `json:"timestamp"`
	Content    any       `json:"content"`
	Importance float64   `json:"importance"` // 0.0 to 1.0
	Tags       []string  `json:"tags"`
}

type TimeMarker struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"` // e.g., "first_meeting", "learned_go"
	Timestamp    time.Time `json:"timestamp"`
	Recurring    bool      `json:"recurring"` // for anniversaries etc
	LastRecalled time.Time `json:"last_recalled"`
}
