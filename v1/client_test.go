package v1

import (
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	mgURL        = os.Getenv("MG_URL")
	mgToken      = os.Getenv("MG_TOKEN")
	channelId, _ = strconv.ParseUint(os.Getenv("MG_CHANNEL"), 10, 64)
	ext          = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
)

func client() *MgClient {
	return New(mgURL, mgToken)
}

func TestMgClient_ActivateTransportChannel(t *testing.T) {
	c := client()
	ch := Channel{
		ID:   channelId,
		Type: "telegram",
		Events: []string{
			"message_sent",
			"message_read",
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
		Events: []string{
			"message_sent",
			"message_updated",
			"message_deleted",
			"message_read",
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

	if deleteData.DectivatedAt.String() == "" {
		t.Errorf("%v", err)
	}

	t.Logf("Deactivate new channel with ID %v", deleteData.ChannelID)
}

func TestMgClient_UpdateTransportChannel(t *testing.T) {
	c := client()
	ch := Channel{
		ID: channelId,
		Events: []string{
			"message_sent",
			"message_updated",
			"message_deleted",
			"message_read",
		},
	}

	data, status, err := c.UpdateTransportChannel(ch)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	t.Logf("Update selected channel: %v", data.ChannelID)
}

func TestMgClient_Messages(t *testing.T) {
	c := client()
	t.Logf("%v", ext)

	snd := SendData{
		SendMessage{
			Message{
				ExternalID: ext,
				Type:       "text",
				Text:       "hello!",
			},
			time.Now(),
		},
		User{
			ExternalID: "6",
			Nickname:   "octopus",
			Firstname:  "Joe",
		},
		channelId,
		"24798237492374",
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

	snd := UpdateData{
		UpdateMessage{
			Message{
				ExternalID: ext,
				Type:       "text",
				Text:       "hello hello!",
			},
			time.Now(),
		},
		channelId,
	}

	data, status, err := c.UpdateMessages(snd)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Time.String() == "" {
		t.Errorf("%v", err)
	}

	t.Logf("Message %v updated", data.MessageID)
}

func TestMgClient_DeleteMessage(t *testing.T) {
	c := client()
	t.Logf("%v", ext)

	snd := DeleteData{
		Message{
			ExternalID: ext,
		},
		channelId,
	}

	data, status, err := c.DeleteMessage(snd)
	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Time.String() == "" {
		t.Errorf("%v", err)
	}

	t.Logf("Message %v updated", data.MessageID)
}

func TestMgClient_DeactivateTransportChannel(t *testing.T) {
	c := client()
	deleteData, status, err := c.DeactivateTransportChannel(channelId)

	if err != nil {
		t.Errorf("%d %v", status, err)
	}

	if deleteData.DectivatedAt.String() == "" {
		t.Errorf("%v", err)
	}

	t.Logf("Deactivate selected channel: %v", deleteData.ChannelID)
}
