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

// SetMode sets the mode the decoder should run Default is ModeImplicit
func (d *Decoder) SetMode(mode Mode) {
	d.mode = mode
}

// SetTagName sets the given tag name to be used by the decoder.
//
// Default is "form".
func (d *Decoder) SetTagName(tagName string) {
	d.tagName = tagName
}

// SetMaxArraySize sets maximum array size that can be created.
// This limit is for the array indexing this library supports to
// avoid potential DOS or man-in-the-middle attacks using an unusually high number.
//
// Default is 10000.
func (d *Decoder) SetMaxArraySize(size uint) {
	d.maxArraySize = int(size)
}

// SetNamespacePrefix sets a struct namespace prefix.
func (d *Decoder) SetNamespacePrefix(namespacePrefix string) {
	d.namespacePrefix = namespacePrefix
}

// SetNamespaceSuffix sets a struct namespace suffix.
func (d *Decoder) SetNamespaceSuffix(namespaceSuffix string) {
	d.namespaceSuffix = namespaceSuffix
}
