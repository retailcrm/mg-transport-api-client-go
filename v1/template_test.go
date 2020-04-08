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
