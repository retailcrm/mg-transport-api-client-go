package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookRequest_IsMessageWebhook(t *testing.T) {
	assert.True(t, WebhookRequest{Type: MessageSendWebhookType}.IsMessageWebhook())
	assert.True(t, WebhookRequest{Type: MessageUpdateWebhookType}.IsMessageWebhook())
	assert.True(t, WebhookRequest{Type: MessageDeleteWebhookType}.IsMessageWebhook())
	assert.True(t, WebhookRequest{Type: MessageReadWebhookType}.IsMessageWebhook())
	assert.False(t, WebhookRequest{}.IsMessageWebhook())
}

func TestWebhookRequest_IsTemplateWebhook(t *testing.T) {
	assert.True(t, WebhookRequest{Type: TemplateCreateWebhookType}.IsTemplateWebhook())
	assert.True(t, WebhookRequest{Type: TemplateUpdateWebhookType}.IsTemplateWebhook())
	assert.True(t, WebhookRequest{Type: TemplateDeleteWebhookType}.IsTemplateWebhook())
	assert.False(t, WebhookRequest{}.IsTemplateWebhook())
}

func TestWebhookData_MessageWebhookData(t *testing.T) {
	wh := WebhookRequest{
		Type: MessageSendWebhookType,
		Data: mustMarshalJSON(MessageWebhookData{
			ExternalUserID:    "1",
			ExternalMessageID: "1",
			ExternalChatID:    "1",
			ChannelID:         1,
			Content:           "test",
			Type:              MsgTypeText,
		}),
	}.MessageWebhookData()
	assert.Equal(t, "test", wh.Content)
}

func TestWebhookData_TemplateCreateWebhookData(t *testing.T) {
	wh := WebhookRequest{
		Type: TemplateCreateWebhookType,
		Data: mustMarshalJSON(TemplateCreateWebhookData{
			TemplateContent: TemplateContent{
				Name: "template",
			},
			ChannelID: 1,
		}),
	}.TemplateCreateWebhookData()
	assert.Equal(t, "template", wh.TemplateContent.Name)
}

func TestWebhookData_TemplateEditWebhookData(t *testing.T) {
	wh := WebhookRequest{
		Type: TemplateUpdateWebhookType,
		Data: mustMarshalJSON(TemplateUpdateWebhookData{
			TemplateContent: TemplateContent{
				Name: "template",
			},
			ChannelID: 1,
			Code:      "code",
		}),
	}.TemplateUpdateWebhookData()
	assert.Equal(t, "template", wh.TemplateContent.Name)
	assert.Equal(t, "code", wh.Code)
}

func TestWebhookData_TemplateDeleteWebhookData(t *testing.T) {
	wh := WebhookRequest{
		Type: TemplateDeleteWebhookType,
		Data: mustMarshalJSON(TemplateDeleteWebhookData{
			ChannelID: 1,
			Code:      "code",
		}),
	}.TemplateDeleteWebhookData()
	assert.Equal(t, "code", wh.Code)
}

func mustMarshalJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
