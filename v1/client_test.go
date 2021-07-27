package v1

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	mgURL               = os.Getenv("MG_URL")
	mgToken             = os.Getenv("MG_TOKEN")
	channelID, _        = strconv.ParseUint(os.Getenv("MG_CHANNEL"), 10, 64)
	ext                 = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	tplCode             = fmt.Sprintf("testTemplate_%d", time.Now().UnixNano())
	tplChannel   uint64 = 0
)

func client() *MgClient {
	c := New(mgURL, mgToken)
	c.Debug = true

	return c
}

func templateChannel(t *testing.T) uint64 {
	if tplChannel == 0 {
		c := client()
		resp, _, err := c.ActivateTransportChannel(Channel{
			Type: "telegram",
			Name: "@test_channel_templates",
			Settings: ChannelSettings{
				SpamAllowed: false,
				Status: Status{
					Delivered: ChannelFeatureBoth,
					Read:      ChannelFeatureBoth,
				},
				Text: ChannelSettingsText{
					Creating:      ChannelFeatureBoth,
					Editing:       ChannelFeatureBoth,
					Quoting:       ChannelFeatureBoth,
					Deleting:      ChannelFeatureBoth,
					MaxCharsCount: 5000,
				},
				Product: Product{
					Creating: ChannelFeatureBoth,
					Editing:  ChannelFeatureBoth,
					Deleting: ChannelFeatureBoth,
				},
				Order: Order{
					Creating: ChannelFeatureBoth,
					Editing:  ChannelFeatureBoth,
					Deleting: ChannelFeatureBoth,
				},
				File: ChannelSettingsFilesBase{
					Creating:             ChannelFeatureBoth,
					Editing:              ChannelFeatureBoth,
					Quoting:              ChannelFeatureBoth,
					Deleting:             ChannelFeatureBoth,
					Max:                  1000000,
					CommentMaxCharsCount: 128,
				},
				Image: ChannelSettingsFilesBase{
					Creating: ChannelFeatureBoth,
					Editing:  ChannelFeatureBoth,
					Quoting:  ChannelFeatureBoth,
					Deleting: ChannelFeatureBoth,
				},
				CustomerExternalID: ChannelFeatureCustomerExternalIDPhone,
				SendingPolicy: SendingPolicy{
					NewCustomer: ChannelFeatureSendingPolicyTemplate,
				},
			},
		})

		if err != nil {
			t.FailNow()
		}

		tplChannel = resp.ChannelID
	}

	return tplChannel
}

func TestMgClient_TransportChannels(t *testing.T) {
	c := client()

	data, status, err := c.TransportChannels(Channels{Active: true})

	if err != nil {
		t.Errorf("%d %v", status, err)
	}

	t.Logf("Channels found: %v", len(data))
}

func TestMgClient_ActivateTransportChannel(t *testing.T) {
	c := client()
	ch := Channel{
		ID:   channelID,
		Type: "telegram",
		Name: "@my_shopping_bot",
		Settings: ChannelSettings{
			SpamAllowed: false,
			Status: Status{
				Delivered: ChannelFeatureNone,
				Read:      ChannelFeatureReceive,
			},
			Text: ChannelSettingsText{
				Creating:      ChannelFeatureBoth,
				Editing:       ChannelFeatureSend,
				Quoting:       ChannelFeatureReceive,
				Deleting:      ChannelFeatureBoth,
				MaxCharsCount: 2000,
			},
			Product: Product{
				Creating: ChannelFeatureSend,
				Deleting: ChannelFeatureSend,
			},
			Order: Order{
				Creating: ChannelFeatureBoth,
				Deleting: ChannelFeatureSend,
			},
			Image: ChannelSettingsFilesBase{
				Creating: ChannelFeatureBoth,
			},
			File: ChannelSettingsFilesBase{
				Creating: ChannelFeatureBoth,
			},
		},
	}

	data, status, err := c.ActivateTransportChannel(ch)

	if err != nil {
		t.Errorf("%d %v", status, err)
	}

	t.Logf("Activate selected channel: %v", data.ChannelID)
}

func TestMgClient_ActivateNewTransportChannel(t *testing.T) {
	c := client()
	ch := Channel{
		Type: "telegram",
		Name: "@my_shopping_bot",
		Settings: ChannelSettings{
			SpamAllowed: false,
			Status: Status{
				Delivered: ChannelFeatureNone,
				Read:      ChannelFeatureBoth,
			},
			Text: ChannelSettingsText{
				Creating: ChannelFeatureBoth,
				Editing:  ChannelFeatureSend,
				Quoting:  ChannelFeatureBoth,
				Deleting: ChannelFeatureSend,
			},
			Product: Product{
				Creating: ChannelFeatureSend,
				Deleting: ChannelFeatureSend,
			},
			Order: Order{
				Creating: ChannelFeatureBoth,
				Deleting: ChannelFeatureSend,
			},
			Image: ChannelSettingsFilesBase{
				Creating: ChannelFeatureBoth,
			},
			File: ChannelSettingsFilesBase{
				Creating: ChannelFeatureBoth,
			},
		},
	}

	data, status, err := c.ActivateTransportChannel(ch)

	if err != nil {
		t.Errorf("%d %v", status, err)
	}

	t.Logf("New channel ID %v", data.ChannelID)

	deleteData, status, err := c.DeactivateTransportChannel(data.ChannelID)

	if err != nil {
		t.Errorf("%d %v", status, err)
	}

	if deleteData.DeactivatedAt.String() == "" {
		t.Errorf("%v", err)
	}

	t.Logf("Deactivate new channel with ID %v", deleteData.ChannelID)
}

func TestMgClient_UpdateTransportChannel(t *testing.T) {
	c := client()
	ch := Channel{
		ID:   channelID,
		Name: "@my_shopping_bot_2",
		Settings: ChannelSettings{
			SpamAllowed: false,
			Status: Status{
				Delivered: ChannelFeatureNone,
				Read:      ChannelFeatureBoth,
			},
			Text: ChannelSettingsText{
				Creating: ChannelFeatureBoth,
				Editing:  ChannelFeatureBoth,
				Quoting:  ChannelFeatureBoth,
				Deleting: ChannelFeatureBoth,
			},
			Product: Product{
				Creating: ChannelFeatureSend,
				Deleting: ChannelFeatureSend,
			},
			Order: Order{
				Creating: ChannelFeatureBoth,
				Deleting: ChannelFeatureSend,
			},
			Image: ChannelSettingsFilesBase{
				Creating: ChannelFeatureBoth,
			},
			File: ChannelSettingsFilesBase{
				Creating: ChannelFeatureBoth,
			},
		},
	}

	data, status, err := c.UpdateTransportChannel(ch)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	t.Logf("Update selected channel: %v", data.ChannelID)
}

func TestMgClient_TransportTemplates(t *testing.T) {
	c := client()

	data, status, err := c.TransportTemplates()
	assert.NoError(t, err, fmt.Sprintf("%d %s", status, err))

	t.Logf("Templates found: %#v", len(data))

	for _, item := range data {
		for _, tpl := range item.Template {
			if tpl.Type == TemplateItemTypeText {
				assert.Empty(t, tpl.VarType)
			} else {
				assert.Empty(t, tpl.Text)
				assert.NotEmpty(t, tpl.VarType)

				if _, ok := templateVarAssoc[tpl.VarType]; !ok {
					t.Errorf("unknown TemplateVar type %s", tpl.VarType)
				}
			}
		}
	}
}

func TestMgClient_ActivateTemplate(t *testing.T) {
	c := client()
	req := ActivateTemplateRequest{
		Code: tplCode,
		Name: tplCode,
		Type: TemplateTypeText,
		Template: []TemplateItem{
			{
				Type: TemplateItemTypeText,
				Text: "Hello ",
			},
			{
				Type:    TemplateItemTypeVar,
				VarType: TemplateVarFirstName,
			},
			{
				Type: TemplateItemTypeText,
				Text: "!",
			},
		},
	}

	status, err := c.ActivateTemplate(templateChannel(t), req)
	assert.NoError(t, err, fmt.Sprintf("%d %s", status, err))

	t.Logf("Activated template with code `%s`", req.Code)
}

func TestMgClient_UpdateTemplate(t *testing.T) {
	c := client()
	tpl := Template{
		Code:      tplCode,
		ChannelID: templateChannel(t),
		Name:      "updated name",
		Enabled:   true,
		Type:      TemplateTypeText,
		Template: []TemplateItem{
			{
				Type: TemplateItemTypeText,
				Text: "Welcome ",
			},
			{
				Type:    TemplateItemTypeVar,
				VarType: TemplateVarFirstName,
			},
			{
				Type: TemplateItemTypeText,
				Text: "!",
			},
		},
	}

	status, err := c.UpdateTemplate(tpl)
	assert.NoError(t, err, fmt.Sprintf("%d %s", status, err))

	templates, status, err := c.TransportTemplates()
	assert.NoError(t, err, fmt.Sprintf("%d %s", status, err))

	for _, template := range templates {
		if template.Code == tpl.Code {
			assert.Equal(t, tpl.Name, template.Name)
		}
	}
}

func TestMgClient_UpdateTemplateFail(t *testing.T) {
	c := client()
	tpl := Template{
		Name:    "updated name",
		Enabled: true,
		Type:    TemplateTypeText,
		Template: []TemplateItem{
			{
				Type: TemplateItemTypeText,
				Text: "Welcome ",
			},
			{
				Type:    TemplateItemTypeVar,
				VarType: TemplateVarFirstName,
			},
			{
				Type: TemplateItemTypeText,
				Text: "!",
			},
		},
	}

	status, err := c.UpdateTemplate(tpl)
	assert.Error(t, err, fmt.Sprintf("%d %s", status, err))
}

func TestMgClient_DeactivateTemplate(t *testing.T) {
	c := client()
	status, err := c.DeactivateTemplate(templateChannel(t), tplCode)
	assert.NoError(t, err, fmt.Sprintf("%d %s", status, err))
}

func TestMgClient_TextMessages(t *testing.T) {
	c := client()
	t.Logf("%v", ext)

	snd := SendData{
		Message: Message{
			ExternalID: ext,
			Type:       MsgTypeText,
			Text:       "hello!",
		},
		Originator: OriginatorCustomer,
		Customer: Customer{
			ExternalID: "6",
			Nickname:   "octopus",
			Firstname:  "Joe",
		},
		Channel:        channelID,
		ExternalChatID: "24798237492374",
	}

	data, status, err := c.Messages(snd)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Time.String() == "" {
		t.Errorf("%v", err)
	}

	t.Logf("Message %v is sent", data.MessageID)
}

func TestMgClient_ImageMessages(t *testing.T) {
	c := client()
	t.Logf("%v", ext)

	uploadFileResponse, st, err := c.UploadFileByURL(UploadFileByUrlRequest{
		Url: "https://via.placeholder.com/1",
	})

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	snd := SendData{
		Message: Message{
			ExternalID: ext + "file",
			Type:       MsgTypeImage,
			Items:      []Item{{ID: uploadFileResponse.ID}},
		},
		Originator: OriginatorCustomer,
		Customer: Customer{
			ExternalID: "6",
			Nickname:   "octopus",
			Firstname:  "Joe",
		},
		Channel:        channelID,
		ExternalChatID: "24798237492374",
	}

	data, status, err := c.Messages(snd)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Time.String() == "" {
		t.Errorf("%v", err)
	}

	t.Logf("Message %v is sent", data.MessageID)
}

func TestMgClient_UpdateMessages(t *testing.T) {
	c := client()
	t.Logf("%v", ext)

	sndU := EditMessageRequest{
		EditMessageRequestMessage{
			ExternalID: ext,
			Text:       "hello hello!",
		},
		channelID,
	}

	dataU, status, err := c.UpdateMessages(sndU)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if dataU.Time.String() == "" {
		t.Errorf("%v", err)
	}

	t.Logf("Message %v updated", dataU.MessageID)
}

func TestMgClient_MarkMessageReadAndDelete(t *testing.T) {
	c := client()
	t.Logf("%v", ext)

	snd := MarkMessageReadRequest{
		MarkMessageReadRequestMessage{
			ExternalID: ext,
		},
		channelID,
	}

	_, status, err := c.MarkMessageRead(snd)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	t.Logf("Message ext marked as read")

	sndD := DeleteData{
		Message{
			ExternalID: ext,
		},
		channelID,
	}

	previousChatMessage, status, err := c.DeleteMessage(sndD)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	t.Logf("Message %v deleted", ext)
	if previousChatMessage != nil {
		t.Logf("Previous chat message %+v", *previousChatMessage)
	}

	sndD = DeleteData{
		Message{
			ExternalID: ext + "file",
		},
		channelID,
	}

	previousChatMessage, status, err = c.DeleteMessage(sndD)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	t.Logf("Message %v deleted", ext+"file")
	if previousChatMessage != nil {
		t.Logf("Previous chat message %+v", *previousChatMessage)
	}
}

func TestMgClient_DeactivateTransportChannel(t *testing.T) {
	c := client()
	deleteData, status, err := c.DeactivateTransportChannel(channelID)

	if err != nil {
		t.Errorf("%d %v", status, err)
	}

	if deleteData.DeactivatedAt.String() == "" {
		t.Errorf("%v", err)
	}

	t.Logf("Deactivate selected channel: %v", deleteData.ChannelID)
}

func TestMgClient_UploadFile(t *testing.T) {
	c := client()
	t.Logf("%v", ext)

	// 1x1 png picture
	img := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEX/TQBcNTh/AAAAAXRSTlPM0jRW/QAAAApJREFUeJxjYgAAAAYAAzY3fKgAAAAASUVORK5CYII="
	binary, err := base64.StdEncoding.DecodeString(img)
	if err != nil {
		t.Errorf("cannot convert base64 to binary: %s", err)
	}

	data, status, err := c.UploadFile(bytes.NewReader(binary))

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	t.Logf("Message %+v is sent", data)
}
