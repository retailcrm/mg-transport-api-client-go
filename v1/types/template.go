package types

import (
	"encoding/json"
	"errors"
	"fmt"
)

// TemplateTypeText is a text template type. There is no other types for now.
const TemplateTypeText = "text"

const (
	// TemplateItemTypeText is a type for text chunk in template
	TemplateItemTypeText = iota
	// TemplateItemTypeVar is a type for variable in template
	TemplateItemTypeVar
)

const (
	// TemplateVarCustom is a custom variable type
	TemplateVarCustom = "custom"
	// TemplateVarName is a name variable type
	TemplateVarName = "name"
	// TemplateVarFirstName is a first name variable type
	TemplateVarFirstName = "first_name"
	// TemplateVarLastName is a last name variable type
	TemplateVarLastName = "last_name"
)

// templateVarAssoc for checking variable validity, only for internal use
var templateVarAssoc = map[string]interface{}{
	TemplateVarCustom:    nil,
	TemplateVarName:      nil,
	TemplateVarFirstName: nil,
	TemplateVarLastName:  nil,
}

// Template struct
type Template struct {
	Code      string         `json:"code"`
	ChannelID int64          `json:"channel_id,omitempty"`
	Name      string         `json:"name"`
	Enabled   bool           `json:"enabled,omitempty"`
	Type      string         `json:"type"`
	Template  []TemplateItem `json:"template"`
}

// TemplateItem is a part of template
type TemplateItem struct {
	Type    uint8
	Text    string
	VarType string
}

// MarshalJSON will correctly marshal TemplateItem
func (t *TemplateItem) MarshalJSON() ([]byte, error) {
	switch t.Type {
	case TemplateItemTypeText:
		return json.Marshal(t.Text)
	case TemplateItemTypeVar:
		return json.Marshal(map[string]interface{}{
			"var": t.VarType,
		})
	}

	return nil, errors.New("unknown TemplateItem type")
}

// UnmarshalJSON will correctly unmarshal TemplateItem
func (t *TemplateItem) UnmarshalJSON(b []byte) error {
	var obj interface{}
	err := json.Unmarshal(b, &obj)
	if err != nil {
		return err
	}

	switch v := obj.(type) {
	case string:
		t.Type = TemplateItemTypeText
		t.Text = v
	case map[string]interface{}:
		// {} case
		if len(v) == 0 {
			t.Type = TemplateItemTypeVar
			t.VarType = TemplateVarCustom
			return nil
		}

		if varTypeCurr, ok := v["var"].(string); ok {
			if _, ok := templateVarAssoc[t.VarType]; !ok {
				return fmt.Errorf("invalid placeholder var '%s'", varTypeCurr)
			}

			t.Type = TemplateItemTypeVar
		} else {
			return errors.New("invalid TemplateItem")
		}
	default:
		return errors.New("invalid TemplateItem")
	}

	return nil
}
