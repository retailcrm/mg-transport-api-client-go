package v1

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

type MGClientTest struct {
	suite.Suite
}

func TestMGClient(t *testing.T) {
	suite.Run(t, new(MGClientTest))
}

func (t *MGClientTest) client() *MgClient {
	c := New("https://mg-test.retailcrm.pro", "mg_token")
	c.Debug = true

	return c
}

func (t *MGClientTest) gock() *gock.Request {
	return gock.New("https://mg-test.retailcrm.pro").MatchHeader("x-transport-token", "mg_token")
}

func (t *MGClientTest) transportURL(path string) string {
	return "/api/transport/v1/" + strings.TrimLeft(path, "/")
}

func (t *MGClientTest) Test_TransportChannels() {
	c := t.client()
	chName := "WhatsApp Channel"
	createdAt := "2021-11-22T08:20:46.479979Z"

	defer gock.Off()
	t.gock().
		Get(t.transportURL("channels")).
		Reply(http.StatusOK).
		JSON([]ChannelListItem{{
			ID:         1,
			ExternalID: "external_id",
			Type:       "whatsapp",
			Name:       &chName,
			Settings: ChannelSettings{
				Status: Status{
					Delivered: ChannelFeatureNone,
					Read:      ChannelFeatureSend,
				},
				Text: ChannelSettingsText{
					Creating:      ChannelFeatureBoth,
					Editing:       ChannelFeatureBoth,
					Quoting:       ChannelFeatureBoth,
					Deleting:      ChannelFeatureReceive,
					MaxCharsCount: 4096,
				},
				Product: Product{
					Creating: ChannelFeatureReceive,
					Editing:  ChannelFeatureReceive,
				},
				Order: Order{
					Creating: ChannelFeatureReceive,
					Editing:  ChannelFeatureReceive,
				},
				File: ChannelSettingsFilesBase{
					Creating: ChannelFeatureBoth,
					Editing:  ChannelFeatureBoth,
					Quoting:  ChannelFeatureBoth,
					Deleting: ChannelFeatureReceive,
					Max:      1,
				},
				Image: ChannelSettingsFilesBase{
					Creating: ChannelFeatureBoth,
					Editing:  ChannelFeatureBoth,
					Quoting:  ChannelFeatureBoth,
					Deleting: ChannelFeatureReceive,
					Max:      1, // nolint:gomnd
				},
				Suggestions: ChannelSettingsSuggestions{
					Text:  ChannelFeatureBoth,
					Phone: ChannelFeatureBoth,
					Email: ChannelFeatureBoth,
				},
				CustomerExternalID: ChannelFeatureCustomerExternalIDPhone,
				SendingPolicy: SendingPolicy{
					NewCustomer: ChannelFeatureSendingPolicyTemplate,
				},
			},
			CreatedAt:     createdAt,
			UpdatedAt:     &createdAt,
			ActivatedAt:   createdAt,
			DeactivatedAt: nil,
			IsActive:      true,
		}})

	data, status, err := c.TransportChannels(Channels{Active: true})
	t.Require().NoError(err)
	t.Assert().Equal(http.StatusOK, status)

	t.Assert().Len(data, 1)
}

func (t *MGClientTest) Test_ActivateTransportChannel() {
	c := t.client()
	ch := Channel{
		ID:   1,
		Type: "telegram",
		Name: "@my_shopping_bot",
		Settings: ChannelSettings{
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

	defer gock.Off()
	t.gock().
		Post(t.transportURL("channels")).
		Reply(http.StatusCreated).
		JSON(ActivateResponse{
			ChannelID:   1,
			ExternalID:  "external_id_1",
			ActivatedAt: time.Now(),
		})

	data, status, err := c.ActivateTransportChannel(ch)
	t.Require().NoError(err)
	t.Assert().Equal(http.StatusCreated, status)

	t.Assert().Equal(uint64(1), data.ChannelID)
	t.Assert().Equal("external_id_1", data.ExternalID)
	t.Assert().NotEmpty(data.ActivatedAt.String())
}

func (t *MGClientTest) Test_ActivateNewTransportChannel() {
	c := t.client()
	ch := Channel{
		Type: "telegram",
		Name: "@my_shopping_bot",
		Settings: ChannelSettings{
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

	defer gock.Off()

	t.gock().
		Post(t.transportURL("channels")).
		Reply(http.StatusCreated).
		JSON(ActivateResponse{
			ChannelID:   1,
			ExternalID:  "external_id_1",
			ActivatedAt: time.Now(),
		})

	t.gock().
		Delete(t.transportURL("channels/1")).
		Reply(http.StatusOK).
		JSON(DeleteResponse{
			ChannelID:     1,
			DeactivatedAt: time.Now(),
		})

	data, status, err := c.ActivateTransportChannel(ch)
	t.Require().NoError(err)
	t.Assert().Equal(http.StatusCreated, status)

	t.Assert().Equal(uint64(1), data.ChannelID)
	t.Assert().Equal("external_id_1", data.ExternalID)
	t.Assert().NotEmpty(data.ActivatedAt.String())

	deleteData, status, err := c.DeactivateTransportChannel(data.ChannelID)
	t.Require().NoError(err)
	t.Assert().Equal(http.StatusOK, status)
	t.Assert().NotEmpty(deleteData.DeactivatedAt.String())
	t.Assert().Equal(uint64(1), deleteData.ChannelID)
}

func (t *MGClientTest) Test_UpdateTransportChannel() {
	c := t.client()
	ch := Channel{
		ID:   1,
		Name: "@my_shopping_bot_2",
		Settings: ChannelSettings{
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

	defer gock.Off()
	t.gock().
		Put(t.transportURL("channels/1")).
		Reply(http.StatusOK).
		JSON(UpdateResponse{
			ChannelID:  uint64(1),
			ExternalID: "external_id_1",
			UpdatedAt:  time.Now(),
		})

	data, status, err := c.UpdateTransportChannel(ch)
	t.Require().NoError(err)
	t.Assert().Equal(http.StatusOK, status)
	t.Assert().Equal(uint64(1), data.ChannelID)
	t.Assert().Equal("external_id_1", data.ExternalID)
	t.Assert().NotEmpty(data.UpdatedAt.String())
}

func (t *MGClientTest) Test_TransportTemplates() {
	c := t.client()

	defer gock.Off()
	t.gock().
		Get(t.transportURL("templates")).
		Reply(http.StatusOK).
		JSON([]Template{{
			Code:      "tpl_code",
			ChannelID: 1,
			Name:      "Test Template",
			Enabled:   true,
			Type:      TemplateTypeText,
			Template: []TemplateItem{
				{
					Type: TemplateItemTypeText,
					Text: "Hello, ",
				},
				{
					Type:    TemplateItemTypeVar,
					VarType: TemplateVarFirstName,
				},
				{
					Type: TemplateItemTypeText,
					Text: "! We're glad to see you back in our store.",
				},
			},
		}})

	data, status, err := c.TransportTemplates()
	t.Assert().NoError(err, fmt.Sprintf("%d %s", status, err))
	t.Assert().Equal(http.StatusOK, status)
	t.Assert().Len(data, 1)

	for _, item := range data {
		for _, tpl := range item.Template {
			if tpl.Type == TemplateItemTypeText {
				t.Assert().Empty(tpl.VarType)
			} else {
				t.Assert().Empty(tpl.Text)
				t.Assert().NotEmpty(tpl.VarType)

				if _, ok := templateVarAssoc[tpl.VarType]; !ok {
					t.T().Errorf("unknown TemplateVar type %s", tpl.VarType)
				}
			}
		}
	}
}

func (t *MGClientTest) Test_ActivateTemplate() {
	c := t.client()
	req := ActivateTemplateRequest{
		Code: "tplCode",
		Name: "tplCode",
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

	defer gock.Off()
	t.gock().
		Post("/channels/1/templates").
		Reply(http.StatusCreated).
		JSON(map[string]interface{}{})

	status, err := c.ActivateTemplate(1, req)
	t.Assert().NoError(err, fmt.Sprintf("%d %s", status, err))
	t.Assert().Equal(http.StatusCreated, status)
}

func (t *MGClientTest) Test_UpdateTemplate() {
	c := t.client()
	tpl := Template{
		Code:      "encodable#code",
		ChannelID: 1,
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

	defer gock.Off()

	t.gock().
		Filter(func(r *http.Request) bool {
			return r.Method == http.MethodPut &&
				r.URL.Path == "/api/transport/v1/channels/1/templates/encodable#code"
		}).
		Reply(http.StatusOK).
		JSON(map[string]interface{}{})

	t.gock().
		Get(t.transportURL("templates")).
		Reply(http.StatusOK).
		JSON([]Template{tpl})

	status, err := c.UpdateTemplate(tpl)
	t.Assert().NoError(err, fmt.Sprintf("%d %s", status, err))

	templates, status, err := c.TransportTemplates()
	t.Assert().NoError(err, fmt.Sprintf("%d %s", status, err))

	for _, template := range templates {
		if template.Code == tpl.Code {
			t.Assert().Equal(tpl.Name, template.Name)
		}
	}
}

func (t *MGClientTest) Test_UpdateTemplateFail() {
	c := t.client()
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

	defer gock.Off()
	t.gock().
		Reply(http.StatusBadRequest).
		JSON(map[string][]string{
			"errors": {"Some weird error message..."},
		})

	status, err := c.UpdateTemplate(tpl)
	t.Assert().Error(err, fmt.Sprintf("%d %s", status, err))
}

func (t *MGClientTest) Test_DeactivateTemplate() {
	c := t.client()

	defer gock.Off()
	t.gock().
		Filter(func(r *http.Request) bool {
			return r.Method == http.MethodDelete &&
				r.URL.Path == t.transportURL("channels/1/templates/test_template#code")
		}).
		Reply(http.StatusOK).
		JSON(map[string]interface{}{})

	status, err := c.DeactivateTemplate(1, "test_template#code")
	t.Assert().NoError(err, fmt.Sprintf("%d %s", status, err))
}

func (t *MGClientTest) Test_TextMessages() {
	c := t.client()

	snd := SendData{
		Message: Message{
			ExternalID: "external_id",
			Type:       MsgTypeText,
			Text:       "hello!",
		},
		Originator: OriginatorCustomer,
		Customer: Customer{
			ExternalID: "6",
			Nickname:   "octopus",
			Firstname:  "Joe",
		},
		Channel:        1,
		ExternalChatID: "24798237492374",
	}

	defer gock.Off()
	t.gock().
		Post(t.transportURL("messages")).
		Reply(http.StatusOK).
		JSON(MessagesResponse{
			MessageID: 1,
			Time:      time.Now(),
		})

	data, status, err := c.Messages(snd)
	t.Require().NoError(err)
	t.Assert().Equal(http.StatusOK, status)
	t.Assert().NotEmpty(data.Time.String())
	t.Assert().Equal(1, data.MessageID)
}

func (t *MGClientTest) Test_ImageMessages() {
	c := t.client()

	defer gock.Off()

	t.gock().
		Post(t.transportURL("files/upload_by_url")).
		Reply(http.StatusOK).
		JSON(UploadFileResponse{
			ID:        "1",
			Hash:      "1",
			Type:      "image/png",
			MimeType:  "",
			Size:      1024,
			CreatedAt: time.Now(),
		})

	t.gock().
		Post(t.transportURL("messages")).
		Reply(http.StatusOK).
		JSON(MessagesResponse{
			MessageID: 1,
			Time:      time.Now(),
		})

	uploadFileResponse, st, err := c.UploadFileByURL(UploadFileByUrlRequest{
		Url: "https://via.placeholder.com/1",
	})
	t.Require().NoError(err)
	t.Assert().Equal(http.StatusOK, st)
	t.Assert().Equal("1", uploadFileResponse.ID)

	snd := SendData{
		Message: Message{
			ExternalID: "file",
			Type:       MsgTypeImage,
			Items:      []Item{{ID: uploadFileResponse.ID}},
		},
		Originator: OriginatorCustomer,
		Customer: Customer{
			ExternalID: "6",
			Nickname:   "octopus",
			Firstname:  "Joe",
		},
		Channel:        1,
		ExternalChatID: "24798237492374",
	}

	data, status, err := c.Messages(snd)
	t.Require().NoError(err)
	t.Assert().Equal(http.StatusOK, status)
	t.Assert().NotEmpty(data.Time.String())
	t.Assert().Equal(1, data.MessageID)
}

func (t *MGClientTest) Test_UpdateMessages() {
	c := t.client()

	sndU := EditMessageRequest{
		EditMessageRequestMessage{
			ExternalID: "editing",
			Text:       "hello hello!",
		},
		1,
	}

	defer gock.Off()
	t.gock().
		Put(t.transportURL("messages")).
		Reply(http.StatusOK).
		JSON(MessagesResponse{
			MessageID: 1,
			Time:      time.Now(),
		})

	dataU, status, err := c.UpdateMessages(sndU)
	t.Require().NoError(err)
	t.Assert().Equal(http.StatusOK, status)
	t.Assert().NotEmpty(dataU.Time.String())
	t.Assert().Equal(1, dataU.MessageID)
}

func (t *MGClientTest) Test_MarkMessageReadAndDelete() {
	c := t.client()

	snd := MarkMessageReadRequest{
		MarkMessageReadRequestMessage{
			ExternalID: "external_1",
		},
		1,
	}

	defer gock.Off()
	t.gock().
		Post(t.transportURL("messages/read")).
		Reply(http.StatusOK).
		JSON(MarkMessageReadResponse{})

	t.gock().
		Delete(t.transportURL("messages")).
		JSON(DeleteData{
			Message: Message{
				ExternalID: "deleted",
			},
			Channel: 1,
		}).
		Reply(http.StatusOK).
		JSON(MessagesResponse{
			MessageID: 2,
			Time:      time.Now(),
		})

	_, status, err := c.MarkMessageRead(snd)
	t.Require().NoError(err)
	t.Assert().Equal(http.StatusOK, status)

	previousChatMessage, status, err := c.DeleteMessage(DeleteData{
		Message{
			ExternalID: "deleted",
		},
		1,
	})
	t.Require().NoError(err)
	t.Assert().Equal(http.StatusOK, status)
	t.Assert().Equal(2, previousChatMessage.MessageID)
}

func (t *MGClientTest) Test_DeactivateTransportChannel() {
	c := t.client()

	defer gock.Off()
	t.gock().
		Delete(t.transportURL("channels/1")).
		Reply(http.StatusOK).
		JSON(DeleteResponse{
			ChannelID:     1,
			DeactivatedAt: time.Now(),
		})

	deleteData, status, err := c.DeactivateTransportChannel(1)
	t.Require().NoError(err)
	t.Assert().Equal(http.StatusOK, status)
	t.Assert().NotEmpty(deleteData.DeactivatedAt.String())
	t.Assert().Equal(uint64(1), deleteData.ChannelID)
}

func (t *MGClientTest) Test_UploadFile() {
	c := t.client()

	// 1x1 png picture
	img := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEX/TQBcNTh/AAAAAXRSTlPM0jRW/QAAAApJREFUeJxjYgAAAAYAAzY3fKgAAAAASUVORK5CYII="
	binary, err := base64.StdEncoding.DecodeString(img)
	if err != nil {
		t.T().Errorf("cannot convert base64 to binary: %s", err)
	}

	resp := UploadFileResponse{
		ID:        "1",
		Hash:      "1",
		Type:      "image/png",
		MimeType:  "",
		Size:      1024,
		CreatedAt: time.Now(),
	}

	defer gock.Off()
	t.gock().
		Post(t.transportURL("files/upload")).
		Body(bytes.NewReader(binary)).
		Reply(http.StatusOK).
		JSON(resp)

	data, status, err := c.UploadFile(bytes.NewReader(binary))
	t.Require().NoError(err)
	t.Assert().Equal(http.StatusOK, status)

	resp.CreatedAt = data.CreatedAt
	t.Assert().Equal(resp, data)
}

func (t *MGClientTest) Test_SuccessHandleAPIError() {
	client := t.client()
	handleError := client.Error([]byte(`{"errors": ["Channel not found"]}`))

	t.Assert().IsType(new(httpClientError), handleError)
	t.Assert().Equal(handleError.Error(), "Channel not found")

	defer gock.Off()
	t.gock().
		Delete(t.transportURL("channels/123")).
		Reply(http.StatusInternalServerError)

	_, statusCode, err := client.DeactivateTransportChannel(123)

	t.Assert().Equal(http.StatusInternalServerError, statusCode)
	t.Assert().IsType(new(httpClientError), err)
	t.Assert().Equal("Internal server error", err.Error())
}
