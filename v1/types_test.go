package v1

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendData_MarshalJSON(t *testing.T) {
	t.Run("WithoutOriginator", func(t *testing.T) {
		data, err := json.Marshal(SendData{})

		assert.NoError(t, err)

		expected := `{
			"channel": 0,
			"customer": {
				"external_id": "",
				"nickname": ""
			},
			"external_chat_id": "",
			"message": {
				"external_id": ""
			}
		}`
		assert.JSONEq(t, expected, string(data))
	})

	t.Run("WithOriginator", func(t *testing.T) {
		cases := []struct {
			originator Originator
			string     string
		}{
			{
				OriginatorCustomer,
				"customer",
			},
			{
				OriginatorChannel,
				"channel",
			},
		}

		pattern := `{
			"originator": "%s",
			"channel": 0,
			"customer": {
				"external_id": "",
				"nickname": ""
			},
			"external_chat_id": "",
			"message": {
				"external_id": ""
			}
		}`

		for _, c := range cases {
			data, err := json.Marshal(SendData{
				Originator: c.originator,
			})

			assert.NoError(t, err)
			assert.JSONEq(t, fmt.Sprintf(pattern, c.string), string(data))
		}
	})
}

func TestOriginator(t *testing.T) {
	b, err := OriginatorCustomer.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "customer", string(b))

	b, err = OriginatorChannel.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "channel", string(b))

	t.Run("MarshalText_Valid", func(t *testing.T) {
		cases := []struct {
			byte     byte
			expected string
		}{
			{1, "customer"},
			{2, "channel"},
		}
		for _, c := range cases {
			result, err := Originator(c.byte).MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, c.expected, string(result))
		}
	})

	t.Run("MarshalText_Invalid", func(t *testing.T) {
		data, err := Originator(0).MarshalText()
		assert.Nil(t, data)
		assert.Equal(t, err, ErrInvalidOriginator)
	})

	t.Run("MarshalJSON_Invalid", OriginatorMarshalJSONInvalid)

	t.Run("UnmarshalText_Valid", func(t *testing.T) {
		cases := []struct {
			value Originator
			text  []byte
		}{
			{1, []byte("customer")},
			{2, []byte("channel")},
		}
		for _, c := range cases {
			var o Originator
			err := o.UnmarshalText(c.text)
			assert.NoError(t, err)
			assert.Equal(t, c.value, o)
		}
	})

	t.Run("UnmarshalText_Invalid", func(t *testing.T) {
		var o Originator
		err := o.UnmarshalText([]byte{})
		assert.Empty(t, o)
		assert.Equal(t, err, ErrInvalidOriginator)
	})

	t.Run("UnmarshalJSON_Invalid", func(t *testing.T) {
		var o Originator
		err := json.Unmarshal([]byte("\"unknown\""), &o)
		assert.Empty(t, o)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidOriginator, err)
	})
}

func TestTransportErrorResponse(t *testing.T) {
	t.Run("NewTransportErrorResponse", func(t *testing.T) {
		expected := TransportResponse{
			Error: &TransportError{
				Code:    MessageErrorGeneral,
				Message: "error",
			},
		}
		resp := NewTransportErrorResponse(MessageErrorGeneral, "error")
		assert.Equal(t, expected, resp)
	})

	t.Run("NewSentMessageResponse", func(t *testing.T) {
		expected := TransportResponse{
			ExternalMessageID: "extID",
		}
		resp := NewSentMessageResponse("extID")
		assert.Equal(t, expected, resp)
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		cases := []struct {
			byte     byte
			expected string
		}{
			{1, "general"},
			{2, "customer_not_exists"},
			{3, "reply_timed_out"},
			{4, "spam_suspicion"},
			{5, "access_restricted"},
		}
		for _, c := range cases {
			result, err := TransportErrorCode(c.expected).MarshalJSON()
			expected := []byte(fmt.Sprintf(`"%s"`, c.expected))
			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		}
	})
}

func TestTemplateInfoUnmarshal(t *testing.T) {
	tmplJSON := `{
		"code": "namespace#BABA_JABA#ru",
		"variables": {
			"header": [
				"header1",
				"header2"
			],
			"attachments": [
				{"caption":"test-caption", "id":"550e8400-e29b-41d4-a716-446655440000"}
			],
			"body": [
				"BABA",
				"JABA"
			],
			"buttons": [
				["button1"],
				[],
				["button2"]
			]
		}
	}`

	var tmpl TemplateInfo
	assert.NoError(t, json.Unmarshal([]byte(tmplJSON), &tmpl))
	assert.Equal(t, "namespace#BABA_JABA#ru", tmpl.Code)
	assert.Equal(t, []string{"header1", "header2"}, tmpl.Variables.Header)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", tmpl.Variables.Attachments[0].ID)
	assert.Equal(t, "test-caption", tmpl.Variables.Attachments[0].Caption)
	assert.Equal(t, []string{"BABA", "JABA"}, tmpl.Variables.Body)
	assert.Equal(t, [][]string{{"button1"}, {}, {"button2"}}, tmpl.Variables.Buttons)
}

func TestUnmarshalMessageWebhook(t *testing.T) {
	msgJSON := `{
	  "type": "message_sent",
	  "meta": {
		"timestamp": 1703523308
	  },
	  "data": {
		"external_user_id": "79998887766",
		"external_chat_id": "",
		"channel_id": 83,
		"type": "text",
		"content": "Thank you for your order\n\nYou have placed order No. 8061C in the amount of 17400. We have already started working on it and will soon notify you of a change in status.\n\nStay with us",
		"quote_external_id": null,
		"quote_content": null,
		"in_app_id": 638,
		"bot": {
		  "id": 760,
		  "name": "Bot system",
		  "avatar": ""
		},
		"customer": {
		  "first_name": "",
		  "last_name": "",
		  "avatar": ""
		},
		"template": {
		  "code": "f87e678f_660b_461a_b60a_a6194e2e9094#thanks_for_order#ru",
		  "args": [
			"8061C",
			"17400"
		  ],
		  "variables": {
			"body": [
			  "8061C",
			  "17400"
			],
			"buttons": [
			  [],
			  [
				"8061"
			  ]
			]
		  }
		},
		"attachments": {
		  "suggestions": [
			{
			  "type": "text",
			  "title": "ОК"
			},
			{
			  "type": "url",
			  "title": "Our site",
			  "payload": "https://site.com/8061"
			}
		  ]
		}
	  }
	}`

	var wh struct {
		Data MessageWebhookData `json:"data"`
	}
	assert.NoError(t, json.Unmarshal([]byte(msgJSON), &wh))

	assert.NotNil(t, wh.Data.Attachments)
	assert.Len(t, wh.Data.Attachments.Suggestions, 2)

	okButton := wh.Data.Attachments.Suggestions[0]
	assert.Equal(t, SuggestionTypeText, okButton.Type)
	assert.Equal(t, "ОК", okButton.Title)
	assert.Empty(t, okButton.Payload)

	urlButton := wh.Data.Attachments.Suggestions[1]
	assert.Equal(t, SuggestionTypeURL, urlButton.Type)
	assert.Equal(t, "Our site", urlButton.Title)
	assert.Equal(t, "https://site.com/8061", urlButton.Payload)
}
