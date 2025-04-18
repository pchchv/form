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

// SetAnonymousMode sets the mode the encoder should run Default is AnonymousEmbed.
func (e *Encoder) SetAnonymousMode(mode AnonymousMode) {
	e.embedAnonymous = mode == AnonymousEmbed
}

// SetNamespacePrefix sets a struct namespace prefix.
func (e *Encoder) SetNamespacePrefix(namespacePrefix string) {
	e.namespacePrefix = namespacePrefix
}

// SetNamespaceSuffix sets a struct namespace suffix.
func (e *Encoder) SetNamespaceSuffix(namespaceSuffix string) {
	e.namespaceSuffix = namespaceSuffix
}

// SetTagName sets the given tag name to be used by the encoder.
// Default is "form"
func (e *Encoder) SetTagName(tagName string) {
	e.tagName = tagName
}

// SetMode sets the mode the encoder should run.
// Default is ModeImplicit.
func (e *Encoder) SetMode(mode Mode) {
	e.mode = mode
}
