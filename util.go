package form

import (
	"reflect"
	"strconv"
)

// ExtractType gets the actual underlying type of field value.
func ExtractType(current reflect.Value) (reflect.Value, reflect.Kind) {
	switch current.Kind() {
	case reflect.Ptr:
		if current.IsNil() {
			return current, reflect.Ptr
		} else {
			return ExtractType(current.Elem())
		}
	case reflect.Interface:
		if current.IsNil() {
			return current, reflect.Interface
		} else {
			return ExtractType(current.Elem())
		}
	default:
		return current, current.Kind()
	}
}

func parseBool(str string) (bool, error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "on", "yes", "ok":
		return true, nil
	case "", "0", "f", "F", "false", "FALSE", "False", "off", "no":
		return false, nil
	default:
		// strconv.NumError exactly mimics the strconv.ParseBool(..)
		// error and type to ensure compatibility with std library and others
		return false, &strconv.NumError{Func: "ParseBool", Num: str, Err: strconv.ErrSyntax}
	}
}
