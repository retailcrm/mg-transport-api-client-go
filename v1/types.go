package v1

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// noinspection ALL.
const (
	// ChannelFeatureNone channel can not implement feature.
	ChannelFeatureNone string = "none"
	// ChannelFeatureReceive channel implement feature on receive.
	ChannelFeatureReceive string = "receive"
	// ChannelFeatureSend channel implement feature on send.
	ChannelFeatureSend string = "send"
	// ChannelFeatureBoth channel implement feature on both directions.
	ChannelFeatureBoth string = "both"
	// ChannelFeatureAny channel implement feature on any.
	ChannelFeatureAny string = "any"
	// ChannelFeatureSendingPolicyNo channel can not implement feature.
	ChannelFeatureSendingPolicyNo string = "no"
	// ChannelFeatureSendingPolicyTemplate channel can implement template.
	ChannelFeatureSendingPolicyTemplate string = "template"
	// ChannelFeatureCustomerExternalIDPhone customer externalId is phone.
	ChannelFeatureCustomerExternalIDPhone string = "phone"

	// MsgTypeText text message.
	MsgTypeText string = "text"
	// MsgTypeSystem system message.
	MsgTypeSystem string = "system"
	// MsgTypeCommand command (for bots).
	MsgTypeCommand string = "command"
	// MsgTypeOrder order card.
	MsgTypeOrder string = "order"
	// MsgTypeProduct product card.
	MsgTypeProduct string = "product"
	// MsgTypeFile file card.
	MsgTypeFile string = "file"
	// MsgTypeImage image card.
	MsgTypeImage string = "image"
	// MsgTypeAudio audio.
	MsgTypeAudio string = "audio"

	// MsgOrderStatusCodeNew order status group new.
	MsgOrderStatusCodeNew = "new"
	// MsgOrderStatusCodeApproval order status group approval.
	MsgOrderStatusCodeApproval = "approval"
	// MsgOrderStatusCodeAssembling order status group assembling.
	MsgOrderStatusCodeAssembling = "assembling"
	// MsgOrderStatusCodeDelivery order status group delivery.
	MsgOrderStatusCodeDelivery = "delivery"
	// MsgOrderStatusCodeComplete order status group complete.
	MsgOrderStatusCodeComplete = "complete"
	// MsgOrderStatusCodeCancel order status group cancel.
	MsgOrderStatusCodeCancel = "cancel"

	FileSizeLimit = 20 * 1024 * 1024
)

const (
	// OriginatorCustomer means message was created by customer.
	OriginatorCustomer Originator = iota + 1
	// OriginatorChannel means message was created by channel, for example via messenger mobile application.
	OriginatorChannel
)

type ErrorType string

const (
	GeneralError           ErrorType = "general"
	CustomerNotExistsError ErrorType = "customer_not_exists"
	ReplyTimedOutError     ErrorType = "reply_timed_out"
	SpamSuspicionError     ErrorType = "spam_suspicion"
	AccessRestrictedError  ErrorType = "access_restricted"
)

// MgClient type.
type MgClient struct {
	URL        string       `json:"url"`
	Token      string       `json:"token"`
	Debug      bool         `json:"debug"`
	httpClient *http.Client `json:"-"`
	logger     BasicLogger  `json:"-"`
	limiter    Limiter      `json:"-"`
}

// Channel type.
type Channel struct {
	ID         uint64          `json:"id,omitempty"`
	ExternalID string          `json:"external_id,omitempty"`
	Type       string          `json:"type,omitempty"`
	Name       string          `json:"name,omitempty"`
	AvatarUrl  string          `json:"avatar_url,omitempty"`
	Settings   ChannelSettings `json:"settings,omitempty"`
}

// ChannelSettings struct.
type ChannelSettings struct {
	Status             Status                     `json:"status"`
	Text               ChannelSettingsText        `json:"text"`
	Product            Product                    `json:"product"`
	Order              Order                      `json:"order"`
	File               ChannelSettingsFilesBase   `json:"file"`
	Image              ChannelSettingsFilesBase   `json:"image"`
	Reaction           Reaction                   `json:"reaction"`
	CustomerExternalID string                     `json:"customer_external_id,omitempty"`
	SendingPolicy      SendingPolicy              `json:"sending_policy,omitempty"`
	Suggestions        ChannelSettingsSuggestions `json:"suggestions,omitempty"`
	Audio              ChannelSettingsAudio       `json:"audio"`
	Template           ChannelSettingsTemplate    `json:"template"`
	WhatsApp           *WhatsAppChannelProperties `json:"whatsapp,omitempty"` // for WhatsApp channels only.
}

// Product type.
type Product struct {
	Creating string `json:"creating,omitempty"`
	Editing  string `json:"editing,omitempty"`
	Deleting string `json:"deleting,omitempty"`
	Reaction string `json:"reaction,omitempty"`
}

type Reaction struct {
	Dictionary []string `json:"dictionary,omitempty"`
	MaxCount   uint16   `json:"max_count,omitempty"`
}

// Order type.
type Order struct {
	Creating string `json:"creating,omitempty"`
	Editing  string `json:"editing,omitempty"`
	Deleting string `json:"deleting,omitempty"`
	Reaction string `json:"reaction,omitempty"`
}

// Status struct.
type Status struct {
	Delivered string `json:"delivered,omitempty"`
	Read      string `json:"read,omitempty"`
}

// ChannelSettingsText struct.
type ChannelSettingsText struct {
	Creating      string `json:"creating,omitempty"`
	Editing       string `json:"editing,omitempty"`
	Quoting       string `json:"quoting,omitempty"`
	Deleting      string `json:"deleting,omitempty"`
	Reaction      string `json:"reaction,omitempty"`
	MaxCharsCount uint16 `json:"max_chars_count,omitempty"`
}

// ChannelSettingsFilesBase struct.
type ChannelSettingsFilesBase struct {
	Creating          string  `json:"creating,omitempty"`
	Editing           string  `json:"editing,omitempty"`
	Quoting           string  `json:"quoting,omitempty"`
	Deleting          string  `json:"deleting,omitempty"`
	Reaction          string  `json:"reaction,omitempty"`
	Max               uint64  `json:"max_items_count,omitempty"`
	NoteMaxCharsCount *uint16 `json:"note_max_chars_count,omitempty"`
	MaxItemSize       *uint64 `json:"max_item_size,omitempty"`
}

// ChannelSettingsAudio struct.
type ChannelSettingsAudio struct {
	Creating      string  `json:"creating,omitempty"`
	Quoting       string  `json:"quoting,omitempty"`
	Deleting      string  `json:"deleting,omitempty"`
	Reaction      string  `json:"reaction,omitempty"`
	MaxItemsCount uint64  `json:"max_items_count,omitempty"`
	MaxItemSize   *uint64 `json:"max_item_size,omitempty"`
}

type SendingPolicy struct {
	NewCustomer       string `json:"new_customer,omitempty"`
	AfterReplyTimeout string `json:"after_reply_timeout,omitempty"`
	Outgoing          string `json:"outgoing,omitempty"`
}

type ChannelSettingsSuggestions struct {
	Text  string `json:"text,omitempty"`
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}

type ChannelSettingsTemplate struct {
	Creation bool `json:"creation,omitempty"`
}

// WhatsAppChannelProperties WhatsApp-specific channel characteristics
type WhatsAppChannelProperties struct {
	Tier           *int                    `json:"tier"`
	ChannelQuality *WhatsAppChannelQuality `json:"channel_quality"`
	ChannelStatus  *WhatsAppChannelStatus  `json:"channel_status"`
}

type WhatsAppChannelQuality string

var (
	WhatsAppChannelQualityHigh   WhatsAppChannelQuality = "high"
	WhatsAppChannelQualityMedium WhatsAppChannelQuality = "medium"
	WhatsAppChannelQualityLow    WhatsAppChannelQuality = "low"
)

type WhatsAppChannelStatus string

var (
	WhatsAppChannelStatusConnected  WhatsAppChannelStatus = "connected"
	WhatsAppChannelStatusFlagged    WhatsAppChannelStatus = "flagged"
	WhatsAppChannelStatusOffline    WhatsAppChannelStatus = "offline"
	WhatsAppChannelStatusPending    WhatsAppChannelStatus = "pending"
	WhatsAppChannelStatusRestricted WhatsAppChannelStatus = "restricted"
)

// FullFileResponse uploaded file data.
type FullFileResponse struct {
	ID       string `json:"id,omitempty"`
	Type     string `json:"type,omitempty"`
	Size     int    `json:"size,omitempty"`
	Url      string `json:"url,omitempty"` //nolint:golint
	MimeType string `json:"mime_type,omitempty"`
}

// UploadFileResponse uploaded file data.
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

// FileMeta file metadata.
type FileMeta struct {
	Width  *int `json:"width,omitempty"`
	Height *int `json:"height,omitempty"`
}

// UploadFileByUrlRequest file url to upload.
type UploadFileByUrlRequest struct {
	Url string `json:"url"`
}

// ActivateResponse channel activation response.
type ActivateResponse struct {
	ChannelID   uint64    `json:"id"`
	ExternalID  string    `json:"external_id"`
	ActivatedAt time.Time `json:"activated_at"`
}

// UpdateResponse channel update response.
type UpdateResponse struct {
	ChannelID  uint64    `json:"id"`
	ExternalID string    `json:"external_id"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// DeleteResponse channel deactivation response.
type DeleteResponse struct {
	ChannelID     uint64    `json:"id"`
	DeactivatedAt time.Time `json:"deactivated_at"`
}

// ChannelListItem response struct.
type ChannelListItem struct {
	ID            uint64          `json:"id"`
	ExternalID    string          `json:"external_id"`
	Type          string          `json:"type"`
	Name          *string         `json:"name"`
	Settings      ChannelSettings `json:"settings"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     *string         `json:"updated_at"`
	ActivatedAt   string          `json:"activated_at"`
	DeactivatedAt *string         `json:"deactivated_at"`
	IsActive      bool            `json:"is_active"`
}

// Channels request type.
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

// Customer struct.
type Customer struct {
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
	Utm        *Utm   `json:"utm,omitempty"`
}

type Utm struct {
	Source   string `json:"source,omitempty"`
	Medium   string `json:"medium,omitempty"`
	Campaign string `json:"campaign,omitempty"`
	Term     string `json:"term,omitempty"`
	Content  string `json:"content,omitempty"`
}

// Message struct.
type Message struct {
	ExternalID string              `json:"external_id"`
	Type       string              `json:"type,omitempty"`
	Text       string              `json:"text,omitempty"`
	Note       string              `json:"note,omitempty"`
	Items      []Item              `json:"items,omitempty"`
	PageLink   string              `json:"page_link,omitempty"`
	Product    *MessageDataProduct `json:"product,omitempty"`
	Order      *MessageDataOrder   `json:"order,omitempty"`
}

// SendMessage struct.
type SendMessage struct {
	Message
	SentAt time.Time `json:"sent_at,omitempty"`
}

// EditMessageRequest type.
type EditMessageRequest struct {
	Message EditMessageRequestMessage `json:"message"`
	Channel uint64                    `json:"channel"`
}

// EditMessageRequestMessage type.
type EditMessageRequestMessage struct {
	ExternalID string `json:"external_id"`
	Text       string `json:"text"`
	EditedAt   int64  `json:"edited_at"`
	PageLink   string `json:"page_link,omitempty"`
}

type SendHistoryMessageRequest struct {
	Message        SendMessageRequestMessage `json:"message"`
	ChannelID      uint64                    `json:"channel_id"`
	ExternalChatID string                    `json:"external_chat_id"`
	Customer       *Customer                 `json:"customer"`
	Quote          *SendMessageRequestQuote  `json:"quote,omitempty"`
	Originator     Originator                `json:"originator,omitempty"`
	ReplyDeadline  *time.Time                `json:"reply_deadline,omitempty"`
}

type ReactionRequest struct {
	ChannelID uint64                   `json:"channel_id"`
	Message   ReactionMessageReference `json:"message"`
	Reaction  string                   `json:"reaction,omitempty"`
}

type ReactionMessageReference struct {
	ExternalID string `json:"external_id"`
}

type SendMessageRequestMessage struct {
	Type       string     `json:"type"`
	ExternalID string     `json:"external_id,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	IsComment  bool       `json:"is_comment,omitempty"`
	Text       string     `json:"text"`
	Items      []Item     `json:"items"`
	Note       string     `json:"note,omitempty"`
}

// SendData struct.
type SendData struct {
	Message                  Message                  `json:"message"`
	Originator               Originator               `json:"originator,omitempty"`
	Customer                 Customer                 `json:"customer"`
	Channel                  uint64                   `json:"channel"`
	ExternalChatID           string                   `json:"external_chat_id"`
	Quote                    *SendMessageRequestQuote `json:"quote,omitempty"`
	ReplyDeadline            *time.Time               `json:"reply_deadline,omitempty"`
	SecondaryExternalChatIDs []string                 `json:"secondary_external_chat_ids,omitempty"`
}

// Item struct.
type Item struct {
	ID      string `json:"id"`
	Caption string `json:"caption"`
}

// SendMessageRequestQuote type.
type SendMessageRequestQuote struct {
	ExternalID string `json:"external_id"`
}

// MarkMessageReadResponse type.
type MarkMessageReadResponse struct{}

type MessageReactionResponse struct{}

// MarkMessageReadRequest type.
type MarkMessageReadRequest struct {
	Message   MarkMessageReadRequestMessage `json:"message"`
	ChannelID uint64                        `json:"channel_id"`
}

// MarkMessageReadRequestMessage type.
type MarkMessageReadRequestMessage struct {
	ExternalID string `json:"external_id"`
}

// AckMessageRequest type.
type AckMessageRequest struct {
	ExternalMessageID  string            `json:"external_message_id"`
	TransportMessageID string            `json:"transport_message_id,omitempty"`
	Channel            uint64            `json:"channel"`
	Error              *MessageSentError `json:"error,omitempty"`
}

// MarkMessagesReadUntilRequest type.
type MarkMessagesReadUntilRequest struct {
	CustomerExternalID string    `json:"customer_external_id"`
	ChannelID          uint64    `json:"channel_id"`
	Until              time.Time `json:"until"`
}

// MarkMessagesReadUntilResponse type.
type MarkMessagesReadUntilResponse struct {
	IDs []int64 `json:"ids"`
}

// DeleteData struct.
type DeleteData struct {
	Message Message `json:"message"`
	Channel uint64  `json:"channel"`
}

// MessagesResponse message event response.
type MessagesResponse struct {
	MessageID int       `json:"message_id,omitempty"`
	Time      time.Time `json:"time,omitempty"`
	Warnings  []string  `json:"warnings"`
}

// WebhookMessageSentResponse type
// Consider using this structure while processing webhook request.
type WebhookMessageSentResponse struct {
	ExternalMessageID  string            `json:"external_message_id"`
	Error              *MessageSentError `json:"error,omitempty"`
	Async              bool              `json:"async"`
	ExternalCustomerID string            `json:"external_customer_id,omitempty"`
	ExternalChatID     string            `json:"external_chat_id,omitempty"`
}

// MessageSentError type.
type MessageSentError struct {
	Code    ErrorType `json:"code"`
	Message string    `json:"message"`
}

// MessageWebhookData request data.
type MessageWebhookData struct {
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
	Template          *TemplateInfo       `json:"template,omitempty"`
	Attachments       *Attachments        `json:"attachments,omitempty"`
	InAppID           int32               `json:"in_app_id,omitempty"`
}

type ReactionWebhookData struct {
	ExternalUserID    string         `json:"external_user_id"`
	ExternalChatID    string         `json:"external_chat_id"`
	ChannelID         uint64         `json:"channel_id"`
	ExternalMessageID string         `json:"external_message_id"`
	NewReaction       string         `json:"new_reaction,omitempty"`
	OldReaction       string         `json:"old_reaction,omitempty"`
	AllReactions      []ReactionInfo `json:"all_reactions,omitempty"`
}

type ReactionInfo struct {
	Reaction string `json:"reaction"`
}

type Attachments struct {
	Suggestions []Suggestion `json:"suggestions,omitempty"`
}

const (
	SuggestionTypeText  SuggestionType = "text"
	SuggestionTypeEmail SuggestionType = "email"
	SuggestionTypePhone SuggestionType = "phone"
	SuggestionTypeURL   SuggestionType = "url"
)

type SuggestionType string

type Suggestion struct {
	Type    SuggestionType `json:"type"`
	Title   string         `json:"title,omitempty"`   // required for type=text and optional for others
	Payload string         `json:"payload,omitempty"` // ignored for type=text
}

type TemplateInfo struct {
	Code      string            `json:"code"`
	Variables TemplateArguments `json:"variables,omitempty"`
}

type TemplateArguments struct {
	Header  *TemplateHeaderArguments  `json:"header,omitempty"`
	Body    TemplateBodyArguments     `json:"body"`
	Buttons []TemplateButtonArguments `json:"buttons,omitempty"`
}

type TemplateHeaderArguments struct {
	Type        string               `json:"type"`
	Args        []string             `json:"args,omitempty"`        // for type="text" only
	Attachments []TemplateAttachment `json:"attachments,omitempty"` // for media templates only
}

type TemplateBodyArguments struct {
	Args []string `json:"args"`
}

type TemplateButtonArguments struct {
	Type  string   `json:"type"`
	Title string   `json:"title"`
	Args  []string `json:"args,omitempty"`
}

type TemplateAttachment struct {
	ID      string `json:"id"`
	Caption string `json:"caption"`
}

// FileItem struct.
type FileItem struct {
	ID      string `json:"id"`
	Size    int    `json:"size"`
	Caption string `json:"caption"`
	Height  *int   `json:"height,omitempty"`
	Width   *int   `json:"width,omitempty"`
}

// MessageDataUser user data from webhook.
type MessageDataUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

// MessageDataBot bot data from webhook.
type MessageDataBot struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// MessageDataProduct product data from webhook.
type MessageDataProduct struct {
	ID      uint64                `json:"id"`
	Name    string                `json:"name"`
	Article string                `json:"article,omitempty"`
	Url     string                `json:"url,omitempty"`
	Img     string                `json:"img,omitempty"`
	Cost    *MessageDataOrderCost `json:"cost,omitempty"`
	Unit    string                `json:"unit,omitempty"`
	// Deprecated: now you need to use Unit instead of this field
	Quantity *MessageDataOrderQuantity `json:"quantity,omitempty"`
}

// MessageDataOrder order data from webhook.
type MessageDataOrder struct {
	ExternalID int64                     `json:"external_id"`
	Number     string                    `json:"number"`
	URL        string                    `json:"url,omitempty"`
	Date       string                    `json:"date,omitempty"`
	Cost       *MessageDataOrderCost     `json:"cost,omitempty"`
	Discount   *MessageDataOrderCost     `json:"discount,omitempty"`
	Status     *MessageDataOrderStatus   `json:"status,omitempty"`
	Delivery   *MessageDataOrderDelivery `json:"delivery"`
	Payments   []MessageDataOrderPayment `json:"payments"`
	Items      []MessageDataOrderItem    `json:"items,omitempty"`
}

// MessageDataOrderStatus type.
type MessageDataOrderStatus struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

// MessageDataOrderItem type.
type MessageDataOrderItem struct {
	Name     string                    `json:"name,omitempty"`
	Url      string                    `json:"url,omitempty"`
	Img      string                    `json:"img,omitempty"`
	Quantity *MessageDataOrderQuantity `json:"quantity,omitempty"`
	Price    *MessageDataOrderCost     `json:"price,omitempty"`
}

// MessageDataOrderCost type.
type MessageDataOrderCost struct {
	Value    float32 `json:"value,omitempty"`
	Currency string  `json:"currency"`
}

// MessageDataOrderQuantity type.
type MessageDataOrderQuantity struct {
	Value float32 `json:"value"`
	Unit  string  `json:"unit"`
}

// MessageDataOrderPayment type.
type MessageDataOrderPayment struct {
	Name   string                         `json:"name"`
	Status *MessageDataOrderPaymentStatus `json:"status"`
	Amount *MessageDataOrderCost          `json:"amount"`
}

// MessageDataOrderPaymentStatus type.
type MessageDataOrderPaymentStatus struct {
	Name string `json:"name"`
	Paid bool   `json:"paid"`
}

// MessageDataOrderDelivery type.
type MessageDataOrderDelivery struct {
	Name    string                `json:"name"`
	Price   *MessageDataOrderCost `json:"price"`
	Address string                `json:"address"`
	Comment string                `json:"comment,omitempty"`
}

// TransportRequestMeta request metadata.
type TransportRequestMeta struct {
	ID        uint64 `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

type UpdateTemplateRequest struct {
	Name               string                     `json:"name"`
	Template           []TemplateItem             `json:"template,omitempty"`
	Body               string                     `json:"body"`
	Lang               string                     `json:"lang,omitempty"`
	Category           string                     `json:"category,omitempty"`
	Example            *TemplateExample           `json:"example,omitempty"`
	VerificationStatus TemplateVerificationStatus `json:"verification_status"`
	Quality            *TemplateQuality           `json:"quality,omitempty"`
	RejectionReason    TemplateRejectionReason    `json:"rejection_reason,omitempty"`
	Header             *TemplateHeader            `json:"header,omitempty"`
	Footer             string                     `json:"footer,omitempty"`
	Buttons            *TemplateButtons           `json:"buttons,omitempty"`
}

type ActivateTemplateRequest struct {
	UpdateTemplateRequest

	Code string       `json:"code"`
	Type TemplateType `json:"type"`
}

var ErrInvalidOriginator = errors.New("invalid originator")

// Originator of message.
type Originator byte

// MarshalText marshals originator to text.
func (o Originator) MarshalText() ([]byte, error) {
	switch o {
	case OriginatorCustomer:
		return []byte("customer"), nil
	case OriginatorChannel:
		return []byte("channel"), nil
	}

	return nil, ErrInvalidOriginator
}

// UnmarshalText unmarshals originator from text to the value.
func (o *Originator) UnmarshalText(text []byte) error {
	switch string(text) {
	case "customer":
		*o = OriginatorCustomer
		return nil
	case "channel":
		*o = OriginatorChannel
		return nil
	}

	return ErrInvalidOriginator
}

type TransportErrorCode string

const (
	MessageErrorGeneral           TransportErrorCode = "general"
	MessageErrorCustomerNotExists TransportErrorCode = "customer_not_exists"
	MessageErrorReplyTimedOut     TransportErrorCode = "reply_timed_out"
	MessageErrorSpamSuspicion     TransportErrorCode = "spam_suspicion"
	MessageErrorAccessRestricted  TransportErrorCode = "access_restricted"
)

type TransportResponse struct {
	ExternalMessageID string          `json:"external_message_id,omitempty"`
	Error             *TransportError `json:"error,omitempty"`
}

type TransportError struct {
	Code    TransportErrorCode `json:"code"`
	Message string             `json:"message,omitempty"`
}

func (t TransportErrorCode) MarshalJSON() ([]byte, error) {
	if t == "" {
		return []byte(fmt.Sprintf(`"%s"`, MessageErrorGeneral)), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t)), nil
}

func NewSentMessageResponse(externalMessageID string) TransportResponse {
	return TransportResponse{ExternalMessageID: externalMessageID}
}

func NewTransportErrorResponse(code TransportErrorCode, message string) TransportResponse {
	return TransportResponse{
		Error: &TransportError{
			Code:    code,
			Message: message,
		},
	}
}

type HeaderParams struct {
	TextVars    []string `json:"textVars,omitempty"`
	ImageURL    string   `json:"imageUrl,omitempty"`
	VideoURL    string   `json:"videoUrl,omitempty"`
	DocumentURL string   `json:"documentUrl,omitempty"`
}

type ButtonParam struct {
	ButtonType   ButtonType `json:"type"`
	Text         string     `json:"text,omitempty"`
	URLParameter string     `json:"urlParameter,omitempty"`
}

type TemplateContent struct {
	Name     string           `json:"name"`
	Lang     string           `json:"lang"`
	Category string           `json:"category"`
	Body     string           `json:"body"`
	Header   *TemplateHeader  `json:"header,omitempty"`
	Footer   string           `json:"footer,omitempty"`
	Buttons  *TemplateButtons `json:"buttons,omitempty"`
	Example  *TemplateExample `json:"example,omitempty"`
}

type TemplateCreateWebhookData struct {
	TemplateContent
	ChannelID int64 `json:"channel_id"`
}

type TemplateCreateWebhookResponse struct {
	Code               string                     `json:"code" binding:"required"`
	VerificationStatus TemplateVerificationStatus `json:"verification_status" binding:"required"`
}

type TemplateUpdateWebhookData struct {
	TemplateContent
	ChannelID int64  `json:"channel_id"`
	Code      string `json:"code"`
}

type TemplateDeleteWebhookData struct {
	ChannelID int64  `json:"channel_id"`
	Code      string `json:"code"`
	Lang      string `json:"lang"`
}

type TemplateVerificationStatus string

const (
	TemplateStatusApproved TemplateVerificationStatus = "approved"
	TemplateStatusPending  TemplateVerificationStatus = "pending"
	TemplateStatusRejected TemplateVerificationStatus = "rejected"
	TemplateStatusNew      TemplateVerificationStatus = "new"
	TemplateStatusPaused   TemplateVerificationStatus = "paused"
	TemplateStatusDisabled TemplateVerificationStatus = "disabled"
)

type TemplateQuality string

const (
	TemplateQualityPending TemplateQuality = "pending"
	TemplateQualityHigh    TemplateQuality = "high"
	TemplateQualityMedium  TemplateQuality = "medium"
	TemplateQualityLow     TemplateQuality = "low"
)

type TemplateRejectionReason string

const (
	ReasonAbusiveContent    TemplateRejectionReason = "abusive_content"
	ReasonIncorrectCategory TemplateRejectionReason = "incorrect_category"
	ReasonInvalidFormat     TemplateRejectionReason = "invalid_format"
	ReasonScam              TemplateRejectionReason = "scam"
)
