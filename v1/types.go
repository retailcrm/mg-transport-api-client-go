package v1

import (
	"net/http"
	"time"
)

const (
	ChannelFeatureNone    string = "none"
	ChannelFeatureReceive string = "receive"
	ChannelFeatureSend    string = "send"
	ChannelFeatureBoth    string = "both"

	MsgModeNever       string = "never"
	MsgModeAlways      string = "always"
)

// MgClient type
type MgClient struct {
	URL        string       `json:"url"`
	Token      string       `json:"token"`
	Debug      bool         `json:"debug"`
	httpClient *http.Client `json:"-"`
}

// Channel type
type Channel struct {
	ID       uint64          `json:"id,omitempty"`
	Type     string          `json:"type,omitempty"`
	Events   []string        `json:"events,omitempty,brackets"`
	Settings ChannelSettings `json:"settings,omitempty,brackets"`
}

// ChannelSettings struct
type ChannelSettings struct {
	Features           ChannelFeatures `json:"features"`
	ReceiveMessageMode string          `json:"receive_message_mode"`
	SpamAllowed        bool            `json:"spam_allowed"`
}

// ChannelFeatures struct
type ChannelFeatures struct {
	StatusDelivered string `json:"status_delivered"`
	StatusRead      string `json:"status_read"`
	MessageDeleting string `json:"message_deleting"`
	MessageEditing  string `json:"message_editing"`
	MessageQuoting  string `json:"message_quoting"`
	ImageMessage    string `json:"image_message"`
	FileMessage     string `json:"file_message"`
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
	ExternalID string `json:"external_id"`
	Nickname   string `json:"nickname"`
	Firstname  string `json:"first_name,omitempty"`
	Lastname   string `json:"last_name,omitempty"`
	Avatar     string `json:"avatar,omitempty"`
	ProfileURL string `json:"profile_url,omitempty"`
	Country    string `json:"country,omitempty"`
	Language   string `json:"language,omitempty"`
	Phone      string `json:"phone,omitempty"`
}

// Message struct
type Message struct {
	ExternalID string `json:"external_id"`
	Type       string `json:"type,omitempty"`
	Text       string `json:"text,omitempty"`
}

// SendMessage struct
type SendMessage struct {
	Message
	SentAt time.Time `json:"sent_at,omitempty"`
}

// UpdateMessage struct
type UpdateMessage struct {
	Message
	EditedAt int64 `json:"edited_at,omitempty"`
}

// SendData struct
type SendData struct {
	Message        SendMessage              `json:"message"`
	User           User                     `json:"user"`
	Channel        uint64                   `json:"channel"`
	ExternalChatID string                   `json:"external_chat_id"`
	Quote          *SendMessageRequestQuote `json:"quote,omitempty"`
}

type SendMessageRequestQuote struct {
	ExternalID string `json:"external_id"`
}

// UpdateData struct
type UpdateData struct {
	Message UpdateMessage `json:"message"`
	Channel uint64        `json:"channel"`
}

// DeleteData struct
type DeleteData struct {
	Message Message `json:"message"`
	Channel uint64  `json:"channel"`
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
	ExternalUserID    string `json:"external_user_id"`
	ExternalMessageID string `json:"external_message_id,omitempty"`
	ExternalChatID    string `json:"external_chat_id"`
	ChannelID         uint64 `json:"channel_id"`
	Content           string `json:"content"`
	QuoteExternalID   string `json:"quote_external_id,omitempty"`
	QuoteContent      string `json:"quote_content,omitempty"`
}

// TransportRequestMeta request metadata
type TransportRequestMeta struct {
	ID        uint64 `json:"id"`
	Timestamp int64  `json:"timestamp"`
}
