package model

// Score - score
type Score struct {
	User     User     `json:"user,omitempty"`
	Question Question `json:"question,omitempty"`
	Created  uint64   `json:"created,omitempty"`
}
