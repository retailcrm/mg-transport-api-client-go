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
			"message_updated",
			"message_deleted",
			"message_read",
		},
	}

	_, status, err := c.ActivateTransportChannel(ch)

	if status > http.StatusCreated {
		t.Errorf("%v", err)
	}
}

func TestMgClient_UpdateTransportChannel(t *testing.T) {
	c := client()
	ch := Channel{
		ID: channelId,
		Events: []string{
			"message_sent",
			"message_read",
		},
	}

	_, status, err := c.UpdateTransportChannel(ch)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}
}

func TestMgClient_DeactivateTransportChannel(t *testing.T) {
	c := client()
	deleteData, status, err := c.DeactivateTransportChannel(channelId)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if deleteData.DectivatedAt.String() == "" {
		t.Errorf("%v", err)
	}
}

func TestMgClient_Messages(t *testing.T) {
	c := client()

	snd := SendData{
		SendMessage{
			ExternalID: "qwerty",
			Channel:    channelId,
			Type:       "text",
			Text:       "hello!",
			SentAt:     time.Now(),
		},
		User{
			ExternalID: "1927391739217391",
			Nickname:   "@octopulus",
			Firstname:  "Alex",
		},
	}

	data, status, err := c.Messages(snd)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Time.String() == "" {
		t.Errorf("%v", err)
	}
}
