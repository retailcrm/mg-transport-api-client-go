package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
)

// New initialize client.
func New(url string, token string) *MgClient {
	return NewWithClient(url, token, &http.Client{Timeout: time.Minute})
}

// NewWithClient initializes client with provided http client.
func NewWithClient(url string, token string, client *http.Client) *MgClient {
	return &MgClient{
		URL:        url,
		Token:      token,
		httpClient: client,
	}
}

// WithLogger sets the provided logger instance into the Client
func (c *MgClient) WithLogger(logger BasicLogger) *MgClient {
	c.logger = logger
	return c
}

// writeLog writes to the log.
func (c *MgClient) writeLog(format string, v ...interface{}) {
	if c.logger != nil {
		c.logger.Printf(format, v...)
		return
	}

	log.Printf(format, v...)
}

// TransportTemplates returns templates list
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	data, status, err := client.TransportTemplates()
//
//	if err != nil {
//		fmt.Printf("%v", err)
//	}
//
//	fmt.Printf("Status: %v, Templates found: %v", status, len(data))
func (c *MgClient) TransportTemplates() ([]Template, int, error) {
	var resp []Template

	data, status, err := c.GetRequest("/templates", []byte{})
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status > http.StatusCreated || status < http.StatusOK {
		return resp, status, NewAPIClientError(data)
	}

	return resp, status, err
}

// ActivateTemplate implements template activation
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	request := v1.ActivateTemplateRequest{
//			Code: "code",
//			Name: "name",
//			Type: v1.TemplateTypeText,
//			Template: []v1.TemplateItem{
//				{
//					Type: v1.TemplateItemTypeText,
//					Text: "Hello, ",
//				},
//				{
//					Type:    v1.TemplateItemTypeVar,
//					VarType: v1.TemplateVarName,
//				},
//				{
//					Type: v1.TemplateItemTypeText,
//					Text: "!",
//				},
//			},
//	}
//
//	_, err := client.ActivateTemplate(uint64(1), request)
//
//	if err != nil {
//		fmt.Printf("%v", err)
//	}
func (c *MgClient) ActivateTemplate(channelID uint64, request ActivateTemplateRequest) (int, error) {
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PostRequest(fmt.Sprintf("/channels/%d/templates", channelID), bytes.NewBuffer(outgoing))
	if err != nil {
		return status, err
	}

	if status > http.StatusCreated || status < http.StatusOK {
		return status, NewAPIClientError(data)
	}

	return status, err
}

// UpdateTemplate implements template updating
// Example:
//
//	var client = New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	request := v1.Template{
//		Code:      "templateCode",
//		ChannelID: 1,
//		Name:      "templateName",
//		Template:  []v1.TemplateItem{
//			{
//				Type: v1.TemplateItemTypeText,
//				Text: "Welcome, ",
//			},
//			{
//				Type: v1.TemplateItemTypeVar,
//				VarType: v1.TemplateVarName,
//			},
//			{
//				Type: v1.TemplateItemTypeText,
//				Text: "!",
//			},
//		},
//	}
//
//	_, err := client.UpdateTemplate(request)
//
//	if err != nil {
//		fmt.Printf("%#v", err)
//	}
func (c *MgClient) UpdateTemplate(channelID uint64, code string, request UpdateTemplateRequest) (int, error) {
	outgoing, _ := json.Marshal(&request)

	if channelID == 0 || code == "" {
		return 0, errors.New("`ChannelID` and `Code` cannot be blank")
	}

	data, status, err := c.PutRequest(
		fmt.Sprintf("/channels/%d/templates/%s", channelID, url.PathEscape(code)), outgoing)
	if err != nil {
		return status, err
	}

	if status != http.StatusOK {
		return status, NewAPIClientError(data)
	}

	return status, err
}

// DeactivateTemplate implements template deactivation
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	_, err := client.DeactivateTemplate(3053450384, "templateCode")
//
//	if err != nil {
//		fmt.Printf("%v", err)
//	}
func (c *MgClient) DeactivateTemplate(channelID uint64, templateCode string) (int, error) {
	data, status, err := c.DeleteRequest(
		fmt.Sprintf("/channels/%d/templates/%s", channelID, url.PathEscape(templateCode)), []byte{})
	if err != nil {
		return status, err
	}

	if status > http.StatusCreated || status < http.StatusOK {
		return status, NewAPIClientError(data)
	}

	return status, err
}

// TransportChannels returns channels list
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	data, status, err := client.TransportChannels(Channels{Active: true})
//
//	if err != nil {
//		fmt.Printf("%v", err)
//	}
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
		return resp, status, NewAPIClientError(data)
	}

	return resp, status, err
}

// ActivateTransportChannel implement channel activation
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	request := ActivateRequest{
//		Type: "telegram",
//		Name: "@my_shopping_bot",
//		Settings: ChannelSettings{
//			Status: Status{
//				Delivered: ChannelFeatureNone,
//				Read: ChannelFeatureReceive,
//			},
//			Text: ChannelSettingsText{
//				Creating: ChannelFeatureBoth,
//				Editing:  ChannelFeatureBoth,
//				Quoting:  ChannelFeatureReceive,
//				Deleting: ChannelFeatureSend,
//				MaxCharsCount: 2000,
//			},
//			Product: Product{
//				Creating: ChannelFeatureSend,
//				Deleting: ChannelFeatureSend,
//			},
//			Order: Order{
//				Creating: ChannelFeatureBoth,
//				Deleting: ChannelFeatureSend,
//			},
//		},
//	}
//
//	data, status, err := client.ActivateTransportChannel(request)
//
//	if err != nil {
//		fmt.Printf("%v", err)
//	}
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
		return resp, status, NewAPIClientError(data)
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
//			Status: Status{
//				Delivered: ChannelFeatureNone,
//				Read: ChannelFeatureReceive,
//			},
//			Text: ChannelSettingsText{
//				Creating: ChannelFeatureBoth,
//				Editing:  ChannelFeatureSend,
//				Quoting:  ChannelFeatureReceive,
//				Deleting: ChannelFeatureBoth,
//			},
//			Product: Product{
//				Creating: ChannelFeatureSend,
//				Deleting: ChannelFeatureSend,
//			},
//			Order: Order{
//				Creating: ChannelFeatureBoth,
//				Deleting: ChannelFeatureSend,
//			},
//		},
//	}
//
//	data, status, err := client.UpdateTransportChannel(request)
//
//	if err != nil {
//		fmt.Printf("%v", err)
//	}
//
//	fmt.Printf("%s\n", data.UpdatedAt)
func (c *MgClient) UpdateTransportChannel(request Channel) (UpdateResponse, int, error) {
	var resp UpdateResponse
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PutRequest(fmt.Sprintf("/channels/%d", request.ID), outgoing)
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status != http.StatusOK {
		return resp, status, NewAPIClientError(data)
	}

	return resp, status, err
}

// DeactivateTransportChannel implement channel deactivation
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	data, status, err := client.DeactivateTransportChannel(3053450384)
//
//	if err != nil {
//		fmt.Printf("%v", err)
//	}
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
		return resp, status, NewAPIClientError(data)
	}

	return resp, status, err
}

// Messages implement send message
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
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
//	data, status, err := client.Messages(msg)
//
//	if err != nil {
//		fmt.Printf("%v", err)
//	}
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
		return resp, status, NewAPIClientError(data)
	}

	return resp, status, err
}

// MessagesHistory implement history message sending.
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//	msg := v1.SendHistoryMessageRequest{
//		Message: v1.SendMessageRequestMessage{
//			Type:       v1.MsgTypeText,
//			ExternalID: "external_id",
//			CreatedAt:  v1.TimePtr(time.Now()),
//			IsComment:  false,
//			Text:       "Test message",
//		},
//		ChannelID:      1,
//		ExternalChatID: "chat_id",
//		Customer: &v1.Customer{
//			ExternalID: "1",
//			Nickname:   "@john_doe",
//			Firstname:  "John",
//			Lastname:   "Doe",
//		},
//		Originator:    v1.OriginatorCustomer,
//		ReplyDeadline: v1.TimePtr(time.Now().Add(time.Hour * 24)),
//	}
//
//	data, status, err := client.MessagesHistory(msg)
//	if err != nil {
//		fmt.Printf("[%d]: %v", status, err)
//	}
//
//	fmt.Printf("%d\n", data.MessageID)
func (c *MgClient) MessagesHistory(request SendHistoryMessageRequest) (MessagesResponse, int, error) {
	var (
		resp     MessagesResponse
		outgoing = &bytes.Buffer{}
	)
	_ = json.NewEncoder(outgoing).Encode(request)

	data, status, err := c.PostRequest("/messages/history", outgoing)
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status != http.StatusOK {
		return resp, status, NewAPIClientError(data)
	}

	return resp, status, err
}

// UpdateMessages implement edit message
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
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
//	data, status, err := client.UpdateMessages(msg)
//
//	if err != nil {
//		fmt.Printf("%v", err)
//	}
//
//	fmt.Printf("%s\n", data.MessageID)
func (c *MgClient) UpdateMessages(request EditMessageRequest) (MessagesResponse, int, error) {
	var resp MessagesResponse
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PutRequest("/messages", outgoing)
	if err != nil {
		return resp, status, err
	}

	if e := json.Unmarshal(data, &resp); e != nil {
		return resp, status, e
	}

	if status != http.StatusOK {
		return resp, status, NewAPIClientError(data)
	}

	return resp, status, err
}

// MarkMessageRead send message read event to MG
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//	msg := MarkMessageReadRequest{
//		Message{
//			ExternalID: "274628",
//		},
//		10,
//	}
//
//	data, status, err := client.MarkMessageRead(msg)
//
//	if err != nil {
//		fmt.Printf("%v", err)
//	}
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
		return resp, status, NewAPIClientError(data)
	}

	return resp, status, err
}

// AckMessage implements ack of message
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	request := AckMessageRequest{
//		ExternalMessageID: "274628",
//		Channel: 10,
//	}
//
//	status, err := client.AckMessage(request)
//
//	if err != nil {
//		fmt.Printf("%v", err)
//	}
func (c *MgClient) AckMessage(request AckMessageRequest) (int, error) {
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PostRequest("/messages/ack", bytes.NewBuffer(outgoing))
	if err != nil {
		return status, err
	}

	if status != http.StatusOK {
		return status, NewAPIClientError(data)
	}

	return status, err
}

// ReadUntil will mark all messages from specified timestamp as read.
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	request := ReadUntilRequest{
//		ExternalMessageID: "274628",
//		Channel: 10,
//	}
//
//	resp, status, err := client.ReadUntil(request)
//	if err != nil {
//		fmt.Printf("%v", err)
//	}
//	if resp != nil {
//		fmt.Printf("Marked these as read: %s", resp.IDs)
//	}
func (c *MgClient) ReadUntil(request MarkMessagesReadUntilRequest) (*MarkMessagesReadUntilResponse, int, error) {
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PostRequest("/messages/read_until", bytes.NewBuffer(outgoing))
	if err != nil {
		return nil, status, err
	}
	if status != http.StatusOK {
		return nil, status, NewAPIClientError(data)
	}

	var resp *MarkMessagesReadUntilResponse
	if e := json.Unmarshal(data, &resp); e != nil {
		return nil, status, e
	}
	return resp, status, nil
}

// DeleteMessage implement delete message
//
// Example:
//
//		var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//		msg := DeleteData{
//			Message{
//				ExternalID: "274628",
//			},
//			10,
//		}
//
//		previousChatMessage, status, err := client.DeleteMessage(msg)
//		if err != nil {
//			fmt.Printf("%v", err)
//		}
//
//	 if previousChatMessage != nil {
//	 	fmt.Printf("Previous chat message id = %d", previousChatMessage.MessageID)
//	 }
func (c *MgClient) DeleteMessage(request DeleteData) (*MessagesResponse, int, error) {
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.DeleteRequest(
		"/messages",
		outgoing,
	)
	if err != nil {
		return nil, status, err
	}
	if status != http.StatusOK {
		return nil, status, NewAPIClientError(data)
	}

	var previousChatMessage *MessagesResponse
	if e := json.Unmarshal(data, &previousChatMessage); e != nil {
		return nil, status, e
	}

	return previousChatMessage, status, nil
}

// GetFile implement get file url
//
// Example:
//
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	data, status, err := client.GetFile("file_ID")
//
//	if err != nil {
//		fmt.Printf("%v", err)
//	}
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
		return resp, status, NewAPIClientError(data)
	}

	return resp, status, err
}

// UploadFile upload file.
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
		return resp, status, NewAPIClientError(data)
	}

	return resp, status, err
}

// UploadFileByURL upload file by url.
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
		return resp, status, NewAPIClientError(data)
	}

	return resp, status, err
}

// MakeTimestamp returns current unix timestamp.
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
