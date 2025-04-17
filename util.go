package form

import "strconv"

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
