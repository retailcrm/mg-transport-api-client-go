package v1

import "encoding/json"

type WebhookType string

const (
	MessageSendWebhookType    WebhookType = "message_sent"
	MessageUpdateWebhookType  WebhookType = "message_updated"
	MessageDeleteWebhookType  WebhookType = "message_deleted"
	MessageReadWebhookType    WebhookType = "message_read"
	ReactionAddWebhookType    WebhookType = "reaction_add"
	ReactionDeleteWebhookType WebhookType = "reaction_delete"
	TemplateCreateWebhookType WebhookType = "template_create"
	TemplateUpdateWebhookType WebhookType = "template_update"
	TemplateDeleteWebhookType WebhookType = "template_delete"
)

// WebhookRequest type.
type WebhookRequest struct {
	Type WebhookType          `json:"type"`
	Meta TransportRequestMeta `json:"meta"`
	Data json.RawMessage      `json:"data"`
}

// IsMessageWebhook returns true if current webhook contains data related to chat messages.
func (w WebhookRequest) IsMessageWebhook() bool {
	return w.Type == MessageReadWebhookType || w.Type == MessageDeleteWebhookType ||
		w.Type == MessageSendWebhookType || w.Type == MessageUpdateWebhookType
}

// IsReactionWebhook returns true if current webhook contains data related to chat reactions.
func (w WebhookRequest) IsReactionWebhook() bool {
	return w.Type == ReactionAddWebhookType || w.Type == ReactionDeleteWebhookType
}

// IsTemplateWebhook returns true if current webhook contains data related to the templates changes.
func (w WebhookRequest) IsTemplateWebhook() bool {
	return w.Type == TemplateCreateWebhookType ||
		w.Type == TemplateUpdateWebhookType ||
		w.Type == TemplateDeleteWebhookType
}

// MessageWebhookData returns the message data from webhook contents.
//
// Note: this call will not fail even if underlying data is not related to the messages.
// Use IsMessageWebhook to mitigate this.
func (w WebhookRequest) MessageWebhookData() (wd MessageWebhookData) {
	_ = json.Unmarshal(w.Data, &wd)
	return
}

// ReactionWebhookData returns the reaction data from webhook contents.
//
// Note: this call will not fail even if underlying data is not related to the reactions.
// Use IsReactionWebhook to mitigate this.
func (w WebhookRequest) ReactionWebhookData() (wd ReactionWebhookData) {
	_ = json.Unmarshal(w.Data, &wd)
	return
}

// TemplateCreateWebhookData returns new template data from webhook contents.
// This method is used if current webhook was initiated because user created a template.
//
// Note: this call will not fail even if underlying data is not related to the templates.
// Use IsTemplateWebhook or direct Type comparison (Type == TemplateCreateWebhookType) to mitigate this.
func (w WebhookRequest) TemplateCreateWebhookData() (wd TemplateCreateWebhookData) {
	_ = json.Unmarshal(w.Data, &wd)
	return
}

// TemplateUpdateWebhookData returns existing template data from webhook contents.
// This method is used if current webhook was initiated because user updated a template.
//
// Note: this call will not fail even if underlying data is not related to the templates.
// Use IsTemplateWebhook or direct Type comparison (Type == TemplateUpdateWebhookData) to mitigate this.
func (w WebhookRequest) TemplateUpdateWebhookData() (wd TemplateUpdateWebhookData) {
	_ = json.Unmarshal(w.Data, &wd)
	return
}

// TemplateDeleteWebhookData returns existing template data from webhook contents.
// This method is used if current webhook was initiated because user deleted a template.
//
// Note: this call will not fail even if underlying data is not related to the templates.
// Use IsTemplateWebhook or direct Type comparison (Type == TemplateDeleteWebhookType) to mitigate this.
func (w WebhookRequest) TemplateDeleteWebhookData() (wd TemplateDeleteWebhookData) {
	_ = json.Unmarshal(w.Data, &wd)
	return
}
