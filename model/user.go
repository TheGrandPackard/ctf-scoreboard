package model

// User - user
type User struct {
	ID           uint64 `json:"id,omitempty"`
	Admin        bool   `json:"admin,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	Team         Team   `json:"team,omitempty"`
	Created      uint64 `json:"created,omitempty"`
}
