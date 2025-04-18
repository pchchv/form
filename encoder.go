package form

import "net/url"

type encoder struct {
	e         *Encoder
	errs      EncodeErrors
	values    url.Values
	namespace []byte
}
