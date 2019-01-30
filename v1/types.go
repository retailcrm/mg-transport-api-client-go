package v1

import (
	"net/http"
	"time"
)

//noinspection ALL
const (
	// ChannelFeatureNone channel can not implement feature
	ChannelFeatureNone string = "none"
	// ChannelFeatureReceive channel implement feature on receive
	ChannelFeatureReceive string = "receive"
	// ChannelFeatureSend channel implement feature on send
	ChannelFeatureSend string = "send"
	// ChannelFeatureBoth channel implement feature on both directions
	ChannelFeatureBoth string = "both"

	// MsgTypeText text message
	MsgTypeText string = "text"
	// MsgTypeSystem system message
	MsgTypeSystem string = "system"
	// MsgTypeCommand command (for bots)
	MsgTypeCommand string = "command"
	// MsgTypeOrder order card
	MsgTypeOrder string = "order"
	// MsgTypeProduct product card
	MsgTypeProduct string = "product"
	// MsgTypeFile file card
	MsgTypeFile string = "file"
	// MsgTypeImage image card
	MsgTypeImage string = "image"

	// MsgOrderStatusCodeNew order status group new
	MsgOrderStatusCodeNew = "new"
	// MsgOrderStatusCodeApproval order status group approval
	MsgOrderStatusCodeApproval = "approval"
	// MsgOrderStatusCodeAssembling order status group assembling
	MsgOrderStatusCodeAssembling = "assembling"
	// MsgOrderStatusCodeDelivery order status group delivery
	MsgOrderStatusCodeDelivery = "delivery"
	// MsgOrderStatusCodeComplete order status group complete
	MsgOrderStatusCodeComplete = "complete"
	// MsgOrderStatusCodeCancel order status group cancel
	MsgOrderStatusCodeCancel = "cancel"

	FileSizeLimit = 20 * 1024 * 1024
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
	SpamAllowed bool                     `json:"spam_allowed"`
	Status      Status                   `json:"status"`
	Text        ChannelSettingsText      `json:"text"`
	Product     Product                  `json:"product"`
	Order       Order                    `json:"order"`
	File        ChannelSettingsFilesBase `json:"file"`
	Image       ChannelSettingsFilesBase `json:"image"`
}

// Product type
type Product struct {
	Creating string `json:"creating,omitempty"`
	Editing  string `json:"editing,omitempty"`
	Deleting string `json:"deleting,omitempty"`
}

// Order type
type Order struct {
	Creating string `json:"creating,omitempty"`
	Editing  string `json:"editing,omitempty"`
	Deleting string `json:"deleting,omitempty"`
}

// Status struct
type Status struct {
	Delivered string `json:"delivered,omitempty"`
	Read      string `json:"read,omitempty"`
}

// ChannelSettingsText struct
type ChannelSettingsText struct {
	Creating      string `json:"creating,omitempty"`
	Editing       string `json:"editing,omitempty"`
	Quoting       string `json:"quoting,omitempty"`
	Deleting      string `json:"deleting,omitempty"`
	MaxCharsCount uint16 `json:"max_chars_count,omitempty"`
}

// ChannelSettingsFilesBase struct
type ChannelSettingsFilesBase struct {
	Creating             string `json:"creating,omitempty"`
	Editing              string `json:"editing,omitempty"`
	Quoting              string `json:"quoting,omitempty"`
	Deleting             string `json:"deleting,omitempty"`
	Max                  uint64 `json:"max_items_count,omitempty"`
	CommentAttribute     string `json:"comment_attribute,omitempty"`
	CommentMaxCharsCount int    `json:"comment_max_chars_count,omitempty"`
}

// FullFileResponse uploaded file data
type FullFileResponse struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Size int    `json:"size,omitempty"`
	Url  string `json:"url,omitempty"`
}

// UploadFileResponse uploaded file data
type UploadFileResponse struct {
	ID        string    `json:"id"`
	Hash      string    `json:"hash"`
	Type      string    `json:"type"`
	Meta      FileMeta  `json:"meta"`
	MimeType  string    `json:"mime_type"`
	Size      int       `json:"size"`
	Url       *string   `json:"source_url"`
	CreatedAt time.Time `json:"created_at"`
}

// FileMeta file metadata
type FileMeta struct {
	Width  *int `json:"width,omitempty"`
	Height *int `json:"height,omitempty"`
}

// UploadFileByUrlRequest file url to upload
type UploadFileByUrlRequest struct {
	Url string `json:"url"`
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
	ChannelID     uint64    `json:"id"`
	DeactivatedAt time.Time `json:"deactivated_at"`
}

// ChannelListItem response struct
type ChannelListItem struct {
	ID            uint64          `json:"id"`
	Type          string          `json:"type"`
	Name          *string         `json:"name"`
	Settings      ChannelSettings `json:"settings"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     *string         `json:"updated_at"`
	ActivatedAt   string          `json:"activated_at"`
	DeactivatedAt *string         `json:"deactivated_at"`
	IsActive      bool            `json:"is_active"`
}

// Channels request type
type Channels struct {
	ID          int       `url:"id,omitempty" json:"id,omitempty"`
	Types       []string  `url:"types,omitempty" json:"types,omitempty"`
	Active      bool      `url:"active,omitempty" json:"active,omitempty"`
	Since       time.Time `url:"since,omitempty" json:"since,omitempty"`
	Until       time.Time `url:"until,omitempty" json:"until,omitempty"`
	TransportID uint64    `url:"transport_id,omitempty" json:"transport_id,omitempty"`
	Sort        string    `url:"sort,omitempty" json:"sort,omitempty"`
	Limit       int       `url:"limit,omitempty" json:"limit,omitempty"`
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
	Items      []Item `json:"items,omitempty"`
}

// SendMessage struct
type SendMessage struct {
	Message
	SentAt time.Time `json:"sent_at,omitempty"`
}

// EditMessageRequest type
type EditMessageRequest struct {
	Message EditMessageRequestMessage `json:"message"`
	Channel uint64                    `json:"channel"`
}

// EditMessageRequestMessage type
type EditMessageRequestMessage struct {
	ExternalID string `json:"external_id"`
	Text       string `json:"text"`
	EditedAt   int64  `json:"edited_at"`
}

// SendData struct
type SendData struct {
	Message        Message                  `json:"message"`
	Originator     string                   `json:"originator"`
	User           User                     `json:"user"`
	Channel        uint64                   `json:"channel"`
	ExternalChatID string                   `json:"external_chat_id"`
	Quote          *SendMessageRequestQuote `json:"quote,omitempty"`
}

// Item struct
type Item struct {
	ID      string `json:"id"`
	Caption string `json:"caption"`
}

// SendMessageRequestQuote type
type SendMessageRequestQuote struct {
	ExternalID string `json:"external_id"`
}

// MarkMessageReadResponse type
type MarkMessageReadResponse struct{}

// MarkMessageReadRequest type
type MarkMessageReadRequest struct {
	Message   MarkMessageReadRequestMessage `json:"message"`
	ChannelID uint64                        `json:"channel_id"`
}

// MarkMessageReadRequestMessage type
type MarkMessageReadRequestMessage struct {
	ExternalID string `json:"external_id"`
}

// DeleteData struct
type DeleteData struct {
	Message Message `json:"message"`
	Channel uint64  `json:"channel"`
}

// MessagesResponse message event response
type MessagesResponse struct {
	MessageID int       `json:"message_id,omitempty"`
	Time      time.Time `json:"time,omitempty"`
}

// WebhookRequest type
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
	Type              string              `json:"type"`
	User              *MessageDataUser    `json:"user,omitempty"`
	Bot               *MessageDataBot     `json:"bot,omitempty"`
	Product           *MessageDataProduct `json:"product,omitempty"`
	Order             *MessageDataOrder   `json:"order,omitempty"`
	Items             *[]FileItem         `json:"items,omitempty"`
}

// FileItem struct
type FileItem struct {
	ID      string `json:"id"`
	Size    int    `json:"size"`
	Caption string `json:"caption"`
	Height  *int   `json:"height,omitempty"`
	Width   *int   `json:"width,omitempty"`
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
	Number   string                    `json:"number"`
	Url      string                    `json:"url,omitempty"`
	Date     string                    `json:"date,omitempty"`
	Cost     *MessageDataOrderCost     `json:"cost,omitempty"`
	Status   *MessageDataOrderStatus   `json:"status,omitempty"`
	Delivery *MessageDataOrderDelivery `json:"delivery"`
	Payments []MessageDataOrderPayment `json:"payments"`
	Items    []MessageDataOrderItem    `json:"items,omitempty"`
}

// MessageDataOrderStatus type
type MessageDataOrderStatus struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

// MessageDataOrderItem type
type MessageDataOrderItem struct {
	Name     string                    `json:"name,omitempty"`
	Url      string                    `json:"url,omitempty"`
	Img      string                    `json:"img,omitempty"`
	Quantity *MessageDataOrderQuantity `json:"quantity,omitempty"`
	Price    *MessageDataOrderCost     `json:"price,omitempty"`
}

// MessageDataOrderCost type
type MessageDataOrderCost struct {
	Value    float32 `json:"value,omitempty"`
	Currency string  `json:"currency"`
}

// MessageDataOrderQuantity type
type MessageDataOrderQuantity struct {
	Value float32 `json:"value"`
	Unit  string  `json:"unit"`
}

// MessageDataOrderPayment type
type MessageDataOrderPayment struct {
	Name   string                         `json:"name"`
	Status *MessageDataOrderPaymentStatus `json:"status"`
	Amount *MessageDataOrderCost          `json:"amount"`
}

// MessageDataOrderPaymentStatus type
type MessageDataOrderPaymentStatus struct {
	Name  string `json:"name"`
	Payed bool   `json:"payed"`
}

// MessageDataOrderDelivery type
type MessageDataOrderDelivery struct {
	Name    string                `json:"name"`
	Price   *MessageDataOrderCost `json:"price"`
	Address string                `json:"address"`
	Comment string                `json:"comment,omitempty"`
}

// TransportRequestMeta request metadata
type TransportRequestMeta struct {
	ID        uint64 `json:"id"`
	Timestamp int64  `json:"timestamp"`
}
