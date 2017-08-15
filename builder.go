package jsonez

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

/**
 * Method to get a child item in an object
 */
func (g *GoJSON) GetObjectEntry(key string) *GoJSON {
	var child *GoJSON

	child = g.Child

	for {
		if child != nil && strings.Compare(child.Key, key) != 0 {
			child = child.Next
		} else {
			break
		}
	}

	return child
}

/**
 * Function to add a child GoJSON object to a GoJSON object
 */
func (g *GoJSON) AddEntryToObject(key string, entry *GoJSON) {
	entry.Key = key

	g.AddEntryToArray(entry)
}

/**
 * Function to child GoJSON object from a GoJSON object
 */
func (g *GoJSON) DelEntryFromObject(key string) error {
	var cur, prev *GoJSON

	/*
	 * For first element, special processing needs
	 * to be done
	 */
	if g.Child.Key == key {
		cur = g.Child
		g.Child = cur.Next
		cur = nil
		return nil
	}

	prev = g.Child
	cur = prev.Next

	for {
		if cur == nil {
			errorStr := fmt.Sprintf("%s: Child object with key %s not found", funcName(), key)
			return errors.New(errorStr)
		} else if cur.Key == key {
			/*
			 * Check for last element in the list
			 * and handle accordingly
			 */
			if cur.Next == nil {
				prev.Next = nil
			} else {
				resolveLink(prev, cur.Next)
			}

			cur = nil
			break
		} else {
			prev = cur
			cur = cur.Next
		}
	}
	return nil
}

/**
 * Functions to Build a GoJSON tree
 */

/**
 * Method to get the array size
 */
func (g *GoJSON) GetArraySize() int {
	var child *GoJSON
	var size int

	child = g.Child

	for {
		if child != nil {
			size++
			child = child.Next
		} else {
			break
		}
	}

	return size
}

/**
 * Method to get an array entry based on index
 */
func (g *GoJSON) GetArrayElemByIndex(loc int) (*GoJSON, error) {
	var child *GoJSON

	if g.Jsontype != JSON_ARRAY {
		errorStr := fmt.Sprintf("%s: Parsing Error", funcName())
		return nil, errors.New(errorStr)
	}

	child = g.Child

	for {
		if child != nil && loc > 0 {
			loc--
			child = child.Next
		} else {
			break
		}
	}

	return child, nil
}

/**
 * Method to get an array entry based on the value of the element
 */
func (g *GoJSON) GetArrayEntry(val interface{}, Jsontype int) (*GoJSON, error) {
	var child *GoJSON
	var jtype string

	if g.Jsontype != JSON_ARRAY {
		errorStr := fmt.Sprintf("%s: Parsing Error", funcName())
		return nil, errors.New(errorStr)
	}

	child = g.Child

	for {
		if child != nil {
			switch Jsontype {
			case JSON_INT:
				if child.Valint == int64(val.(int)) {
					return child, nil
				}
				jtype = "JSON_INT"
			case JSON_UINT:
				if child.Valuint == uint64(val.(uint)) {
					return child, nil
				}
				jtype = "JSON_UINT"
			case JSON_DOUBLE:
				if child.Valdouble == val.(float64) {
					return child, nil
				}
				jtype = "JSON_DOUBLE"
			case JSON_BOOL:
				if child.Valbool == val.(bool) {
					return child, nil
				}
				jtype = "JSON_BOOL"
			case JSON_STRING:
				if child.Valstr == val.(string) {
					return child, nil
				}
				jtype = "JSON_STRING"

			}
			child = child.Next
		} else {
			errorStr := fmt.Sprintf("%s: Value not found for type %s",
				funcName(), jtype)
			return nil, errors.New(errorStr)
		}
	}
}

/**
 * Function to add an entry to a GoJSON array object.
 */
func (g *GoJSON) AddEntryToArray(entry *GoJSON) {
	child := g.Child

	if child == nil {
		g.Child = entry
		entry.Next = nil
	} else {
		for {
			if child != nil && child.Next != nil {
				child = child.Next
			} else {
				break
			}
		}

		resolveLink(child, entry)
	}
}

/**
 * Function to delete an entry to a GoJSON array object based on index.
 * Index starts at 0.
 */
func (g *GoJSON) DelIndexFromArray(index int) error {

	var cur, prev *GoJSON
	var size int = g.GetArraySize()
	var i int = 1

	if index == 0 {
		cur = g.Child
		g.Child = cur.Next
		cur = nil
		return nil
	} else if index >= size {
		errorStr := fmt.Sprintf("%s: Index exceeds array size of %d", funcName(), size)
		return errors.New(errorStr)
	} else {
		prev = g.Child
		cur = prev.Next

		for {
			if i == index {
				/*
				 * Check for last element in the list
				 * and handle accordingly
				 */
				if cur.Next == nil {
					prev.Next = nil
				} else {
					resolveLink(prev, cur.Next)
				}
				cur = nil
				break
			} else {
				prev = cur
				cur = cur.Next
				i++
			}
		}
		return nil
	}
}

/**
 * Method to get an array entry based on the value of the element
 */
func (g *GoJSON) DelArrayEntry(val interface{}, Jsontype int) error {
	var elem *GoJSON
	elem, err := g.GetArrayEntry(val, Jsontype)

	if err != nil {
		return err
	}

	elem.Prev.Next = elem.Next
	elem = nil

	return nil
}

/**
 * Functions to query the tree based on a path and
 * get the corresponding GoJSON object
 */
func (g *GoJSON) Get(keys ...string) (*GoJSON, error) {
	var cur *GoJSON = g
	for _, key := range keys {
		cur = cur.GetObjectEntry(key)

		if cur == nil {
			errorStr := fmt.Sprintf("%s: Path not found", funcName())
			return nil, errors.New(errorStr)
		}
	}

	return cur, nil
}

/**
 * Functions to query the tree based on a path and
 * get the integer value of the key if exists
 */
func (g *GoJSON) GetIntVal(keys ...string) (int64, error) {
	var cur *GoJSON = g
	for _, key := range keys {
		cur = cur.GetObjectEntry(key)

		if cur == nil {
			errorStr := fmt.Sprintf("%s: Path not found", funcName())
			return 0, errors.New(errorStr)
		}
	}

	if cur.Jsontype != JSON_INT {
		errorStr := fmt.Sprintf("%s: key %s is not of type int", funcName(), cur.Key)
		return 0, errors.New(errorStr)
	}

	return cur.Valint, nil
}

/**
 * Functions to query the tree based on a path and
 * get the unsigned integer value of the key if exists
 */
func (g *GoJSON) GetUIntVal(keys ...string) (uint64, error) {
	var cur *GoJSON = g
	for _, key := range keys {
		cur = cur.GetObjectEntry(key)

		if cur == nil {
			errorStr := fmt.Sprintf("%s: Path not found", funcName())
			return 0, errors.New(errorStr)
		}
	}

	if cur.Jsontype != JSON_UINT {
		errorStr := fmt.Sprintf("%s: key %s is not of type uint", funcName(), cur.Key)
		return 0, errors.New(errorStr)
	}

	return cur.Valuint, nil
}

/**
 * Functions to query the tree based on a path and
 * get the double value of the key if exists
 */
func (g *GoJSON) GetDoubleVal(keys ...string) (float64, error) {
	var cur *GoJSON = g
	for _, key := range keys {
		cur = cur.GetObjectEntry(key)

		if cur == nil {
			errorStr := fmt.Sprintf("%s: Path not found", funcName())
			return 0, errors.New(errorStr)
		}
	}

	if cur.Jsontype != JSON_DOUBLE {
		errorStr := fmt.Sprintf("%s: key %s is not of type double", funcName(), cur.Key)
		return 0, errors.New(errorStr)
	}

	return cur.Valdouble, nil
}

/**
 * Functions to query the tree based on a path and
 * get the bool value of the key if exists
 */
func (g *GoJSON) GetBoolVal(keys ...string) (bool, error) {
	var cur *GoJSON = g
	for _, key := range keys {
		cur = cur.GetObjectEntry(key)

		if cur == nil {
			errorStr := fmt.Sprintf("%s: Path not found", funcName())
			return false, errors.New(errorStr)
		}
	}

	if cur.Jsontype != JSON_BOOL {
		errorStr := fmt.Sprintf("%s: key %s is not of type double", funcName(), cur.Key)
		return false, errors.New(errorStr)
	}

	return cur.Valbool, nil
}

/**
 * Functions to query the tree based on a path and
 * get the string value of the key if exists
 */
func (g *GoJSON) GetStringVal(keys ...string) (string, error) {
	var cur *GoJSON = g
	for _, key := range keys {
		cur = cur.GetObjectEntry(key)

		if cur == nil {
			errorStr := fmt.Sprintf("%s: Path not found", funcName())
			return "", errors.New(errorStr)
		}
	}

	if cur.Jsontype != JSON_STRING {
		errorStr := fmt.Sprintf("%s: key %s is not of type double", funcName(), cur.Key)
		return "", errors.New(errorStr)
	}

	return cur.Valstr, nil
}

/**
 * Functions to query the tree based on a path and
 * add a new int, double, bool or sting value
 */
func (g *GoJSON) AddVal(val interface{}, paths ...string) error {
	var cur, prev *GoJSON
	size := len(paths)
	var key string

	prev = g

	for i, k := range paths {
		cur = prev.GetObjectEntry(k)
		key = k

		if cur == nil {
			if i == size-1 {
				cur = new(GoJSON)
				break
			} else {
				errorStr := fmt.Sprintf("%s: Path %s not found", funcName(), k)
				return errors.New(errorStr)
			}
		}
		prev = cur
	}

	/*
	 * Get the json type of the value to be added
	 */
	t, err := resolveInterface(val)
	if err != nil {
		return err
	}

	switch t {
	case JSON_INT:
		v := reflect.ValueOf(val)
		if val.(int) < 0 {
			cur = AllocNumber(float64(v.Int()), JSON_INT)
		} else {
			cur = AllocNumber(float64(uint64(v.Int())), JSON_UINT)
		}

		prev.AddEntryToObject(key, cur)
	case JSON_UINT:
		v := reflect.ValueOf(val)
		cur = AllocNumber(float64(v.Uint()), JSON_UINT)
		prev.AddEntryToObject(key, cur)
	case JSON_DOUBLE:
		cur = AllocNumber(val.(float64), JSON_DOUBLE)
		prev.AddEntryToObject(key, cur)
	case JSON_BOOL:
		cur = AllocBool(val.(bool))
		prev.AddEntryToObject(key, cur)
	case JSON_STRING:
		cur = AllocString(val.(string))
		prev.AddEntryToObject(key, cur)
	}

	return nil
}

/**
 * Function to append an entry to an array. The array will
 * be created if it doesn't exist
 */
func (g *GoJSON) AddToArray(val interface{}, paths ...string) error {
	var prev, arr, cur *GoJSON
	size := len(paths)
	var key string

	prev = g

	for i, k := range paths {
		cur = prev.GetObjectEntry(k)
		key = k

		if cur == nil {
			if i == size-1 {
				break
			} else {
				errorStr := fmt.Sprintf("%s: Path %s not found", funcName(), k)
				return errors.New(errorStr)
			}
		} else {
			prev = cur
		}
	}

	if cur == nil {
		arr = AllocArray()
		prev.AddEntryToObject(key, arr)
	} else {
		arr = cur

		if arr.Jsontype != JSON_ARRAY {
			errorStr := fmt.Sprintf("%s: GoJSON object with key %s is not of type array", funcName(), key)
			return errors.New(errorStr)
		}
	}

	/*
	 * Get the json type of the value to be added
	 */
	t, err := resolveInterface(val)
	if err != nil {
		return err
	}

	switch t {
	case JSON_INT:
		v := reflect.ValueOf(val)
		if val.(int) < 0 {
			cur = AllocNumber(float64(v.Int()), JSON_INT)
		} else {
			cur = AllocNumber(float64(uint64(v.Int())), JSON_UINT)
		}
		arr.AddEntryToArray(cur)
	case JSON_UINT:
		v := reflect.ValueOf(val)
		cur = AllocNumber(float64(v.Uint()), JSON_UINT)
		arr.AddEntryToArray(cur)
	case JSON_DOUBLE:
		cur = AllocNumber(val.(float64), JSON_DOUBLE)
		arr.AddEntryToArray(cur)
	case JSON_BOOL:
		cur = AllocBool(val.(bool))
		arr.AddEntryToArray(cur)
	case JSON_STRING:
		cur = AllocString(val.(string))
		arr.AddEntryToArray(cur)
	}

	return nil
}

/**
 * Functions to query the tree based on a path and
 * delete a child object
 */
func (g *GoJSON) DelVal(paths ...string) error {
	var cur, prev *GoJSON
	prev = g
	var key string
	size := len(paths)

	for i, k := range paths {
		cur = prev.GetObjectEntry(k)

		if cur == nil {
			errorStr := fmt.Sprintf("%s: Path %s not found", funcName(), k)
			return errors.New(errorStr)
		} else if i == size-1 && cur.Key == k {
			key = k
			break
		} else {
			prev = cur
		}
	}

	if prev.Jsontype != JSON_OBJECT {
		errorStr := fmt.Sprintf("%s: %s is not a json object", funcName(), prev.Key)
		return errors.New(errorStr)
	}

	return prev.DelEntryFromObject(key)
}

/**
 * Function to delete an entry from an array based on the value. The last entry
 * in the path should be an array
 */
func (g *GoJSON) DelFromArray(val interface{}, paths ...string) error {
	var cur, prev *GoJSON
	prev = g
	size := len(paths)

	for i, k := range paths {
		cur = prev.GetObjectEntry(k)

		if cur == nil {
			errorStr := fmt.Sprintf("%s: Path %s not found", funcName(), k)
			return errors.New(errorStr)
		} else if i == size-1 && cur.Key == k {
			break
		} else {
			prev = cur
		}
	}

	if cur.Jsontype != JSON_ARRAY {
		errorStr := fmt.Sprintf("%s: %s is not a json array", funcName(), cur.Key)
		return errors.New(errorStr)
	}

	/*
	 * Get the json type of the value to be added
	 */
	t, err := resolveInterface(val)
	if err != nil {
		return err
	}

	switch t {
	case JSON_INT:
		if val.(int) < 0 {
			return cur.DelArrayEntry(val, JSON_INT)
		} else {
			return cur.DelArrayEntry(uint(val.(int)), JSON_UINT)
		}

	case JSON_UINT:
		return cur.DelArrayEntry(val, JSON_UINT)

	case JSON_DOUBLE:
		return cur.DelArrayEntry(val, JSON_DOUBLE)

	case JSON_BOOL:
		return cur.DelArrayEntry(val, JSON_BOOL)

	case JSON_STRING:
		return cur.DelArrayEntry(val, JSON_STRING)
	}

	return nil
}
