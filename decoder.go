package form

import "net/url"

type decoder struct {
	d         *Decoder
	dm        dataMap
	errs      DecodeErrors
	values    url.Values
	maxKeyLen int
	namespace []byte
}
