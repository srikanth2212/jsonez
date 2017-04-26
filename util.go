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
	child.jsontype = JSON_NULL

	return child
}

/**
 * Function to create a GoJSON string objext
 */ /**
 * Function to create a GoJSON bool object
 */
func AllocString(val string) *GoJSON {
	child := new(GoJSON)

	child.valuestring = val
	child.jsontype = JSON_STRING

	return child
}

/**
 * Function to create a GoJSON bool object
 */
func AllocBool(val bool) *GoJSON {
	child := new(GoJSON)

	child.jsontype = JSON_BOOL
	if val == false {
		child.valuebool = false
	} else {
		child.valuebool = true
	}

	return child
}

/**
 * Function to create a GoJSON number object
 */
func AllocNumber(val float64, numtype int) *GoJSON {
	child := new(GoJSON)

	if numtype == JSON_INT {
		child.valueint = int(val)
		child.jsontype = JSON_INT
	} else {
		child.valuedouble = val
		child.jsontype = JSON_DOUBLE
	}

	return child
}

/**
 * Function to create a GoJSON array object
 */
func AllocArray() *GoJSON {
	var child *GoJSON
	child = new(GoJSON)
	child.jsontype = JSON_ARRAY

	return child
}

/**
 * Function to create a GoJSON object for
 * an JSON object
 */
func AllocObject() *GoJSON {
	child := new(GoJSON)
	child.jsontype = JSON_OBJECT

	return child
}

/**
 * Function to resolve a link when
 * new entries are added to an array or an
 * object.
 */
func resolveLink(prev *GoJSON, cur *GoJSON) {
	prev.next = cur
	cur.prev = prev
}

/**
 * Function to resolve an interface to a json type
 */
func resolveInterface(v interface{}) (int, error) {
	switch v.(type) {
	case int:
		return JSON_INT, nil
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
