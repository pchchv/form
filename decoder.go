package form

import (
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"
)

const (
	errMissingEndBracket   = "invalid formatting for key '%s' missing ']' bracket"
	errMissingStartBracket = "invalid formatting for key '%s' missing '[' bracket"
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

func (d *decoder) parseMapData() {
	if len(d.dm) > 0 {
		// already parsed
		return
	}

	var i, idx, l int
	var rd *recursiveData
	var isNum, insideBracket bool
	d.maxKeyLen = 0
	d.dm = d.dm[0:0]
	for k := range d.values {
		if len(k) > d.maxKeyLen {
			d.maxKeyLen = len(k)
		}

		for i = 0; i < len(k); i++ {
			switch k[i] {
			case '[':
				idx = i
				insideBracket = true
				isNum = true
			case ']':
				if !insideBracket {
					log.Panicf(errMissingStartBracket, k)
				}

				if rd = d.findAlias(k[:idx]); rd == nil {
					l = len(d.dm) + 1
					if l > cap(d.dm) {
						dm := make(dataMap, l)
						copy(dm, d.dm)
						rd = new(recursiveData)
						dm[len(d.dm)] = rd
						d.dm = dm
					} else {
						l = len(d.dm)
						d.dm = d.dm[:l+1]
						rd = d.dm[l]
						rd.sliceLen = 0
						rd.keys = rd.keys[0:0]
					}

					rd.alias = k[:idx]
				}

				// is map + key
				ke := key{
					ivalue:      -1,
					value:       k[idx+1 : i],
					searchValue: k[idx : i+1],
				}

				// is key is number, most likely array key, keep track of just in case an array/slice
				if isNum {
					// no need to check for error, it will always pass
					// as we have done the checking to ensure
					// the value is a number ahead of time
					var err error
					ke.ivalue, err = strconv.Atoi(ke.value)
					if err != nil {
						ke.ivalue = -1
					}

					if ke.ivalue > rd.sliceLen {
						rd.sliceLen = ke.ivalue

					}
				}

				rd.keys = append(rd.keys, ke)
				insideBracket = false
			default:
				// checking if not a number, 0-9 is 48-57 in byte, see for yourself fmt.Println('0', '1', '2', '3', '4', '5', '6', '7', '8', '9')
				if insideBracket && (k[i] > 57 || k[i] < 48) {
					isNum = false
				}
			}
		}

		// if still inside bracket, that means no ending bracket was ever specified
		if insideBracket {
			log.Panicf(errMissingEndBracket, k)
		}
	}
}

func (d *decoder) findAlias(ns string) *recursiveData {
	for i := 0; i < len(d.dm); i++ {
		if d.dm[i].alias == ns {
			return d.dm[i]
		}
	}
	return nil
}
