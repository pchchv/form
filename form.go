package form

import (
	"reflect"
	"time"
)

const (
	blank             = ""
	ignore            = "-"
	fieldNS           = "Field Namespace:"
	errorText         = " ERROR:"
	ModeImplicit Mode = iota // ModeImplicit tries to parse values for all fields that do not have an ignore '-' tag.
	ModeExplicit             // ModeExplicit parses values for field with a field tag and that tag is not the ignore '-' tag.
	// AnonymousEmbed embeds anonymous data when encoding
	// eg. type A struct { Field string }
	//     type B struct { A, Field string }
	//     encode results: url.Values{"Field":[]string{"B FieldVal", "A FieldVal"}}
	AnonymousEmbed AnonymousMode = iota

	// AnonymousSeparate does not embed anonymous data when encoding
	// eg. type A struct { Field string }
	//     type B struct { A, Field string }
	//     encode results: url.Values{"Field":[]string{"B FieldVal"}, "A.Field":[]string{"A FieldVal"}}
	AnonymousSeparate
)

var timeType = reflect.TypeOf(time.Time{})

// Mode specifies which mode the form decoder is to run.
type Mode uint8

// AnonymousMode specifies how data should be rolled up or separated from anonymous structs.
type AnonymousMode uint8
