package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateItem_MarshalJSON(t *testing.T) {
	text := TemplateItem{
		Type: TemplateItemTypeText,
		Text: "text item",
	}

	variable := TemplateItem{
		Type:    TemplateItemTypeVar,
		VarType: TemplateVarFirstName,
	}

	emptyVariable := TemplateItem{
		Type:    TemplateItemTypeVar,
		VarType: "",
	}

	data, err := json.Marshal(text)
	assert.NoError(t, err)
	assert.Equal(t, "\""+text.Text+"\"", string(data))

	data, err = json.Marshal(variable)
	assert.NoError(t, err)
	assert.Equal(t, `{"var":"first_name"}`, string(data))

	data, err = json.Marshal(emptyVariable)
	assert.NoError(t, err)
	assert.Equal(t, "{}", string(data))
}

func TestTemplateItem_UnmarshalJSON(t *testing.T) {
	var (
		textResult          TemplateItem
		variableResult      TemplateItem
		emptyVariableResult TemplateItem
	)

	text := []byte("\"text block\"")
	variable := []byte(`{"var":"first_name"}`)
	emptyVariable := []byte("{}")

	require.NoError(t, json.Unmarshal(text, &textResult))
	require.NoError(t, json.Unmarshal(variable, &variableResult))
	require.NoError(t, json.Unmarshal(emptyVariable, &emptyVariableResult))

	assert.Equal(t, TemplateItemTypeText, textResult.Type)
	assert.Equal(t, string(text)[1:11], textResult.Text)

	assert.Equal(t, TemplateItemTypeVar, variableResult.Type)
	assert.Equal(t, TemplateVarFirstName, variableResult.VarType)
	assert.Empty(t, variableResult.Text)

	assert.Equal(t, TemplateItemTypeVar, emptyVariableResult.Type)
	assert.Equal(t, TemplateVarCustom, emptyVariableResult.VarType)
	assert.Empty(t, emptyVariableResult.Text)
}

func TestUnmarshalInteractiveTemplate_TextHeader(t *testing.T) {
	var template Template
	input := `{
	"code":"aaa#bbb#ru",
    "phone": "79252223456",
    "channel_id": 1,
    "header": {
        "content": {
			"type": "text",
	        "body": "Hello, {{1}}!"
		}
    },
    "footer": "Scooter",
    "buttons": {
		"items": [
			{
	            "type": "url",
				"label": "Go to website",
	            "url": "222ddd"
	        },
	        {
	            "type": "plain",
	            "label": "Yes"
	        }
		]
	},
    "verification_status": "approved"
}`
	assert.NoError(t, json.Unmarshal([]byte(input), &template))

	assert.Equal(t, "aaa#bbb#ru", template.Code)
	assert.Equal(t, HeaderContentTypeText, template.Header.Content.HeaderContentType())

	h := template.Header.TextContent()
	assert.Equal(t, "Hello, {{1}}!", h.Body)
	assert.Equal(t, "Scooter", template.Footer)
	assert.Equal(t, TemplateStatusApproved, template.VerificationStatus)
	assert.Equal(t, ButtonTypeURL, template.Buttons.Items[0].ButtonType())
	assert.Equal(t, "222ddd", template.Buttons.Items[0].(*URLButton).URL)
	assert.Equal(t, "Go to website", template.Buttons.Items[0].(*URLButton).Label)
	assert.Equal(t, ButtonTypePlain, template.Buttons.Items[1].ButtonType())
	assert.Equal(t, "Yes", template.Buttons.Items[1].(*PlainButton).Label)

	input = `{"footer": "Scooter"}`
	template = Template{}
	assert.NoError(t, json.Unmarshal([]byte(input), &template))
	assert.Nil(t, template.Header)
	assert.Empty(t, template.Buttons)
}

func TestUnmarshalInteractiveTemplate_DocumentHeader(t *testing.T) {
	var template Template
	input := `{
	"code":"aaa#bbb#ru",
    "phone": "79252223456",
    "channel_id": 1,
    "header": {
        "content": {
			"type": "document"
		}
    },
    "footer": "Scooter",
    "buttons": {
		"items": [
			{
	            "type": "url",
				"label": "Go to website",
	            "url": "222ddd"
	        },
	        {
	            "type": "plain",
	            "label": "Yes"
	        }
		]
	},
    "verification_status": "approved"
}`
	assert.NoError(t, json.Unmarshal([]byte(input), &template))

	assert.Equal(t, "aaa#bbb#ru", template.Code)
	assert.Equal(t, HeaderContentTypeDocument, template.Header.Content.HeaderContentType())
	assert.NotNil(t, template.Header.DocumentContent())
	assert.Equal(t, "Scooter", template.Footer)
	assert.Equal(t, TemplateStatusApproved, template.VerificationStatus)
	assert.Equal(t, ButtonTypeURL, template.Buttons.Items[0].ButtonType())
	assert.Equal(t, "222ddd", template.Buttons.Items[0].(*URLButton).URL)
	assert.Equal(t, "Go to website", template.Buttons.Items[0].(*URLButton).Label)
	assert.Equal(t, ButtonTypePlain, template.Buttons.Items[1].ButtonType())
	assert.Equal(t, "Yes", template.Buttons.Items[1].(*PlainButton).Label)

	input = `{"footer": "Scooter"}`
	template = Template{}
	assert.NoError(t, json.Unmarshal([]byte(input), &template))
	assert.Nil(t, template.Header)
	assert.Empty(t, template.Buttons)
}

func TestUnmarshalInteractiveTemplate_ImageHeader(t *testing.T) {
	var template Template
	input := `{
	"code":"aaa#bbb#ru",
    "phone": "79252223456",
    "channel_id": 1,
    "header": {
        "content": {
			"type": "image"
		}
    },
    "footer": "Scooter",
    "buttons": {
		"items": [
			{
	            "type": "url",
				"label": "Go to website",
	            "url": "222ddd"
	        },
	        {
	            "type": "plain",
	            "label": "Yes"
	        }
		]
	},
    "verification_status": "approved"
}`
	assert.NoError(t, json.Unmarshal([]byte(input), &template))

	assert.Equal(t, "aaa#bbb#ru", template.Code)
	assert.Equal(t, HeaderContentTypeImage, template.Header.Content.HeaderContentType())
	assert.NotNil(t, template.Header.ImageContent())
	assert.Equal(t, "Scooter", template.Footer)
	assert.Equal(t, TemplateStatusApproved, template.VerificationStatus)
	assert.Equal(t, ButtonTypeURL, template.Buttons.Items[0].ButtonType())
	assert.Equal(t, "222ddd", template.Buttons.Items[0].(*URLButton).URL)
	assert.Equal(t, "Go to website", template.Buttons.Items[0].(*URLButton).Label)
	assert.Equal(t, ButtonTypePlain, template.Buttons.Items[1].ButtonType())
	assert.Equal(t, "Yes", template.Buttons.Items[1].(*PlainButton).Label)

	input = `{"footer": "Scooter"}`
	template = Template{}
	assert.NoError(t, json.Unmarshal([]byte(input), &template))
	assert.Nil(t, template.Header)
	assert.Empty(t, template.Buttons)
}

func TestUnmarshalInteractiveTemplate_VideoHeader(t *testing.T) {
	var template Template
	input := `{
	"code":"aaa#bbb#ru",
    "phone": "79252223456",
    "channel_id": 1,
    "header": {
        "content": {
			"type": "video"
		}
    },
    "footer": "Scooter",
    "buttons": {
		"items": [
			{
	            "type": "url",
				"label": "Go to website",
	            "url": "222ddd"
	        },
	        {
	            "type": "plain",
	            "label": "Yes"
	        }
		]
	},
    "verification_status": "approved"
}`
	assert.NoError(t, json.Unmarshal([]byte(input), &template))

	assert.Equal(t, "aaa#bbb#ru", template.Code)
	assert.Equal(t, HeaderContentTypeVideo, template.Header.Content.HeaderContentType())
	assert.NotNil(t, template.Header.VideoContent())
	assert.Equal(t, "Scooter", template.Footer)
	assert.Equal(t, TemplateStatusApproved, template.VerificationStatus)
	assert.Equal(t, ButtonTypeURL, template.Buttons.Items[0].ButtonType())
	assert.Equal(t, "222ddd", template.Buttons.Items[0].(*URLButton).URL)
	assert.Equal(t, "Go to website", template.Buttons.Items[0].(*URLButton).Label)
	assert.Equal(t, ButtonTypePlain, template.Buttons.Items[1].ButtonType())
	assert.Equal(t, "Yes", template.Buttons.Items[1].(*PlainButton).Label)

	input = `{"footer": "Scooter"}`
	template = Template{}
	assert.NoError(t, json.Unmarshal([]byte(input), &template))
	assert.Nil(t, template.Header)
	assert.Empty(t, template.Buttons)
}

func TestUnmarshalInteractiveTemplate_Examples(t *testing.T) {
	var template Template
	input := `{
	"code":"aaa#bbb#ru",
    "phone": "79252223456",
    "channel_id": 1,
    "header": {
        "content": {
			"type": "text",
			"body": "Hello, {{1}}!"
		}
    },
    "body": "Order {{1}} successfully delivered",
    "buttons": {
		"items": [
			{
	            "type": "url",
				"label": "Go to website",
	            "url": "https://test.com/{{1}}"
	        },
	        {
	            "type": "plain",
	            "label": "OK"
	        }
		]
	},
    "verification_status": "approved",
	"example": {
		"header": ["Alex"],
		"body": ["ORDER-111"],
		"buttons": [["id123"], []]
	}
}`
	assert.NoError(t, json.Unmarshal([]byte(input), &template))

	assert.NotNil(t, template.Example)
	assert.Equal(t, []string{"Alex"}, template.Example.Header)
	assert.Equal(t, []string{"ORDER-111"}, template.Example.Body)
	assert.Equal(t, [][]string{{"id123"}, {}}, template.Example.Buttons)
}

func TestUnmarshalInteractiveTemplate_Attachments(t *testing.T) {
	var template Template
	input := `{
	"code":"aaa#bbb#ru",
    "phone": "79252223456",
    "channel_id": 1,
    "header": {
        "content": {
			"type": "image"
		}
    },
    "body": "Welcome to new delivery point",
    "verification_status": "approved",
	"example": {
		"attachments": [{"id": "a6cf882e-6915-410a-8672-ed3a28a7875d", "caption": "test-cats.png"}]
	}
}`
	assert.NoError(t, json.Unmarshal([]byte(input), &template))

	assert.NotNil(t, template.Example)
	assert.Len(t, template.Example.Attachments, 1)
	assert.Equal(t, "a6cf882e-6915-410a-8672-ed3a28a7875d", template.Example.Attachments[0].ID)
	assert.Equal(t, "test-cats.png", template.Example.Attachments[0].Caption)
}
