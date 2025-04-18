package form

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

type encoder struct {
	e         *Encoder
	errs      EncodeErrors
	values    url.Values
	namespace []byte
}

func (e *encoder) getMapKey(key reflect.Value, namespace []byte) (string, bool) {
	v, kind := ExtractType(key)
	if e.e.customTypeFuncs != nil {
		if cf, ok := e.e.customTypeFuncs[v.Type()]; ok {
			if arr, err := cf(v.Interface()); err != nil {
				e.setError(namespace, err)
				return "", false
			} else {
				return arr[0], true
			}
		}
	}

	switch kind {
	case reflect.Interface, reflect.Ptr:
		return "", false
	case reflect.String:
		return v.String(), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), true
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32), true
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), true
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), true
	default:
		e.setError(namespace, fmt.Errorf("Unsupported Map Key '%v' Namespace '%s'", v.String(), namespace))
		return "", false
	}
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
