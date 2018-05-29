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
	ID     uint64   `url:"id,omitempty"`
	Type   string   `url:"type,omitempty"`
	Events []string `url:"events,omitempty,brackets"`
}

// ActivateResponse channel activation response
type ActivateResponse struct {
	ChannelID   uint64    `json:"id"`
	ActivatedAt time.Time `json:"activated_at"`
}

// UpdateResponse channel update response
type UpdateResponse struct {
	ChannelID uint64    `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DeleteResponse channel deactivation response
type DeleteResponse struct {
	ChannelID    uint64    `json:"id"`
	DectivatedAt time.Time `json:"deactivated_at"`
}

// User struct
type User struct {
	ExternalID string `url:"external_id" json:"external_id"`
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
	ExternalID string `url:"external_id" json:"external_id"`
	Type       string `url:"type,omitempty"`
	Text       string `url:"text,omitempty"`
}

// SendMessage struct
type SendMessage struct {
	Message
	SentAt time.Time `url:"sent_at,omitempty"`
}

// UpdateMessage struct
type UpdateMessage struct {
	Message
	EditedAt time.Time `url:"edited_at,omitempty"`
}

// SendData struct
type SendData struct {
	Message SendMessage `url:"message"`
	User    User        `url:"user"`
	Channel uint64      `url:"channel"`
}

// UpdateData struct
type UpdateData struct {
	Message UpdateMessage `url:"message"`
	Channel uint64        `url:"channel"`
}

// DeleteData struct
type DeleteData struct {
	Message Message `url:"message"`
	Channel uint64  `url:"channel"`
}

// MessagesResponse message event response
type MessagesResponse struct {
	MessageID int       `json:"message_id"`
	Time      time.Time `json:"time"`
}

// Webhook request
type WebhookRequest struct {
	Type string               `json:"type"`
	Meta TransportRequestMeta `json:"meta"`
	Data WebhookData          `json:"data"`
}

// WebhookData request data
type WebhookData struct {
	ExternalUserID         string `json:"external_user_id"`
	ExternalMessageID      string `json:"external_message_id,omitempty"`
	ChannelID              uint64 `json:"channel_id"`
	ChatID                 int64  `json:"chat_id"` // TODO: add to webhook request
	Content                string `json:"content"`
	QuoteMessageExternalID string `json:"quote_message_external_id,omitempty"`
}

type TransportRequestMeta struct {
	ID        uint64 `json:"id"`
	Timestamp int64  `json:"timestamp"`
}
