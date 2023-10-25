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

var UnknownTypeValue = errors.New("unknown TemplateType")

func (e TemplateType) MarshalText() (text []byte, err error) {
	if e.isValid() {
		return TypeMap[e], nil
	}

	return nil, UnknownTypeValue
}

func (e TemplateType) String() string {
	if e.isValid() {
		return string(TypeMap[e])
	}

	panic(UnknownTypeValue)
}

func (e *TemplateType) UnmarshalText(text []byte) error {
	for f, v := range TypeMap {
		if !bytes.Equal(text, v) {
			continue
		}

		*e = TemplateType(f)
		return nil
	}

	return UnknownTypeValue
}

func (e TemplateType) isValid() bool {
	return int(e) < len(TypeMap)
}
