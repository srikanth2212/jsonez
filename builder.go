package jsonez

import (
	"errors"
	"fmt"
	"strings"
)

/**
 * Method to get a child item in an object
 */
func (g *GoJSON) GetObjectEntry(key string) *GoJSON {
	var child *GoJSON

	child = g.child

	for {
		if child != nil && strings.Compare(child.key, key) != 0 {
			child = child.next
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
	entry.key = key

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
	if g.child.key == key {
		cur = g.child
		g.child = cur.next
		cur = nil
		return nil
	}

	prev = g.child
	cur = prev.next

	for {
		if cur == nil {
			errorStr := fmt.Sprintf("%s: Child object with key %s not found", funcName(), key)
			return errors.New(errorStr)
		} else if cur.key == key {
			/*
			 * Check for last element in the list
			 * and handle accordingly
			 */
			if cur.next == nil {
				prev.next = nil
			} else {
				resolveLink(prev, cur.next)
			}

			cur = nil
			break
		} else {
			prev = cur
			cur = cur.next
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

	child = g.child

	for {
		if child != nil {
			size++
			child = child.next
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

	if g.jsontype != JSON_ARRAY {
		errorStr := fmt.Sprintf("%s: Parsing Error", funcName())
		return nil, errors.New(errorStr)
	}

	child = g.child

	for {
		if child != nil && loc > 0 {
			loc--
			child = child.next
		} else {
			break
		}
	}

	return child, nil
}

/**
 * Method to get an array entry based on the value of the element
 */
func (g *GoJSON) GetArrayEntry(val interface{}, jsontype int) (*GoJSON, error) {
	var child *GoJSON

	if g.jsontype != JSON_ARRAY {
		errorStr := fmt.Sprintf("%s: Parsing Error", funcName())
		return nil, errors.New(errorStr)
	}

	child = g.child

	for {
		if child != nil {
			switch jsontype {
			case JSON_INT:
				if child.valueint == val.(int) {
					return child, nil
				}
			case JSON_DOUBLE:
				if child.valuedouble == val.(float64) {
					return child, nil
				}
			case JSON_BOOL:
				if child.valuebool == val.(bool) {
					return child, nil
				}
			case JSON_STRING:
				if child.valuestring == val.(string) {
					return child, nil
				}

			}
			child = child.next
		} else {
			errorStr := fmt.Sprintf("%s: Value not found", funcName())
			return nil, errors.New(errorStr)
		}
	}
}

/**
 * Function to add an entry to a GoJSON array object.
 */
func (g *GoJSON) AddEntryToArray(entry *GoJSON) {
	child := g.child

	if child == nil {
		g.child = entry
		entry.next = nil
	} else {
		for {
			if child != nil && child.next != nil {
				child = child.next
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
		cur = g.child
		g.child = cur.next
		cur = nil
		return nil
	} else if index >= size {
		errorStr := fmt.Sprintf("%s: Index exceeds array size of %d", funcName(), size)
		return errors.New(errorStr)
	} else {
		prev = g.child
		cur = prev.next

		for {
			if i == index {
				/*
				 * Check for last element in the list
				 * and handle accordingly
				 */
				if cur.next == nil {
					prev.next = nil
				} else {
					resolveLink(prev, cur.next)
				}
				cur = nil
				break
			} else {
				prev = cur
				cur = cur.next
				i++
			}
		}
		return nil
	}
}

/**
 * Method to get an array entry based on the value of the element
 */
func (g *GoJSON) DelArrayEntry(val interface{}, jsontype int) error {
	var elem *GoJSON
	elem, err := g.GetArrayEntry(val, jsontype)

	if err != nil {
		return err
	}

	elem.prev.next = elem.next
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
func (g *GoJSON) GetIntVal(keys ...string) (int, error) {
	var cur *GoJSON = g
	for _, key := range keys {
		cur = cur.GetObjectEntry(key)

		if cur == nil {
			errorStr := fmt.Sprintf("%s: Path not found", funcName())
			return 0, errors.New(errorStr)
		}
	}

	if cur.jsontype != JSON_INT {
		errorStr := fmt.Sprintf("%s: key %s is not of type int", funcName(), cur.key)
		return 0, errors.New(errorStr)
	}

	return cur.valueint, nil
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

	if cur.jsontype != JSON_DOUBLE {
		errorStr := fmt.Sprintf("%s: key %s is not of type double", funcName(), cur.key)
		return 0, errors.New(errorStr)
	}

	return cur.valuedouble, nil
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

	if cur.jsontype != JSON_BOOL {
		errorStr := fmt.Sprintf("%s: key %s is not of type double", funcName(), cur.key)
		return false, errors.New(errorStr)
	}

	return cur.valuebool, nil
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

	if cur.jsontype != JSON_STRING {
		errorStr := fmt.Sprintf("%s: key %s is not of type double", funcName(), cur.key)
		return "", errors.New(errorStr)
	}

	return cur.valuestring, nil
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
		cur = AllocNumber(float64(val.(int)), JSON_INT)
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

		if arr.jsontype != JSON_ARRAY {
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
		cur = AllocNumber(float64(val.(int)), JSON_INT)
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
		} else if i == size-1 && cur.key == k {
			key = k
			break
		} else {
			prev = cur
		}
	}

	if prev.jsontype != JSON_OBJECT {
		errorStr := fmt.Sprintf("%s: %s is not a json object", funcName(), prev.key)
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
		} else if i == size-1 && cur.key == k {
			break
		} else {
			prev = cur
		}
	}

	if cur.jsontype != JSON_ARRAY {
		errorStr := fmt.Sprintf("%s: %s is not a json array", funcName(), cur.key)
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
		return cur.DelArrayEntry(val, JSON_INT)

	case JSON_DOUBLE:
		return cur.DelArrayEntry(val, JSON_DOUBLE)

	case JSON_BOOL:
		return cur.DelArrayEntry(val, JSON_BOOL)

	case JSON_STRING:
		return cur.DelArrayEntry(val, JSON_STRING)
	}

	return nil
}