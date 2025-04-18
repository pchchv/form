package form

import (
	"bytes"
	"reflect"
	"strings"
)

// EncodeErrors is a map of errors encountered during form encoding.
type EncodeErrors map[string]error

func (e EncodeErrors) Error() string {
	buff := bytes.NewBufferString(blank)
	for k, err := range e {
		buff.WriteString(fieldNS)
		buff.WriteString(k)
		buff.WriteString(errorText)
		buff.WriteString(err.Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

// InvalidEncodeError describes an invalid argument passed to Encode.
type InvalidEncodeError struct {
	Type reflect.Type
}

func (e *InvalidEncodeError) Error() string {
	if e.Type == nil {
		return "form: Encode(nil)"
	}

	return "form: Encode(nil " + e.Type.String() + ")"
}
