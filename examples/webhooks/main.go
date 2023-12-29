package main

import (
	"encoding/json"
	"errors"
	"fmt"
	v1 "github.com/retailcrm/mg-transport-api-client-go/v1"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// H is a basic hashmap type.
type H map[string]interface{}

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Println("listening on", addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/webhook", HandleWebhook)
	err := http.ListenAndServe(addr, mux)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen %s err: %s", addr, err)
	}
}

func HandleWebhook(rw http.ResponseWriter, req *http.Request) {
	// You should authenticate the request. Usually it's done by some sort of token header.
	var wh v1.WebhookRequest
	if err := readJSON(req, &wh); err != nil {
		serveError(rw, http.StatusBadRequest, err)
		return
	}

	switch wh.Type {
	case v1.MessageSendWebhookType:
		HandleSendWebhook(rw, wh.MessageWebhookData())
	case v1.MessageReadWebhookType:
		HandleReadWebhook(rw, wh.MessageWebhookData())
	case v1.MessageDeleteWebhookType:
		HandleDeleteWebhook(rw, wh.MessageWebhookData())
	case v1.TemplateCreateWebhookType:
		HandleTemplateCreate(rw, wh.TemplateCreateWebhookData())
	case v1.TemplateUpdateWebhookType:
		HandleTemplateUpdate(rw, wh.TemplateUpdateWebhookData())
	case v1.TemplateDeleteWebhookType:
		HandleTemplateDelete(rw, wh.TemplateDeleteWebhookData())
	default:
		serveError(rw, http.StatusUnprocessableEntity, fmt.Errorf("unknown webhook type: %s", wh.Type))
	}
}

func HandleSendWebhook(rw http.ResponseWriter, msg v1.MessageWebhookData) {
	log.Printf("incoming message: %#v", msg)
	serveJSON(rw, http.StatusOK, v1.NewSentMessageResponse(strconv.FormatInt(time.Now().UnixNano(), 10)))
}

func HandleReadWebhook(rw http.ResponseWriter, msg v1.MessageWebhookData) {
	log.Printf("incoming message read status: %#v", msg)
	serveJSON(rw, http.StatusOK, H{})
}

func HandleDeleteWebhook(rw http.ResponseWriter, msg v1.MessageWebhookData) {
	log.Printf("incoming message removal: %#v", msg)
	serveJSON(rw, http.StatusOK, H{})
}

func HandleTemplateCreate(rw http.ResponseWriter, tpl v1.TemplateCreateWebhookData) {
	log.Printf("new template: %#v", tpl)
	serveJSON(rw, http.StatusOK, H{})
}

func HandleTemplateUpdate(rw http.ResponseWriter, tpl v1.TemplateUpdateWebhookData) {
	log.Printf("updated template: %#v", tpl)
	serveJSON(rw, http.StatusOK, H{})
}

func HandleTemplateDelete(rw http.ResponseWriter, tpl v1.TemplateDeleteWebhookData) {
	log.Printf("template removal: %#v", tpl)
	serveJSON(rw, http.StatusOK, H{})
}

func readJSON(req *http.Request, out any) error {
	defer func() { _ = req.Body.Close() }()
	return json.NewDecoder(req.Body).Decode(out)
}

func serveJSON(rw http.ResponseWriter, st int, data any) {
	resp, _ := json.Marshal(data)
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(st)
	_, _ = fmt.Fprintln(rw, string(resp))
}

func serveError(rw http.ResponseWriter, st int, err error) {
	serveJSON(rw, st, H{"error": err.Error()})
}
