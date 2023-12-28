package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	v1 "github.com/retailcrm/mg-transport-api-client-go/v1"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const ChatIDPrefix = "tg_"

func Listen() {
	router := gin.Default()
	router.POST("/api/v1/mg", MGWebhookHandler)
	router.POST("/api/v1/tg", TGWebhookHandler)

	srv := &http.Server{
		Addr:    AppConfig.Listen,
		Handler: router,
	}
	go func() {
		log.Printf("listening on %s", AppConfig.Listen)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("shutting down:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("quitting...")
	}
}

func MGWebhookHandler(c *gin.Context) {
	var wh v1.WebhookRequest
	if err := c.ShouldBindJSON(&wh); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			v1.NewTransportErrorResponse(v1.MessageErrorGeneral, "invalid webhook data"))
		return
	}

	if wh.Type != v1.MessageSendWebhookType {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity,
			v1.NewTransportErrorResponse(v1.MessageErrorGeneral, "unsupported webhook type"))
		return
	}

	whMsg := wh.MessageWebhookData()
	if strings.HasPrefix(whMsg.ExternalChatID, ChatIDPrefix) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity,
			v1.NewTransportErrorResponse(v1.MessageErrorGeneral, "unexpected chat ID"))
		return
	}

	chatID, err := strconv.ParseInt(whMsg.ExternalChatID[len(ChatIDPrefix):], 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity,
			v1.NewTransportErrorResponse(v1.MessageErrorGeneral, "unparsable chat ID"))
		return
	}

	resp, err := TG.Send(tgbotapi.NewMessage(chatID, whMsg.Content))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity,
			v1.NewTransportErrorResponse(v1.MessageErrorGeneral, err.Error()))
		return
	}
	c.JSON(http.StatusOK, v1.NewSentMessageResponse(strconv.Itoa(resp.MessageID)))
}

func TGWebhookHandler(c *gin.Context) {
	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if update.Message == nil {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	_, _, err := MG.Messages(v1.SendData{
		Message: v1.Message{
			ExternalID: newExternalMessageID(update.Message.MessageID),
			Type:       v1.MsgTypeText,
			Text:       update.Message.Text,
		},
		Originator: v1.OriginatorCustomer,
		Customer: v1.Customer{
			ExternalID: strconv.FormatInt(update.Message.From.ID, 10),
			Nickname:   update.Message.From.UserName,
			Firstname:  update.Message.From.FirstName,
			Lastname:   update.Message.From.LastName,
			ProfileURL: fmt.Sprintf("https://t.me/%s", update.Message.From.UserName),
			Language:   update.Message.From.LanguageCode,
		},
		Channel:        Channel.ID,
		ExternalChatID: ChatIDPrefix + strconv.FormatInt(update.Message.Chat.ID, 10),
	})
	if err != nil {
		log.Printf("error: cannot send message: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "cannot send message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func newExternalMessageID(messageID int) string {
	return fmt.Sprintf("%d-%d", TG.Self.ID, messageID)
}
