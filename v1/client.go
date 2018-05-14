package v1

import (
	"encoding/json"
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
// 	var client = v1.New("https://demo.url", "09jIJ")
//
//  request := ActivateRequest{
//		Type: "telegram",
//		Events: [2]int{"message_sent", "message_sent"}
//  }
//
// 	data, status, err := client.ActivateTransportChannel(request)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	fmt.Printf("%s\n", data.CreatedAt)
func (c *MgClient) ActivateTransportChannel(request Channel) (ActivateResponse, int, Failure) {
	var resp ActivateResponse
	outgoing, _ := json.Marshal(&request)
	p := []byte(outgoing)

	data, status, err := c.PostRequest("/transports/channels", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if status != 200 && status != 201 {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// UpdateTransportChannel implement channel activation
//
// Example:
//
// 	var client = v1.New("https://demo.url", "09jIJ")
//
//  request := ActivateRequest{
//		ID:   3053450384,
//		Type: "telegram",
//		Events: [2]int{"message_sent", "message_sent"}
//  }
//
// 	data, status, err := client.UpdateTransportChannel(request)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	fmt.Printf("%s\n", data.UpdatedAt)
func (c *MgClient) UpdateTransportChannel(request Channel) (UpdateResponse, int, Failure) {
	var resp UpdateResponse
	var url = fmt.Sprintf("/transports/channels/%d", request.ID)
	outgoing, _ := json.Marshal(&request)
	p := []byte(outgoing)

	data, status, err := c.PutRequest(url, p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if status != 200 {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// DeactivateTransportChannel implement channel deactivation
//
// Example:
//
// 	var client = v1.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.DeactivateTransportChannel(3053450384)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	fmt.Printf("%s\n", data.DectivatedAt)
func (c *MgClient) DeactivateTransportChannel(id uint64) (DeleteResponse, int, Failure) {
	var resp DeleteResponse
	var url = fmt.Sprintf("/transports/channels/%s", strconv.FormatUint(id, 10))

	data, status, err := c.DeleteRequest(url)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if status != 200 {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Messages implement send message
//
// Example:
//
// 	var client = v1.New("https://demo.url", "09jIJ")
//	var message = Message{
//		ExternalID: "92784982374239847293",
//		Сhannel: "3053450384",
//		Type: "text",
//		Text: "Hello!",
//	}
//
//	var user = User{
//		ExternalID: "453535434535",
//		Nickname: "John_Doe",
//	}
//
// 	data, status, err := client.Messages(message, user)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	fmt.Printf("%s\n", data.MessageID)
func (c *MgClient) Messages(request SendData) (MessagesResponse, int, Failure) {
	var resp MessagesResponse
	outgoing, _ := json.Marshal(&request)
	p := []byte(outgoing)

	data, status, err := c.PostRequest("/messages", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if status != 200 {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// UpdateMessages implement edit message
//
// Example:
//
// 	var client = v1.New("https://demo.url", "09jIJ")
//	var message = Message{
//		ExternalID: "92784982374239847293",
//		Сhannel: "3053450384",
//		Type: "text",
//		Text: "Hello!",
//	}
//
//	var user = User{
//		ExternalID: "453535434535",
//		Nickname: "John_Doe",
//	}
//
// 	data, status, err := client.UpdateMessages(message, user)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	fmt.Printf("%s\n", data.MessageID)
/*func (c *MgClient) UpdateMessages(message UpdateMessage, user User) (MessagesResponse, int, Failure) {
	var resp MessagesResponse

	msg, _ := json.Marshal(&message)
	usr, _ := json.Marshal(&user)

	p := url.Values{
		"message": {string(msg[:])},
		"user":    {string(usr[:])},
	}

	data, status, err := c.PutRequest("/messages", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if status != 200 {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}*/

// DeleteMessage implement delete message
//
// Example:
//
// 	var client = v1.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.DeleteMessage("3053450384")
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	fmt.Printf("%s\n", data.MessageID)
func (c *MgClient) DeleteMessage(id string) (MessagesResponse, int, Failure) {
	var resp MessagesResponse

	data, status, err := c.DeleteRequest(fmt.Sprintf("/messages/%s", id))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if status != 200 {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}
