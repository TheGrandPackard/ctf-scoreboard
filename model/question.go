package model

// Question - question
type Question struct {
	ID       uint64   `json:"id,omitempty"`
	Category Category `json:"category,omitempty"`
	Team     Team     `json:"team,omitempty"`
	Text     string   `json:"text,omitempty"`
	Answer   string   `json:"answer,omitempty"`
	Hint     string   `json:"hint,omitempty"`
	File     string   `json:"file,omitempty"`
	Points   uint64   `json:"points,omitempty"`
	Created  uint64   `json:"created,omitempty"`
}
