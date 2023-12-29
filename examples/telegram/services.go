package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v1 "github.com/retailcrm/mg-transport-api-client-go/v1"
)

var (
	Channel *v1.Channel
	MG      *v1.MgClient
	TG      *tgbotapi.BotAPI
)
