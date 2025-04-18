package form

import (
	"bytes"
	"reflect"
	"strings"
	"sync"
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

// EncodeCustomTypeFunc allows for registering/overriding types to be parsed.
type EncodeCustomTypeFunc func(x interface{}) ([]string, error)

// Encoder is the main encode instance.
type Encoder struct {
	mode            Mode
	tagName         string
	dataPool        *sync.Pool
	structCache     *structCacheMap
	embedAnonymous  bool
	namespacePrefix string
	namespaceSuffix string
	customTypeFuncs map[reflect.Type]EncodeCustomTypeFunc
}
