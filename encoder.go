package form

import "net/url"

type encoder struct {
	e         *Encoder
	errs      EncodeErrors
	values    url.Values
	namespace []byte
}

func (e *encoder) setError(namespace []byte, err error) {
	if e.errs == nil {
		e.errs = make(EncodeErrors)
	}

	e.errs[string(namespace)] = err
}

func (e *encoder) setVal(namespace []byte, idx int, vals ...string) {
	arr, ok := e.values[string(namespace)]
	if ok {
		arr = append(arr, vals...)
	} else {
		arr = vals
	}

	e.values[string(namespace)] = arr
}
