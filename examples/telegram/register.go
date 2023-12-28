package main

import (
	retailcrm "github.com/retailcrm/api-client-go/v2"
	v1 "github.com/retailcrm/mg-transport-api-client-go/v1"
	"log"
)

func RegisterSystem() {
	client := retailcrm.New(AppConfig.System, AppConfig.APIKey)
	resp, _, err := client.IntegrationModuleEdit(retailcrm.IntegrationModule{
		Code:            "telegram-bot-integration-example",
		IntegrationCode: "telegram-bot-integration-example",
		Active:          v1.BoolPtr(true),
		Name:            "Telegram Bot Integration Example",
		ClientID:        "telegram-bot-integration-example",
		BaseURL:         AppConfig.BaseURL,
		AccountURL:      AppConfig.BaseURL,
		Integrations: &retailcrm.Integrations{
			MgTransport: &retailcrm.MgTransport{
				WebhookURL: AppConfig.BaseURL + "/api/v1/webhook",
			},
		},
	})
	if err != nil {
		log.Fatalln("cannot edit integration module:", err)
	}
	MG = v1.New(resp.Info.MgTransportInfo.EndpointURL, resp.Info.MgTransportInfo.Token)
	log.Println("updated integration module")
}

func RegisterChannel() {
	channels, _, err := MG.TransportChannels(v1.Channels{})
	if err != nil {
		log.Fatalln("cannot get channels:", err)
	}
	channel := v1.Channel{
		Type:     "telegram",
		Name:     "@" + TG.Self.UserName,
		Settings: getChannelSettings(),
	}
	for _, ch := range channels {
		if ch.Name != nil && *ch.Name == "@"+TG.Self.UserName {
			channel.ID = ch.ID
			_, _, err := MG.UpdateTransportChannel(channel)
			if err != nil {
				log.Fatalln("cannot update channel:", err)
			}
			Channel = &channel
			log.Println("updated MG channel")
			return
		}
	}
	resp, _, err := MG.ActivateTransportChannel(channel)
	if err != nil {
		log.Fatalln("cannot activate channel:", err)
	}
	log.Println("activated MG channel with id:", resp.ChannelID)
	channel.ID = resp.ChannelID
	Channel = &channel
}

func getChannelSettings() v1.ChannelSettings {
	return v1.ChannelSettings{
		Status: v1.Status{
			Delivered: v1.ChannelFeatureNone,
			Read:      v1.ChannelFeatureNone,
		},
		Text: v1.ChannelSettingsText{
			Creating:      v1.ChannelFeatureBoth,
			Editing:       v1.ChannelFeatureNone,
			Quoting:       v1.ChannelFeatureNone,
			Deleting:      v1.ChannelFeatureNone,
			MaxCharsCount: 2000,
		},
		Product: v1.Product{
			Creating: v1.ChannelFeatureNone,
			Editing:  v1.ChannelFeatureNone,
			Deleting: v1.ChannelFeatureNone,
		},
		Order: v1.Order{
			Creating: v1.ChannelFeatureNone,
			Editing:  v1.ChannelFeatureNone,
			Deleting: v1.ChannelFeatureNone,
		},
		File: v1.ChannelSettingsFilesBase{
			Creating: v1.ChannelFeatureNone,
			Editing:  v1.ChannelFeatureNone,
			Quoting:  v1.ChannelFeatureNone,
			Deleting: v1.ChannelFeatureNone,
		},
		Image: v1.ChannelSettingsFilesBase{
			Creating: v1.ChannelFeatureNone,
			Editing:  v1.ChannelFeatureNone,
			Quoting:  v1.ChannelFeatureNone,
			Deleting: v1.ChannelFeatureNone,
		},
		Audio: v1.ChannelSettingsAudio{
			Creating: v1.ChannelFeatureNone,
			Quoting:  v1.ChannelFeatureNone,
			Deleting: v1.ChannelFeatureNone,
		},
	}
}
