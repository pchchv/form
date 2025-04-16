package form

const (
	blank             = ""
	ignore            = "-"
	fieldNS           = "Field Namespace:"
	errorText         = " ERROR:"
	ModeImplicit Mode = iota // ModeImplicit tries to parse values for all fields that do not have an ignore '-' tag.
	ModeExplicit             // ModeExplicit parses values for field with a field tag and that tag is not the ignore '-' tag.
)

// Mode specifies which mode the form decoder is to run
type Mode uint8
