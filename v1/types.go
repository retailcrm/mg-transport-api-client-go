package v1

import (
	"net/http"
	"time"
)

const (
	// ChannelFeatureNone channel can not implement feature
	ChannelFeatureNone string = "none"
	// ChannelFeatureReceive channel implement feature on receive
	ChannelFeatureReceive string = "receive"
	// ChannelFeatureSend channel implement feature on send
	ChannelFeatureSend string = "send"
	// ChannelFeatureBoth channel implement feature on both directions
	ChannelFeatureBoth string = "both"
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
	Name     string          `json:"name,omitempty"`
	Settings ChannelSettings `json:"settings,omitempty,brackets"`
}

// ChannelSettings struct
type ChannelSettings struct {
	SpamAllowed bool                `json:"spam_allowed"`
	Status      Status              `json:"status"`
	Text        ChannelSettingsText `json:"text"`
}

// Status struct
type Status struct {
	Delivered string `json:"delivered"`
	Read      string `json:"read"`
}

// ChannelSettingsText struct
type ChannelSettingsText struct {
	Creating string `json:"creating"`
	Editing  string `json:"editing"`
	Quoting  string `json:"quoting"`
	Deleting string `json:"deleting"`
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
	Email      string `json:"email,omitempty"`
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
	Message        Message                  `json:"message"`
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
	ExternalUserID    string              `json:"external_user_id"`
	ExternalMessageID string              `json:"external_message_id,omitempty"`
	ExternalChatID    string              `json:"external_chat_id"`
	ChannelID         uint64              `json:"channel_id"`
	Content           string              `json:"content"`
	QuoteExternalID   string              `json:"quote_external_id,omitempty"`
	QuoteContent      string              `json:"quote_content,omitempty"`
	User              *MessageDataUser    `json:"user,omitempty"`
	Bot               *MessageDataBot     `json:"bot,omitempty"`
	Product           *MessageDataProduct `json:"product,omitempty"`
	Order             *MessageDataOrder   `json:"order,omitempty"`
}

// MessageDataUser user data from webhook
type MessageDataUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

// MessageDataBot bot data from webhook
type MessageDataBot struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// MessageDataProduct product data from webhook
type MessageDataProduct struct {
	ID       uint64                    `json:"id"`
	Name     string                    `json:"name"`
	Article  string                    `json:"article,omitempty"`
	Url      string                    `json:"url,omitempty"`
	Img      string                    `json:"img,omitempty"`
	Cost     *MessageDataOrderCost     `json:"cost,omitempty"`
	Quantity *MessageDataOrderQuantity `json:"quantity,omitempty"`
}

// MessageDataOrder order data from webhook
type MessageDataOrder struct {
	Number string                  `json:"number"`
	Url    string                  `json:"url,omitempty"`
	Date   string                  `json:"date,omitempty"`
	Cost   *MessageDataOrderCost   `json:"cost,omitempty"`
	Status *MessageDataOrderStatus `json:"status,omitempty"`
	Items  []MessageDataOrderItem  `json:"items,omitempty"`
}

type MessageDataOrderStatus struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

type MessageDataOrderItem struct {
	Name     string `json:"name,omitempty"`
	Url      string `json:"url,omitempty"`
	Quantity string `json:"quantity,omitempty"`
	Price    string `json:"price,omitempty"`
}

type MessageDataOrderCost struct {
	Value    float32 `json:"value,omitempty"`
	Currency string  `json:"currency"`
}

type MessageDataOrderQuantity struct {
	Value float32 `json:"value"`
	Unit  string  `json:"unit"`
}

// TransportRequestMeta request metadata
type TransportRequestMeta struct {
	ID        uint64 `json:"id"`
	Timestamp int64  `json:"timestamp"`
}
