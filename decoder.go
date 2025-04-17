package form

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

type decoder struct {
	d         *Decoder
	dm        dataMap
	errs      DecodeErrors
	values    url.Values
	maxKeyLen int
	namespace []byte
}

func (d *decoder) getMapKey(key string, current reflect.Value, namespace []byte) (err error) {
	v, kind := ExtractType(current)
	if d.d.customTypeFuncs != nil {
		if cf, ok := d.d.customTypeFuncs[v.Type()]; ok {
			val, e := cf([]string{key})
			if e != nil {
				return e
			}

			v.Set(reflect.ValueOf(val))
			return
		}
	}

	switch kind {
	case reflect.Interface:
		// If interface would have been set on the struct before decoding,
		// say to a struct value we would not get here but kind would be struct
		v.Set(reflect.ValueOf(key))
		return
	case reflect.Ptr:
		newVal := reflect.New(v.Type().Elem())
		if err = d.getMapKey(key, newVal.Elem(), namespace); err == nil {
			v.Set(newVal)
		}
	case reflect.String:
		v.SetString(key)
	case reflect.Uint, reflect.Uint64:
		u64, e := strconv.ParseUint(key, 10, 64)
		if e != nil {
			return fmt.Errorf("invalid unsigned integer value '%s' type '%v' namespace '%s'", key, v.Type(), string(namespace))
		}
		v.SetUint(u64)
	case reflect.Uint8:
		u64, e := strconv.ParseUint(key, 10, 8)
		if e != nil {
			return fmt.Errorf("invalid unsigned integer value '%s' type '%v' namespace '%s'", key, v.Type(), string(namespace))
		}
		v.SetUint(u64)
	case reflect.Uint16:
		u64, e := strconv.ParseUint(key, 10, 16)
		if e != nil {
			return fmt.Errorf("invalid unsigned integer value '%s' type '%v' namespace '%s'", key, v.Type(), string(namespace))
		}
		v.SetUint(u64)
	case reflect.Uint32:
		u64, e := strconv.ParseUint(key, 10, 32)
		if e != nil {
			return fmt.Errorf("invalid unsigned integer value '%s' type '%v' namespace '%s'", key, v.Type(), string(namespace))
		}
		v.SetUint(u64)
	case reflect.Int, reflect.Int64:
		i64, e := strconv.ParseInt(key, 10, 64)
		if e != nil {
			return fmt.Errorf("invalid integer value '%s' type '%v' namespace '%s'", key, v.Type(), string(namespace))
		}
		v.SetInt(i64)
	case reflect.Int8:
		i64, e := strconv.ParseInt(key, 10, 8)
		if e != nil {
			return fmt.Errorf("invalid integer value '%s' type '%v' namespace '%s'", key, v.Type(), string(namespace))
		}
		v.SetInt(i64)
	case reflect.Int16:
		i64, e := strconv.ParseInt(key, 10, 16)
		if e != nil {
			return fmt.Errorf("invalid integer value '%s' type '%v' namespace '%s'", key, v.Type(), string(namespace))
		}
		v.SetInt(i64)
	case reflect.Int32:
		i64, e := strconv.ParseInt(key, 10, 32)
		if e != nil {
			return fmt.Errorf("invalid integer value '%s' type '%v' namespace '%s'", key, v.Type(), string(namespace))
		}
		v.SetInt(i64)
	case reflect.Float32:
		f, e := strconv.ParseFloat(key, 32)
		if e != nil {
			return fmt.Errorf("invalid float value '%s' type '%v' namespace '%s'", key, v.Type(), string(namespace))
		}
		v.SetFloat(f)
	case reflect.Float64:
		f, e := strconv.ParseFloat(key, 64)
		if e != nil {
			return fmt.Errorf("invalid float value '%s' type '%v' namespace '%s'", key, v.Type(), string(namespace))
		}
		v.SetFloat(f)
	case reflect.Bool:
		b, e := parseBool(key)
		if e != nil {
			return fmt.Errorf("Invalid Boolean Value '%s' Type '%v' Namespace '%s'", key, v.Type(), string(namespace))
		}
		v.SetBool(b)
	default:
		return fmt.Errorf("unsupported map key '%s', type '%v' namespace '%s'", key, v.Type(), string(namespace))
	}

	return
}
