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

func TestUnmarshalMediaInteractiveTemplate(t *testing.T) {
	var template Template
	input := `{
	"code":"aaa#bbb#ru",
    "phone": "79252223456",
    "channel_id": 1,
    "headerParams": {
        "textVars": [
            "Johny",
            "1234C"
        ],
        "imageUrl": "http://example.com/intaro/d2354125",
        "videoUrl": "http://example.com/intaro/d2222",
        "documentUrl": "http://example.com/intaro/d4444"
    },
    "footer": "Scooter",
    "buttonParams": [
        {
            "type": "URL",
            "urlParameter": "222ddd"
        },
        {
            "type": "QUICK_REPLY",
            "text": "Yes"
        }
    ]
}`
	assert.NoError(t, json.Unmarshal([]byte(input), &template))

	assert.Equal(t, "aaa#bbb#ru", template.Code)
	assert.Equal(t, []string{"Johny", "1234C"}, template.HeaderParams.TextVars)
	assert.Equal(t, "http://example.com/intaro/d2354125", template.HeaderParams.ImageURL)
	assert.Equal(t, "http://example.com/intaro/d2222", template.HeaderParams.VideoURL)
	assert.Equal(t, "http://example.com/intaro/d4444", template.HeaderParams.DocumentURL)
	assert.Equal(t, "Scooter", *template.Footer)
	assert.Equal(t, URLButton, template.ButtonParams[0].ButtonType)
	assert.Equal(t, "222ddd", template.ButtonParams[0].URLParameter)
	assert.Equal(t, QuickReplyButton, template.ButtonParams[1].ButtonType)
	assert.Equal(t, "Yes", template.ButtonParams[1].Text)

	input = `{"footer": "Scooter"}`
	template = Template{}
	assert.NoError(t, json.Unmarshal([]byte(input), &template))
	assert.Nil(t, template.HeaderParams)
	assert.Empty(t, template.ButtonParams)
}
