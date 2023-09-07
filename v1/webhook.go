package v1

import "encoding/json"

type WebhookType string

const (
	MessageSendWebhookType    WebhookType = "message_sent"
	MessageUpdateWebhookType  WebhookType = "message_updated"
	MessageDeleteWebhookType  WebhookType = "message_deleted"
	MessageReadWebhookType    WebhookType = "message_read"
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

func (w WebhookRequest) IsTemplateWebhook() bool {
	return w.Type == TemplateCreateWebhookType ||
		w.Type == TemplateUpdateWebhookType ||
		w.Type == TemplateDeleteWebhookType
}

func (w WebhookRequest) MessageWebhookData() (wd MessageWebhookData) {
	_ = json.Unmarshal(w.Data, &wd)
	return
}

func (w WebhookRequest) TemplateCreateWebhookData() (wd TemplateCreateWebhookData) {
	_ = json.Unmarshal(w.Data, &wd)
	return
}

func (w WebhookRequest) TemplateUpdateWebhookData() (wd TemplateUpdateWebhookData) {
	_ = json.Unmarshal(w.Data, &wd)
	return
}

func (w WebhookRequest) TemplateDeleteWebhookData() (wd TemplateDeleteWebhookData) {
	_ = json.Unmarshal(w.Data, &wd)
	return
}
