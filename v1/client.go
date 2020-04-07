package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
)

// New initialize client
func New(url string, token string) *MgClient {
	return NewWithClient(url, token, &http.Client{Timeout: time.Minute})
}

// NewWithClient initializes client with provided http client
func NewWithClient(url string, token string, client *http.Client) *MgClient {
	return &MgClient{
		URL:        url,
		Token:      token,
		httpClient: client,
	}
}

// TransportChannels returns channels list
//
// Example:
//
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
// 	data, status, err := client.TransportChannels(Channels{Active: true})
//
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
//
//	fmt.Printf("Status: %v, Channels found: %v", status, len(data))
func (c *MgClient) TransportChannels(request Channels) ([]ChannelListItem, int, error) {
	var resp []ChannelListItem
	var b []byte
	outgoing, _ := query.Values(request)

	data, status, err := c.GetRequest(fmt.Sprintf("/channels?%s", outgoing.Encode()), b)
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status > http.StatusCreated || status < http.StatusOK {
		return resp, status, c.Error(data)
	}

	return resp, status, err
}

// TransportTemplates returns templates list
//
// Example:
//
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
// 	data, status, err := client.TransportTemplates()
//
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
//
//	fmt.Printf("Status: %v, Templates found: %v", status, len(data))
func (c *MgClient) TransportTemplates() ([]TemplateItem, int, error) {
	var resp []TemplateItem

	data, status, err := c.GetRequest("/templates", []byte{})
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status > http.StatusCreated || status < http.StatusOK {
		return resp, status, c.Error(data)
	}

	return resp, status, err
}

// ActivateTransportChannel implements template activation
//
// Example:
// 		var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
// 		request := v1.ActivateTemplateRequest{
// 				Code: "code",
// 				Name: "name",
// 				Type: v1.TemplateTypeText,
// 				Template: []v1.TemplateItem{
// 					{
// 						Type: v1.TemplateItemTypeText,
// 						Text: "Hello, ",
// 					},
// 					{
// 						Type:    v1.TemplateItemTypeVar,
// 						VarType: v1.TemplateVarName,
// 					},
// 					{
// 						Type: v1.TemplateItemTypeText,
// 						Text: "!",
// 					},
// 				},
// 		}
//
// 		_, err := client.ActivateTemplate(uint64(1), request)
//
// 		if err != nil {
// 			fmt.Printf("%v", err)
// 		}
func (c *MgClient) ActivateTemplate(channelID uint64, request ActivateTemplateRequest) (int, error) {
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PostRequest(fmt.Sprintf("/channels/%d/templates", channelID), bytes.NewBuffer(outgoing))
	if err != nil {
		return status, err
	}

	if status > http.StatusCreated || status < http.StatusOK {
		return status, c.Error(data)
	}

	return status, err
}

// UpdateTemplate implements template updating
// Example:
// 		var client = New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
// 		request := v1.Template{
// 			Code:      "templateCode",
// 			ChannelID: 1,
// 			Name:      "templateName",
// 			Template:  []v1.TemplateItem{
// 				{
// 					Type: v1.TemplateItemTypeText,
// 					Text: "Welcome, ",
// 				},
// 				{
// 					Type: v1.TemplateItemTypeVar,
// 					VarType: v1.TemplateVarName,
// 				},
// 				{
// 					Type: v1.TemplateItemTypeText,
// 					Text: "!",
// 				},
// 			},
// 		}
//
// 		_, err := client.UpdateTemplate(request)
//
// 		if err != nil {
// 			fmt.Printf("%#v", err)
// 		}
func (c *MgClient) UpdateTemplate(request Template) (int, error) {
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PutRequest(fmt.Sprintf("/channels/%d/templates/%s", request.ChannelID, request.Code), outgoing)
	if err != nil {
		return status, err
	}

	if status != http.StatusOK {
		return status, c.Error(data)
	}

	return status, err
}

// DeactivateTemplate implements template deactivation
//
// Example:
//
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
// 	_, err := client.DeactivateTemplate(3053450384, "templateCode")
//
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
func (c *MgClient) DeactivateTemplate(channelID uint64, templateCode string) (int, error) {
	data, status, err := c.DeleteRequest(
		fmt.Sprintf("/channels/%d/templates/%s", channelID, templateCode), []byte{})
	if err != nil {
		return status, err
	}

	if status > http.StatusCreated || status < http.StatusOK {
		return status, c.Error(data)
	}

	return status, err
}

// ActivateTransportChannel implement channel activation
//
// Example:
//
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	request := ActivateRequest{
//		Type: "telegram",
//		Name: "@my_shopping_bot",
//		Settings: ChannelSettings{
//			SpamAllowed: false,
//			Status: Status{
//				Delivered: ChannelFeatureNone,
//				Read: ChannelFeatureReceive,
// 			},
//			Text: ChannelSettingsText{
//				Creating: ChannelFeatureBoth,
//				Editing:  ChannelFeatureBoth,
//				Quoting:  ChannelFeatureReceive,
//				Deleting: ChannelFeatureSend,
//				MaxCharsCount: 2000,
// 			},
//			Product: Product{
//				Creating: ChannelFeatureSend,
//				Deleting: ChannelFeatureSend,
// 			},
//			Order: Order{
//				Creating: ChannelFeatureBoth,
//				Deleting: ChannelFeatureSend,
// 			},
// 		},
//	}
//
// 	data, status, err := client.ActivateTransportChannel(request)
//
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
//
//	fmt.Printf("%s\n", data.CreatedAt)
func (c *MgClient) ActivateTransportChannel(request Channel) (ActivateResponse, int, error) {
	var resp ActivateResponse
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PostRequest("/channels", bytes.NewBuffer(outgoing))
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status > http.StatusCreated || status < http.StatusOK {
		return resp, status, c.Error(data)
	}

	return resp, status, err
}

// UpdateTransportChannel implement channel activation
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	request := ActivateRequest{
//		ID:   3053450384,
//		Type: "telegram",
//		Name: "@my_shopping_bot",
//		Settings: ChannelSettings{
//			SpamAllowed: false,
//			Status: Status{
//				Delivered: ChannelFeatureNone,
//				Read: ChannelFeatureReceive,
// 			},
//			Text: ChannelSettingsText{
//				Creating: ChannelFeatureBoth,
//				Editing:  ChannelFeatureSend,
//				Quoting:  ChannelFeatureReceive,
//				Deleting: ChannelFeatureBoth,
// 			},
//			Product: Product{
//				Creating: ChannelFeatureSend,
//				Deleting: ChannelFeatureSend,
// 			},
//			Order: Order{
//				Creating: ChannelFeatureBoth,
//				Deleting: ChannelFeatureSend,
// 			},
// 		},
//	}
//
// 	data, status, err := client.UpdateTransportChannel(request)
//
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
//
//	fmt.Printf("%s\n", data.UpdatedAt)
func (c *MgClient) UpdateTransportChannel(request Channel) (UpdateResponse, int, error) {
	var resp UpdateResponse
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PutRequest(fmt.Sprintf("/channels/%d", request.ID), []byte(outgoing))
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status != http.StatusOK {
		return resp, status, c.Error(data)
	}

	return resp, status, err
}

// DeactivateTransportChannel implement channel deactivation
//
// Example:
//
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
// 	data, status, err := client.DeactivateTransportChannel(3053450384)
//
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
//
//	fmt.Printf("%s\n", data.DeactivatedAt)
func (c *MgClient) DeactivateTransportChannel(id uint64) (DeleteResponse, int, error) {
	var resp DeleteResponse
	var buf []byte

	data, status, err := c.DeleteRequest(
		fmt.Sprintf("/channels/%s", strconv.FormatUint(id, 10)),
		buf,
	)
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status != http.StatusOK {
		return resp, status, c.Error(data)
	}

	return resp, status, err
}

// Messages implement send message
//
// Example:
//
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//	msg := SendData{
//		SendMessage{
//			Message{
//				ExternalID: "274628",
//				Type:       "text",
//				Text:       "hello!",
//			},
//			time.Now(),
//		},
//		User{
//			ExternalID: "8",
//			Nickname:   "@octopus",
//			Firstname:  "Joe",
//		},
//		10,
//	}
//
// 	data, status, err := client.Messages(msg)
//
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
//
//	fmt.Printf("%s\n", data.MessageID)
func (c *MgClient) Messages(request SendData) (MessagesResponse, int, error) {
	var resp MessagesResponse
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PostRequest("/messages", bytes.NewBuffer(outgoing))
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status != http.StatusOK {
		return resp, status, c.Error(data)
	}

	return resp, status, err
}

// UpdateMessages implement edit message
//
// Example:
//
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//	msg := UpdateData{
//		UpdateMessage{
//			Message{
//				ExternalID: "274628",
//				Type:       "text",
//				Text:       "hello hello!",
//			},
//			MakeTimestamp(),
//		},
//		10,
//	}
//
// 	data, status, err := client.UpdateMessages(msg)
//
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
//
//	fmt.Printf("%s\n", data.MessageID)
func (c *MgClient) UpdateMessages(request EditMessageRequest) (MessagesResponse, int, error) {
	var resp MessagesResponse
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PutRequest("/messages", []byte(outgoing))
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status != http.StatusOK {
		return resp, status, c.Error(data)
	}

	return resp, status, err
}

// MarkMessageRead send message read event to MG
//
// Example:
//
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//	msg := MarkMessageReadRequest{
//		Message{
//			ExternalID: "274628",
//		},
//		10,
//	}
//
// 	data, status, err := client.MarkMessageRead(msg)
//
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
//
//	fmt.Printf("%v %v\n", status, data)
func (c *MgClient) MarkMessageRead(request MarkMessageReadRequest) (MarkMessageReadResponse, int, error) {
	var resp MarkMessageReadResponse
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PostRequest("/messages/read", bytes.NewBuffer(outgoing))
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status != http.StatusOK {
		return resp, status, c.Error(data)
	}

	return resp, status, err
}

// DeleteMessage implement delete message
//
// Example:
//
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	msg := DeleteData{
//		Message{
//			ExternalID: "274628",
//		},
//		10,
//	}
//
// 	data, status, err := client.DeleteMessage(msg)
//
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
//
//	fmt.Printf("%s\n", data.MessageID)
func (c *MgClient) DeleteMessage(request DeleteData) (MessagesResponse, int, error) {
	var resp MessagesResponse
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.DeleteRequest(
		"/messages",
		[]byte(outgoing),
	)
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status != http.StatusOK {
		return resp, status, c.Error(data)
	}

	return resp, status, err
}

// GetFile implement get file url
//
// Example:
//
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
// 	data, status, err := client.GetFile("file_ID")
//
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
//
//	fmt.Printf("%s\n", data.MessageID)
func (c *MgClient) GetFile(request string) (FullFileResponse, int, error) {
	var resp FullFileResponse
	var b []byte

	data, status, err := c.GetRequest(fmt.Sprintf("/files/%s", request), b)

	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status != http.StatusOK {
		return resp, status, c.Error(data)
	}

	return resp, status, err
}

// UploadFile upload file
func (c *MgClient) UploadFile(request io.Reader) (UploadFileResponse, int, error) {
	var resp UploadFileResponse

	data, status, err := c.PostRequest("/files/upload", request)
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status != http.StatusOK {
		return resp, status, c.Error(data)
	}

	return resp, status, err
}

// UploadFileByURL upload file by url
func (c *MgClient) UploadFileByURL(request UploadFileByUrlRequest) (UploadFileResponse, int, error) {
	var resp UploadFileResponse
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PostRequest("/files/upload_by_url", bytes.NewBuffer(outgoing))
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status != http.StatusOK {
		return resp, status, c.Error(data)
	}

	return resp, status, err
}

func (c *MgClient) Error(info []byte) error {
	var data map[string]interface{}

	if err := json.Unmarshal(info, &data); err != nil {
		return err
	}

	values := data["errors"].([]interface{})

	return errors.New(values[0].(string))
}

// MakeTimestamp returns current unix timestamp
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
