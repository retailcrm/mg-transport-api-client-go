package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// New initialize client
func New(url string, token string) *MgClient {
	return &MgClient{
		url,
		token,
		&http.Client{Timeout: 20 * time.Second},
	}
}

// ActivateTransportChannel implement channel activation
//
// Example:
//
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d49bcba99be73bff503ea6")
//
//	request := ActivateRequest{
//		Type: "telegram",
//		Events: [2]int{"message_sent", "message_sent"}
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

	data, status, err := c.PostRequest("/transport/channels", []byte(outgoing))
	if err != nil {
		return resp, status, err
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, status, err
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
//	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d49bcba99be73bff503ea6")
//
//	request := ActivateRequest{
//		ID:   3053450384,
//		Type: "telegram",
//		Events: [2]int{"message_sent", "message_sent"}
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

	data, status, err := c.PutRequest(fmt.Sprintf("/transport/channels/%d", request.ID), []byte(outgoing))
	if err != nil {
		return resp, status, err
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, status, err
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
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d49bcba99be73bff503ea6")
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
		fmt.Sprintf("/transport/channels/%s", strconv.FormatUint(id, 10)),
		buf,
	)
	if err != nil {
		return resp, status, err
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, status, err
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
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d49bcba99be73bff503ea6")
//	snd := SendData{
//		SendMessage{
//			Message{
//				ExternalID: "23e23e23",
//				Channel:    channelId,
//				Type:       "text",
//				Text:       "hello!",
//			},
//			time.Now(),
//		},
//		User{
//			ExternalID: "8",
//			Nickname:   "@octopulus",
//			Firstname:  "Joe",
//		},
//		channelId,
//	}
//
// 	data, status, err := client.Messages(snd)
//
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
//
//	fmt.Printf("%s\n", data.MessageID)
func (c *MgClient) Messages(request SendData) (MessagesResponse, int, error) {
	var resp MessagesResponse
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PostRequest("/transport/messages", []byte(outgoing))
	if err != nil {
		return resp, status, err
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, status, err
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
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d49bcba99be73bff503ea6")
//	snd := SendData{
//		SendMessage{
//			Message{
//				ExternalID: "23e23e23",
//				Channel:    channelId,
//				Type:       "text",
//				Text:       "hello!",
//			},
//			time.Now(),
//		},
//		User{
//			ExternalID: "8",
//			Nickname:   "@octopulus",
//			Firstname:  "Joe",
//		},
//		channelId,
//	}
//
// 	data, status, err := client.UpdateMessages(snd)
//
// 	if err != nil {
// 		fmt.Printf("%v", err)
// 	}
//
//	fmt.Printf("%s\n", data.MessageID)
func (c *MgClient) UpdateMessages(request UpdateData) (MessagesResponse, int, error) {
	var resp MessagesResponse
	outgoing, _ := json.Marshal(&request)

	data, status, err := c.PutRequest("/transport/messages", []byte(outgoing))
	if err != nil {
		return resp, status, err
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, status, err
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
// 	var client = v1.New("https://token.url", "cb8ccf05e38a47543ad8477d49bcba99be73bff503ea6")
//
// 	data, status, err := client.DeleteMessage("3053450384")
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
		"/transport/messages/",
		[]byte(outgoing),
	)
	if err != nil {
		return resp, status, err
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return resp, status, err
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
