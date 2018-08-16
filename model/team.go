package model

// Team - team
type Team struct {
	ID      uint64 `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Created uint64 `json:"created,omitempty"`
}
