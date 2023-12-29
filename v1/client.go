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

// New initializes the MgClient.
func New(url string, token string) *MgClient {
	return NewWithClient(url, token, &http.Client{Timeout: time.Minute})
}

// NewWithClient initializes the MgClient with specified *http.Client.
func NewWithClient(url string, token string, client *http.Client) *MgClient {
	return &MgClient{
		URL:        url,
		Token:      token,
		httpClient: client,
	}
}

// WithLogger sets the provided logger instance into the Client.
func (c *MgClient) WithLogger(logger BasicLogger) *MgClient {
	c.logger = logger
	return c
}

// writeLog writes a message to the log.
func (c *MgClient) writeLog(format string, v ...interface{}) {
	if c.logger != nil {
		c.logger.Printf(format, v...)
		return
	}

	log.Printf(format, v...)
}

// TransportTemplates returns templates list.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	data, status, err := client.TransportTemplates()
//	if err != nil {
//	    log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, response: %#v", status, data)
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

// ActivateTemplate activates template with provided structure.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	status, err := client.ActivateTemplate(1, ActivateTemplateRequest{
//		UpdateTemplateRequest: UpdateTemplateRequest{
//			Name:     "New Template",
//			Body:     "Hello, {{1}}! Welcome to our store!",
//			Lang:     "en",
//			Category: "marketing",
//			Example: &TemplateExample{
//				Header: []string{"https://example.com/image.png"},
//				Body:   []string{"John"},
//			},
//			VerificationStatus: TemplateStatusApproved,
//			Header: &TemplateHeader{
//				Content: HeaderContentImage{},
//			},
//		},
//		Code: "new_template",
//		Type: TemplateTypeMedia,
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d", status)
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

// UpdateTemplate updates existing template by its code.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	status, err := client.UpdateTemplate(1, "new_template", UpdateTemplateRequest{
//		Name:     "New Template",
//		Body:     "Hello, {{1}}! Welcome to our store!",
//		Lang:     "en",
//		Category: "marketing",
//		Example: &TemplateExample{
//			Header: []string{"https://example.com/image.png"},
//			Body:   []string{"John"},
//		},
//		VerificationStatus: TemplateStatusApproved,
//		Header: &TemplateHeader{
//			Content: HeaderContentImage{},
//		},
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d", status)
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

// DeactivateTemplate deactivates the template by its code.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	status, err := client.DeactivateTemplate(1, "new_template")
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d", status)
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

// TransportChannels returns channels for current transport.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	resp, status, err := client.TransportChannels(Channels{
//		Active: true,
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, channels: %#v", status, resp)
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

// ActivateTransportChannel activates the channel with provided settings.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//	uint16Ptr := func(val uint16) *uint16 {
//		return &val
//	}
//	mbToBytes := func(val uint64) *uint64 {
//		val = val * 1024 * 1024
//		return &val
//	}
//
//	resp, status, err := client.ActivateTransportChannel(Channel{
//		Type: "telegram",
//		Name: "@my_shopping_bot",
//		Settings: ChannelSettings{
//			Status: Status{
//				Delivered: ChannelFeatureNone,
//				Read:      ChannelFeatureReceive,
//			},
//			Text: ChannelSettingsText{
//				Creating:      ChannelFeatureBoth,
//				Editing:       ChannelFeatureBoth,
//				Quoting:       ChannelFeatureReceive,
//				Deleting:      ChannelFeatureSend,
//				MaxCharsCount: 2000,
//			},
//			Product: Product{
//				Creating: ChannelFeatureSend,
//				Editing:  ChannelFeatureNone,
//				Deleting: ChannelFeatureSend,
//			},
//			Order: Order{
//				Creating: ChannelFeatureBoth,
//				Editing:  ChannelFeatureNone,
//				Deleting: ChannelFeatureSend,
//			},
//			File: ChannelSettingsFilesBase{
//				Creating:          ChannelFeatureBoth,
//				Editing:           ChannelFeatureBoth,
//				Quoting:           ChannelFeatureBoth,
//				Deleting:          ChannelFeatureBoth,
//				Max:               10,
//				NoteMaxCharsCount: uint16Ptr(256),
//				MaxItemSize:       mbToBytes(50),
//			},
//			Image: ChannelSettingsFilesBase{
//				Creating:          ChannelFeatureBoth,
//				Editing:           ChannelFeatureBoth,
//				Quoting:           ChannelFeatureBoth,
//				Deleting:          ChannelFeatureBoth,
//				Max:               10,
//				NoteMaxCharsCount: uint16Ptr(256),
//				MaxItemSize:       mbToBytes(10),
//			},
//			Suggestions: ChannelSettingsSuggestions{
//				Text:  ChannelFeatureBoth,
//				Phone: ChannelFeatureBoth,
//				Email: ChannelFeatureBoth,
//			},
//			Audio: ChannelSettingsAudio{
//				Creating:    ChannelFeatureBoth,
//				Quoting:     ChannelFeatureBoth,
//				Deleting:    ChannelFeatureBoth,
//				MaxItemSize: mbToBytes(10),
//			},
//		},
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, channel external_id: %s", status, resp.ExternalID)
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

// UpdateTransportChannel updates an existing channel with provided settings.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//	uint16Ptr := func(val uint16) *uint16 {
//		return &val
//	}
//	mbToBytes := func(val uint64) *uint64 {
//		val = val * 1024 * 1024
//		return &val
//	}
//
//	resp, status, err := client.UpdateTransportChannel(Channel{
//		ID:   305,
//		Type: "telegram",
//		Name: "@my_shopping_bot",
//		Settings: ChannelSettings{
//			Status: Status{
//				Delivered: ChannelFeatureNone,
//				Read:      ChannelFeatureReceive,
//			},
//			Text: ChannelSettingsText{
//				Creating:      ChannelFeatureBoth,
//				Editing:       ChannelFeatureBoth,
//				Quoting:       ChannelFeatureReceive,
//				Deleting:      ChannelFeatureSend,
//				MaxCharsCount: 2000,
//			},
//			Product: Product{
//				Creating: ChannelFeatureSend,
//				Editing:  ChannelFeatureNone,
//				Deleting: ChannelFeatureSend,
//			},
//			Order: Order{
//				Creating: ChannelFeatureBoth,
//				Editing:  ChannelFeatureNone,
//				Deleting: ChannelFeatureSend,
//			},
//			File: ChannelSettingsFilesBase{
//				Creating:          ChannelFeatureBoth,
//				Editing:           ChannelFeatureBoth,
//				Quoting:           ChannelFeatureBoth,
//				Deleting:          ChannelFeatureBoth,
//				Max:               10,
//				NoteMaxCharsCount: uint16Ptr(256),
//				MaxItemSize:       mbToBytes(50),
//			},
//			Image: ChannelSettingsFilesBase{
//				Creating:          ChannelFeatureBoth,
//				Editing:           ChannelFeatureBoth,
//				Quoting:           ChannelFeatureBoth,
//				Deleting:          ChannelFeatureBoth,
//				Max:               10,
//				NoteMaxCharsCount: uint16Ptr(256),
//				MaxItemSize:       mbToBytes(10),
//			},
//			Suggestions: ChannelSettingsSuggestions{
//				Text:  ChannelFeatureBoth,
//				Phone: ChannelFeatureBoth,
//				Email: ChannelFeatureBoth,
//			},
//			Audio: ChannelSettingsAudio{
//				Creating:    ChannelFeatureBoth,
//				Quoting:     ChannelFeatureBoth,
//				Deleting:    ChannelFeatureBoth,
//				MaxItemSize: mbToBytes(10),
//			},
//		},
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, channel_id: %d", status, resp.ChannelID)
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

// DeactivateTransportChannel deactivates the channel by its ID.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	resp, status, err := client.DeactivateTransportChannel(305)
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, deactivated at: %s", status, resp.DeactivatedAt)
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

// Messages sends new message.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//	getReplyDeadline := func(after time.Duration) *time.Time {
//		deadline := time.Now().Add(after)
//		return &deadline
//	}
//
//	resp, status, err := client.Messages(SendData{
//		Message: Message{
//			ExternalID: "uid_1",
//			Type:       MsgTypeText,
//			Text:       "Hello customer!",
//			PageLink:   "https://example.com",
//		},
//		Originator: OriginatorCustomer,
//		Customer: Customer{
//			ExternalID: "client_id_1",
//			Nickname:   "customer",
//			Firstname:  "Tester",
//			Lastname:   "Tester",
//			Avatar:     "https://example.com/image.png",
//			ProfileURL: "https://example.com/user/client_id_1",
//			Language:   "en",
//			Utm: &Utm{
//				Source:   "myspace.com",
//				Medium:   "social",
//				Campaign: "something",
//				Term:     "fedora",
//				Content:  "autumn_collection",
//			},
//		},
//		Channel:        305,
//		ExternalChatID: "chat_id_1",
//		ReplyDeadline: getReplyDeadline(24 * time.Hour),
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, message ID: %d", status, resp.MessageID)
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

// MessagesHistory sends history message.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//	getModifiedNow := func(after time.Duration) *time.Time {
//		deadline := time.Now().Add(after)
//		return &deadline
//	}
//
//	resp, status, err := client.MessagesHistory(SendHistoryMessageRequest{
//		Message: SendMessageRequestMessage{
//			ExternalID: "uid_1",
//			Type:       MsgTypeText,
//			Text:       "Hello customer!",
//			CreatedAt:  getModifiedNow(-time.Hour),
//		},
//		Originator: OriginatorCustomer,
//		Customer: &Customer{
//			ExternalID: "client_id_1",
//			Nickname:   "customer",
//			Firstname:  "Tester",
//			Lastname:   "Tester",
//			Avatar:     "https://example.com/image.png",
//			ProfileURL: "https://example.com/user/client_id_1",
//			Language:   "en",
//			Utm: &Utm{
//				Source:   "myspace.com",
//				Medium:   "social",
//				Campaign: "something",
//				Term:     "fedora",
//				Content:  "autumn_collection",
//			},
//		},
//		ChannelID:      305,
//		ExternalChatID: "chat_id_1",
//		ReplyDeadline:  getModifiedNow(24 * time.Hour),
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, message ID: %d", status, resp.MessageID)
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

// UpdateMessages edits existing message. Only text messages are supported.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	resp, status, err := client.UpdateMessages(EditMessageRequest{
//		Message: EditMessageRequestMessage{
//			ExternalID: "message_id_1",
//			Text:       "This is a new text!",
//		},
//		Channel: 305,
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, message ID: %d", status, resp.MessageID)
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

// MarkMessageRead send message read event to MG.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	_, status, err := client.MarkMessageRead(MarkMessageReadRequest{
//		Message: MarkMessageReadRequestMessage{
//			ExternalID: "message_id_1",
//		},
//		ChannelID: 305,
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d", status)
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

// AckMessage sets success status for message or appends an error to message.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	status, err := client.AckMessage(AckMessageRequest{
//		ExternalMessageID: "message_id_1",
//		Channel: 305,
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d", status)
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
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	resp, status, err := client.ReadUntil(MarkMessagesReadUntilRequest{
//		CustomerExternalID: "customer_id_1",
//		ChannelID: 305,
//		Until: time.Now().Add(-time.Hour),
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, marked messages: %+v", status, resp.IDs)
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

// DeleteMessage removes the message.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	resp, status, err := client.DeleteMessage(DeleteData{
//		Message: Message{
//			ExternalID: "message_id_1",
//		},
//		Channel: 305,
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, message ID: %d", status, resp.MessageID)
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

// GetFile returns file information by its ID.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	resp, status, err := client.GetFile("file_id")
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, file URL: %s", status, resp.Url)
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

// UploadFile uploads a file.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	file, err := os.Open("/tmp/file.png")
//	if err != nil {
//		log.Fatalf("cannot open file for reading: %s", err)
//	}
//	defer func() { _ = file.Close() }()
//
//	data, err := io.ReadAll(file)
//	if err != nil {
//		log.Fatalf("cannot read file data: %s", err)
//	}
//
//	resp, status, err := client.UploadFile(bytes.NewReader(data))
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, file ID: %s", status, resp.ID)
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

// UploadFileByURL uploads a file from provided URL.
//
// Example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//
//	resp, status, err := client.UploadFileByURL(UploadFileByUrlRequest{
//		Url: "https://example.com/file.png",
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, file ID: %s", status, resp.ID)
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

// MakeTimestamp returns current unix timestamp in milliseconds.
//
// Example:
//
//	fmt.Printf("UNIX timestamp in milliseconds: %d", MakeTimestamp())
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
