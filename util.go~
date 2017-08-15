package jsonez

import (
	"errors"
	"fmt"
	"runtime"
)

/**
 * Utility functions used by other functions/methods
 */

/**
 * Function to get the current function name for
 * Error reporting
 */
func funcName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

/**
 * Function to create a null object
 */
func AllocNull() *GoJSON {
	child := new(GoJSON)
	child.Jsontype = JSON_NULL

	return child
}

/**
 * Function to create a GoJSON string objext
 */ /**
 * Function to create a GoJSON bool object
 */
func AllocString(val string) *GoJSON {
	child := new(GoJSON)

	child.Valstr = val
	child.Jsontype = JSON_STRING

	return child
}

/**
 * Function to create a GoJSON bool object
 */
func AllocBool(val bool) *GoJSON {
	child := new(GoJSON)

	child.Jsontype = JSON_BOOL
	if val == false {
		child.Valbool = false
	} else {
		child.Valbool = true
	}

	return child
}

/**
 * Function to create a GoJSON number object
 */
func AllocNumber(val float64, numtype int) *GoJSON {
	child := new(GoJSON)

	if numtype == JSON_INT {
		child.Valint = int64(val)
		child.Jsontype = JSON_INT
	} else if numtype == JSON_UINT {
		child.Valuint = uint64(val)
		child.Jsontype = JSON_UINT
	} else {
		child.Valdouble = val
		child.Jsontype = JSON_DOUBLE
	}

	return child
}

/**
 * Function to create a GoJSON array object
 */
func AllocArray() *GoJSON {
	var child *GoJSON
	child = new(GoJSON)
	child.Jsontype = JSON_ARRAY

	return child
}

/**
 * Function to create a GoJSON object for
 * an JSON object
 */
func AllocObject() *GoJSON {
	child := new(GoJSON)
	child.Jsontype = JSON_OBJECT

	return child
}

/**
 * Function to resolve a link when
 * new entries are added to an array or an
 * object.
 */
func resolveLink(prev *GoJSON, cur *GoJSON) {
	prev.Next = cur
	cur.Prev = prev
}

/**
 * Function to resolve an interface to a json type
 */
func resolveInterface(v interface{}) (int, error) {
	switch v.(type) {
	case int:
		return JSON_INT, nil
	case int8:
		return JSON_INT, nil
	case int16:
		return JSON_INT, nil
	case int32:
		return JSON_INT, nil
	case int64:
		return JSON_INT, nil
	case uint:
		return JSON_UINT, nil
	case uint8:
		return JSON_UINT, nil
	case uint16:
		return JSON_UINT, nil
	case uint32:
		return JSON_UINT, nil
	case uint64:
		return JSON_UINT, nil
	case float64:
		return JSON_DOUBLE, nil
	case bool:
		return JSON_BOOL, nil
	case string:
		return JSON_STRING, nil
	default:
		errorStr := fmt.Sprintf("%s: Unknown data type", funcName())
		return -1, errors.New(errorStr)
	}

}
