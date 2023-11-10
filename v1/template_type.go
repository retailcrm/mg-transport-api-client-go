package v1

import (
	"bytes"
	"errors"
)

type TemplateType uint8

const (
	TemplateTypeText TemplateType = iota + 1
	TemplateTypeMedia
)

var TypeMap = [][]byte{
	TemplateTypeText:  []byte("text"),
	TemplateTypeMedia: []byte("media"),
}

var ErrUnknownTypeValue = errors.New("unknown TemplateType")

func (e TemplateType) MarshalText() (text []byte, err error) {
	if e.isValid() {
		return TypeMap[e], nil
	}

	return nil, ErrUnknownTypeValue
}

func (e TemplateType) String() string {
	if e.isValid() {
		return string(TypeMap[e])
	}

	panic(ErrUnknownTypeValue)
}

func (e *TemplateType) UnmarshalText(text []byte) error {
	for f, v := range TypeMap {
		if !bytes.Equal(text, v) {
			continue
		}

		*e = TemplateType(f)
		return nil
	}

	return ErrUnknownTypeValue
}

func (e TemplateType) isValid() bool {
	return int(e) < len(TypeMap)
}
