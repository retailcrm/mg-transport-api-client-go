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
}
