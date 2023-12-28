package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func InitTGBotAPI() {
	botAPI, err := tgbotapi.NewBotAPI(AppConfig.TGBotToken)
	if err != nil {
		log.Fatalln("cannot initialize TG Bot API:", err)
	}
	TG = botAPI
	log.Printf("initialized Telegram Bot API for bot @%s", botAPI.Self.UserName)
}

func SetTGWebhook() {
	wh, err := tgbotapi.NewWebhook(fmt.Sprintf("%s/api/v1/tg", AppConfig.BaseURL))
	if err != nil {
		log.Fatalln("cannot initialize webhook data:", err)
	}
	_, err = TG.Request(wh)
	if err != nil {
		log.Fatalln("cannot register TG webhook:", err)
	}
}
