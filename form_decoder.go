package form

import (
	"bytes"
	"reflect"
	"strings"
	"sync"
)

// DecodeErrors is a map of errors encountered during form decoding
type DecodeErrors map[string]error

func (d DecodeErrors) Error() string {
	buff := bytes.NewBufferString(blank)
	for k, err := range d {
		buff.WriteString(fieldNS)
		buff.WriteString(k)
		buff.WriteString(errorText)
		buff.WriteString(err.Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

// InvalidDecoderError describes an invalid argument passed to Decode.
// Argument passed to Decode must be a non-nil pointer.
type InvalidDecoderError struct {
	Type reflect.Type
}

func (e *InvalidDecoderError) Error() string {
	if e.Type == nil {
		return "form: Decode(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "form: Decode(non-pointer " + e.Type.String() + ")"
	}

	return "form: Decode(nil " + e.Type.String() + ")"
}

// DecodeCustomTypeFunc allows for registering/overriding types to be parsed.
type DecodeCustomTypeFunc func([]string) (interface{}, error)

// Decoder is the main decode instance
type Decoder struct {
	mode            Mode
	tagName         string
	dataPool        *sync.Pool
	structCache     *structCacheMap
	maxArraySize    int
	namespacePrefix string
	namespaceSuffix string
	customTypeFuncs map[reflect.Type]DecodeCustomTypeFunc
}
