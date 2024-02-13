package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	// TemplateItemTypeText is a type for text chunk in template.
	TemplateItemTypeText uint8 = iota
	// TemplateItemTypeVar is a type for variable in template.
	TemplateItemTypeVar
)

const (
	// TemplateVarCustom is a custom variable type.
	TemplateVarCustom = "custom"
	// TemplateVarName is a name variable type.
	TemplateVarName = "name"
	// TemplateVarFirstName is a first name variable type.
	TemplateVarFirstName = "first_name"
	// TemplateVarLastName is a last name variable type.
	TemplateVarLastName = "last_name"
)

// templateVarAssoc for checking variable validity, only for internal use.
var templateVarAssoc = map[string]interface{}{
	TemplateVarCustom:    nil,
	TemplateVarName:      nil,
	TemplateVarFirstName: nil,
	TemplateVarLastName:  nil,
}

// Template struct.
type Template struct {
	ID                 int64                      `json:"id,omitempty"`
	Code               string                     `json:"code,omitempty"`
	ChannelID          uint64                     `json:"channel_id"`
	Name               string                     `json:"name"`
	Enabled            bool                       `json:"enabled"`
	Type               TemplateType               `json:"type"`
	Template           []TemplateItem             `json:"template"`
	Body               string                     `json:"body"`
	Lang               string                     `json:"lang,omitempty"`
	Category           string                     `json:"category,omitempty"`
	Example            *TemplateExample           `json:"example,omitempty"`
	VerificationStatus TemplateVerificationStatus `json:"verification_status"`
	Quality            *TemplateQuality           `json:"quality,omitempty"`
	RejectionReason    TemplateRejectionReason    `json:"rejection_reason,omitempty"`
	Header             *TemplateHeader            `json:"header,omitempty"`
	Footer             string                     `json:"footer,omitempty"`
	Buttons            *TemplateButtons           `json:"buttons,omitempty"`
}

type TemplateExample struct {
	Body        []string                    `json:"body,omitempty"`
	Header      []string                    `json:"header,omitempty"`
	Buttons     [][]string                  `json:"buttons,omitempty"`
	Attachments []TemplateExampleAttachment `json:"attachments,omitempty"`
}

type TemplateButtons struct {
	Items []Button `json:"items"`
}

func (b *TemplateButtons) UnmarshalJSON(value []byte) error {
	var ScanType struct {
		Items []json.RawMessage `json:"items"`
	}

	if err := json.Unmarshal(value, &ScanType); err != nil {
		return err
	}

	var ButtonType struct {
		Type ButtonType `json:"type"`
	}

	for _, r := range ScanType.Items {
		if err := json.Unmarshal(r, &ButtonType); err != nil {
			return err
		}

		var btn Button
		switch ButtonType.Type {
		case ButtonTypePlain:
			btn = &PlainButton{}
		case ButtonTypePhone:
			btn = &PhoneButton{}
		case ButtonTypeURL:
			btn = &URLButton{}
		default:
			return errors.New("undefined type of button")
		}

		if err := json.Unmarshal(r, btn); err != nil {
			return err
		}

		b.Items = append(b.Items, btn)
	}

	return nil
}

func (b TemplateButtons) MarshalJSON() ([]byte, error) {
	var ValueType struct {
		Items []json.RawMessage `json:"items"`
	}

	for _, btn := range b.Items {
		btnData, err := json.Marshal(btn)
		if err != nil {
			return nil, err
		}

		buffer := bytes.NewBuffer(btnData[:len(btnData)-1])
		buffer.WriteByte(',')
		buffer.WriteString(fmt.Sprintf(`"type":"%s"`, btn.ButtonType()))
		buffer.WriteByte('}')

		ValueType.Items = append(ValueType.Items, buffer.Bytes())
	}

	d, err := json.Marshal(ValueType)
	if err != nil {
		return nil, err
	}

	return d, nil
}

type Button interface {
	ButtonType() ButtonType
}

type ButtonType string

const (
	ButtonTypePlain ButtonType = "plain"
	ButtonTypePhone ButtonType = "phone"
	ButtonTypeURL   ButtonType = "url"
)

type PlainButton struct {
	Label string `json:"label"`
}

func (PlainButton) ButtonType() ButtonType { return ButtonTypePlain }

type PhoneButton struct {
	Label string `json:"label"`
	Phone string `json:"phone"`
}

func (PhoneButton) ButtonType() ButtonType { return ButtonTypePhone }

type URLButton struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

func (URLButton) ButtonType() ButtonType { return ButtonTypeURL }

type HeaderContent interface {
	HeaderContentType() HeaderContentType
}

type HeaderContentType string

const (
	HeaderContentTypeText     HeaderContentType = "text"
	HeaderContentTypeDocument HeaderContentType = "document"
	HeaderContentTypeImage    HeaderContentType = "image"
	HeaderContentTypeVideo    HeaderContentType = "video"
)

type HeaderContentText struct {
	Body string `json:"body"`
}

func (HeaderContentText) HeaderContentType() HeaderContentType { return HeaderContentTypeText }

type HeaderContentDocument struct{}

func (HeaderContentDocument) HeaderContentType() HeaderContentType { return HeaderContentTypeDocument }

type HeaderContentImage struct{}

func (HeaderContentImage) HeaderContentType() HeaderContentType { return HeaderContentTypeImage }

type HeaderContentVideo struct{}

func (HeaderContentVideo) HeaderContentType() HeaderContentType { return HeaderContentTypeVideo }

type TemplateHeader struct {
	Content HeaderContent `json:"content"`
}

func (h *TemplateHeader) TextContent() *HeaderContentText {
	if h.Content.HeaderContentType() != HeaderContentTypeText {
		return nil
	}
	return h.Content.(*HeaderContentText)
}

func (h *TemplateHeader) DocumentContent() *HeaderContentDocument {
	if h.Content.HeaderContentType() != HeaderContentTypeDocument {
		return nil
	}
	return h.Content.(*HeaderContentDocument)
}

func (h *TemplateHeader) ImageContent() *HeaderContentImage {
	if h.Content.HeaderContentType() != HeaderContentTypeImage {
		return nil
	}
	return h.Content.(*HeaderContentImage)
}

func (h *TemplateHeader) VideoContent() *HeaderContentVideo {
	if h.Content.HeaderContentType() != HeaderContentTypeVideo {
		return nil
	}
	return h.Content.(*HeaderContentVideo)
}

func (h *TemplateHeader) UnmarshalJSON(value []byte) error {
	var ScanType struct {
		Content json.RawMessage `json:"content"`
	}

	if err := json.Unmarshal(value, &ScanType); err != nil {
		return err
	}

	var ContentType struct {
		Type HeaderContentType `json:"type"`
	}

	if err := json.Unmarshal(ScanType.Content, &ContentType); err != nil {
		return err
	}

	var c HeaderContent
	switch ContentType.Type {
	case HeaderContentTypeText:
		c = &HeaderContentText{}
	case HeaderContentTypeDocument:
		c = &HeaderContentDocument{}
	case HeaderContentTypeImage:
		c = &HeaderContentImage{}
	case HeaderContentTypeVideo:
		c = &HeaderContentVideo{}
	default:
		return errors.New("undefined type of header content")
	}

	if err := json.Unmarshal(ScanType.Content, c); err != nil {
		return err
	}

	h.Content = c
	return nil
}

func (h TemplateHeader) MarshalJSON() ([]byte, error) {
	content, err := json.Marshal(h.Content)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(content[:len(content)-1])
	if buffer.Len() > 1 {
		buffer.WriteByte(',')
	}
	buffer.WriteString(fmt.Sprintf(`"type":"%s"`, h.Content.HeaderContentType()))
	buffer.WriteByte('}')

	var ValueType struct {
		Content json.RawMessage `json:"content"`
	}

	ValueType.Content = buffer.Bytes()

	d, err := json.Marshal(ValueType)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// TemplateItem is a part of template.
type TemplateItem struct {
	Type    uint8
	Text    string
	VarType string
}

// MarshalJSON controls how TemplateItem will be marshaled into JSON.
func (t TemplateItem) MarshalJSON() ([]byte, error) {
	switch t.Type {
	case TemplateItemTypeText:
		return json.Marshal(t.Text)
	case TemplateItemTypeVar:
		// {} case, fast output without marshaling
		if t.VarType == "" || t.VarType == TemplateVarCustom {
			return []byte("{}"), nil
		}

		return json.Marshal(map[string]interface{}{
			"var": t.VarType,
		})
	}

	return nil, errors.New("unknown TemplateItem type")
}

// UnmarshalJSON will correctly unmarshal TemplateItem.
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
			if _, ok := templateVarAssoc[varTypeCurr]; !ok {
				return fmt.Errorf("invalid placeholder var '%s'", varTypeCurr)
			}

			t.Type = TemplateItemTypeVar
			t.VarType = varTypeCurr
		} else {
			return errors.New("invalid TemplateItem")
		}
	default:
		return errors.New("invalid TemplateItem")
	}

	return nil
}

type TemplateExampleAttachment struct {
	ID      string `json:"id"`
	Caption string `json:"caption"`
}
