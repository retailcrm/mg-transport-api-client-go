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
  "type": "message_sent",
  "meta": {
    "timestamp": 1703686050
  },
  "data": {
    "external_user_id": "79998887766",
    "external_chat_id": "",
    "channel_id": 83,
    "type": "image",
    "content": "Thank you for your order\n\nYou have placed order No. 8061C in the amount of 17400. We have already started working on it and will soon notify you of a change in status.\n\nStay with us",
    "quote_external_id": null,
    "quote_content": null,
    "in_app_id": 999,
    "user": {
      "id": 222,
      "first_name": "Alex",
      "last_name": "",
      "avatar": ""
    },
    "customer": {
      "first_name": "",
      "last_name": "",
      "avatar": ""
    },
    "items": [
      {
        "id": "aa4ff988-bafb-43a0-9551-b8e00f677e34",
        "size": 71984,
        "caption": "test.png",
        "height": 742,
        "width": 305
      }
    ],
    "template": {
      "code": "namespace#BABA_JABA#ru",
      "args": [
        "var0",
        "var1"
      ],
      "variables": {
        "header": {
          "type": "image",
          "attachments": [
            {
              "id": "aa4ff988-bafb-43a0-9551-b8e00f677e34",
              "caption": "test.png"
            }
          ]
        },
        "body": {
          "args": [
            "var0",
            "var1"
          ]
        },
        "buttons": [
          {
            "type": "plain",
            "title": "OK"
          },
          {
            "type": "url",
            "title": "Our site",
            "args": [
              "id0"
            ]
          },
          {
            "type": "phone",
            "title": "Our phone"
          }
        ]
      }
    }
  }
}`

	var wh struct {
		Data struct {
			Template TemplateInfo `json:"template"`
		} `json:"data"`
	}
	assert.NoError(t, json.Unmarshal([]byte(tmplJSON), &wh))

	tmpl := wh.Data.Template
	assert.Equal(t, "namespace#BABA_JABA#ru", tmpl.Code)

	assert.NotNil(t, tmpl.Variables.Header)
	assert.Empty(t, tmpl.Variables.Header.Args)
	assert.Equal(t, "aa4ff988-bafb-43a0-9551-b8e00f677e34", tmpl.Variables.Header.Attachments[0].ID)
	assert.Equal(t, "test.png", tmpl.Variables.Header.Attachments[0].Caption)

	assert.Equal(t, []string{"var0", "var1"}, tmpl.Variables.Body.Args)

	assert.Len(t, tmpl.Variables.Buttons, 3)
	assert.Equal(t, "plain", tmpl.Variables.Buttons[0].Type)
	assert.Equal(t, "OK", tmpl.Variables.Buttons[0].Title)
	assert.Empty(t, tmpl.Variables.Buttons[0].Args)
	assert.Equal(t, "url", tmpl.Variables.Buttons[1].Type)
	assert.Equal(t, "Our site", tmpl.Variables.Buttons[1].Title)
	assert.Equal(t, []string{"id0"}, tmpl.Variables.Buttons[1].Args)
	assert.Equal(t, "phone", tmpl.Variables.Buttons[2].Type)
	assert.Equal(t, "Our phone", tmpl.Variables.Buttons[2].Title)
	assert.Empty(t, tmpl.Variables.Buttons[2].Args)
}
