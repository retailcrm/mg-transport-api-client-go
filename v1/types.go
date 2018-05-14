package v1

import (
	"net/http"
	"time"
)

// MgClient type
type MgClient struct {
	URL        string
	Token      string
	httpClient *http.Client
}

// Channel type
type Channel struct {
	ID           uint64    `url:"id,omitempty"`
	Type         string    `url:"type,omitempty"`
	Events       []string  `url:"events,omitempty,brackets"`
	CreatedAt    time.Time `url:"created_at,omitempty"`
	UpdatedAt    time.Time `url:"updated_at,omitempty"`
	ActivatedAt  time.Time `url:"activated_at,omitempty"`
	DectivatedAt time.Time `url:"deactivated_at,omitempty"`
}

// ActivateResponse channel activation response
type ActivateResponse struct {
	ChannelID   uint64    `json:"channel_id"`
	ActivatedAt time.Time `json:"activated_at"`
}

// UpdateResponse channel update response
type UpdateResponse struct {
	ChannelID uint64    `json:"channel_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DeleteResponse channel deactivation response
type DeleteResponse struct {
	ChannelID    uint64    `json:"channel_id"`
	DectivatedAt time.Time `json:"deactivated_at"`
}

type SendData struct {
	Message SendMessage `url:"message"`
	User    User        `url:"user"`
}

// User struct
type User struct {
	ExternalID string `url:"external_id"`
	Nickname   string `url:"nickname"`
	Firstname  string `url:"first_name,omitempty"`
	Lastname   string `url:"last_name,omitempty"`
	Avatar     string `url:"avatar,omitempty"`
	ProfileURL string `url:"profile_url,omitempty"`
	Country    string `url:"country,omitempty"`
	Language   string `url:"language,omitempty"`
	Phone      string `url:"phone,omitempty"`
}

// Message struct
type Message struct {
	ExternalID string `url:"external_id"`
	Channel    uint64 `url:"channel"`
	Type       string `url:"type"`
	Text       string `url:"text,omitempty"`
}

// SendMessage struct
type SendMessage struct {
	ExternalID string    `url:"external_id"`
	Channel    uint64    `url:"channel"`
	Type       string    `url:"type"`
	Text       string    `url:"text,omitempty"`
	SentAt     time.Time `url:"sent_at,omitempty"`
}

// UpdateMessage struct
type UpdateMessage struct {
	ExternalID string    `url:"external_id"`
	Channel    uint64    `url:"channel"`
	Type       string    `url:"type"`
	Text       string    `url:"text,omitempty"`
	EditedAt   time.Time `url:"edited_at,omitempty"`
}

// MessagesResponse message event response
type MessagesResponse struct {
	MessageID string    `json:"message_id"`
	Time      time.Time `json:"time"`
}
