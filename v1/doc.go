// Package v1 provides Go API Client implementation for MessageGateway Transport API.
//
// You can use v1.New or v1.NewWithClient to initialize API client. github.com/retailcrm/mg-transport-api-client-go/examples
// package contains some examples on how to use this library properly.
//
// Basic usage example:
//
//	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
//	getReplyDeadline := func(after time.Duration) *time.Time {
//		deadline := time.Now().Add(after)
//		return &deadline
//	}
//
//	resp, status, err := client.Messages(SendData{
//		Message: Message{
//			ExternalID: "uid_1",
//			Type:       MsgTypeText,
//			Text:       "Hello customer!",
//			PageLink:   "https://example.com",
//		},
//		Originator: OriginatorCustomer,
//		Customer: Customer{
//			ExternalID: "client_id_1",
//			Nickname:   "customer",
//			Firstname:  "Tester",
//			Lastname:   "Tester",
//			Avatar:     "https://example.com/image.png",
//			ProfileURL: "https://example.com/user/client_id_1",
//			Language:   "en",
//			Utm: &Utm{
//				Source:   "myspace.com",
//				Medium:   "social",
//				Campaign: "something",
//				Term:     "fedora",
//				Content:  "autumn_collection",
//			},
//		},
//		Channel:        305,
//		ExternalChatID: "chat_id_1",
//		ReplyDeadline: getReplyDeadline(24 * time.Hour),
//	})
//	if err != nil {
//		log.Fatalf("request error: %s (%d)", err, status)
//	}
//
//	log.Printf("status: %d, message ID: %d", status, resp.MessageID)
package v1

import (
	"log"
	"time"
)

func ooga() {
	client := New("https://message-gateway.url", "cb8ccf05e38a47543ad8477d4999be73bff503ea6")
	getReplyDeadline := func(after time.Duration) *time.Time {
		deadline := time.Now().Add(after)
		return &deadline
	}
	resp, status, err := client.Messages(SendData{
		Message: Message{
			ExternalID: "uid_1",
			Type:       MsgTypeText,
			Text:       "Hello customer!",
			PageLink:   "https://example.com",
		},
		Originator: OriginatorCustomer,
		Customer: Customer{
			ExternalID: "client_id_1",
			Nickname:   "customer",
			Firstname:  "Tester",
			Lastname:   "Tester",
			Avatar:     "https://example.com/image.png",
			ProfileURL: "https://example.com/user/client_id_1",
			Language:   "en",
			Utm: &Utm{
				Source:   "myspace.com",
				Medium:   "social",
				Campaign: "something",
				Term:     "fedora",
				Content:  "autumn_collection",
			},
		},
		Channel:        305,
		ExternalChatID: "chat_id_1",
		ReplyDeadline:  getReplyDeadline(24 * time.Hour),
	})
	if err != nil {
		log.Fatalf("request error: %s (%d)", err, status)
	}

	log.Printf("status: %d, message ID: %d", status, resp.MessageID)
}
