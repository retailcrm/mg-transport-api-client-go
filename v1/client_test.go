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
	channelID, _ = strconv.ParseUint(os.Getenv("MG_CHANNEL"), 10, 64)
	ext          = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
)

func client() *MgClient {
	c := New(mgURL, mgToken)
	c.Debug = true

	return c
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

func TestMgClient_TextMessages(t *testing.T) {
	c := client()
	t.Logf("%v", ext)

	snd := SendData{
		Message: Message{
			ExternalID: ext,
			Type:       MsgTypeText,
			Text:       "hello!",
		},
		Originator: "user",
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
		Url: "https://via.placeholder.com/300",
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
		Originator: "user",
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

	data, status, err := c.DeleteMessage(sndD)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	t.Logf("Message %v deleted", data.MessageID)

	sndD = DeleteData{
		Message{
			ExternalID: ext + "file",
		},
		channelID,
	}

	data, status, err = c.DeleteMessage(sndD)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	t.Logf("Message %v deleted", data.MessageID)
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

	resp, err := http.Get("https://via.placeholder.com/300")
	if err != nil {
		t.Errorf("%v", err)
	}

	defer resp.Body.Close()

	data, status, err := c.UploadFile(resp.Body)

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	t.Logf("Message %+v is sent", data)
}
