package jsonez

import "strconv"

/**
 * Functions to print the contents
 */

/**
 * Function to print a number
 */
func printNumber(cur *GoJSON) []byte {

	switch cur.Jsontype {
	case JSON_DOUBLE:
		return []byte(strconv.FormatFloat(cur.Valdouble, 'E', -1, 64))

	case JSON_INT:
		return []byte(strconv.FormatInt(int64(cur.Valint), 10))

	case JSON_UINT:
		return []byte(strconv.FormatUint(cur.Valuint, 10))
	}

	return []byte{}
}

/**
 * Function to print an array
 */
func printArray(cur *GoJSON, depth, fmt int) []byte {
	var entryCount int = 0
	var output []byte
	var child *GoJSON

	/*
	 * Get the child count
	 */
	for child = cur.Child; child != nil; child = child.Next {
		entryCount++
	}

	if entryCount == 0 {
		return []byte("[]")
	}

	/*
	 * Print the child entries
	 */
	child = cur.Child
	output = append(output, '[', '\n')

	for i := 0; i < entryCount; i++ {
		if fmt != 0 {
			for j := 0; j < depth; j++ {
				output = append(output, '\t')
			}
		}

		output = append(output, printValue(child, depth, fmt)...)

		/*
		 * Add a "," if this not the last entry
		 */
		if i != entryCount-1 {
			output = append(output, ',')
		}

		output = append(output, '\n')

		child = child.Next
	}

	if fmt != 0 {
		for j := 0; j < depth-1; j++ {
			output = append(output, '\t')
		}
	}

	output = append(output, ']')

	return output
}

/**
 * Function to print an object
 */
func printObject(cur *GoJSON, depth, fmt int) []byte {
	var entryCount int = 0
	var output []byte
	var child *GoJSON

	/*
	 * Get the child count
	 */
	for child = cur.Child; child != nil; child = child.Next {
		entryCount++
	}

	output = append(output, '{')
	output = append(output, '\n')

	if entryCount != 0 {
		/*
		 * Walk the child entries
		 */
		child = cur.Child

		for i := 0; i < entryCount && child != nil; i++ {
			if fmt != 0 {
				for j := 0; j < depth; j++ {
					output = append(output, '\t')
				}
			}

			output = append(output, '"')
			output = append(output, []byte(child.Key)...)
			output = append(output, '"')

			output = append(output, ':')

			if fmt != 0 {
				output = append(output, ' ')
			}

			output = append(output, printValue(child, depth, fmt)...)

			if i != entryCount-1 {
				output = append(output, ',')
			}

			if fmt != 0 {
				output = append(output, '\n')
			}

			child = child.Next
		}
	}

	if fmt != 0 {
		for j := 0; j < depth-1; j++ {
			output = append(output, '\t')
		}
	}

	output = append(output, '}')

	return output
}

/**
 * Function to print the current item
 */
func printValue(cur *GoJSON, depth, fmt int) []byte {
	var output []byte

	switch cur.Jsontype {
	case JSON_NULL:
		output = append(output, []byte("null")...)
		return output

	case JSON_BOOL:
		if cur.Valbool == false {
			output = append(output, []byte("false")...)
		} else {
			output = append(output, []byte("true")...)
		}
		return output

	case JSON_INT:
		fallthrough
	case JSON_UINT:
		fallthrough
	case JSON_DOUBLE:
		output = append(output, printNumber(cur)...)
		return output

	case JSON_STRING:
		output = append(output, '"')
		output = append(output, []byte(cur.Valstr)...)
		output = append(output, '"')
		return output

	case JSON_ARRAY:
		output = append(output, printArray(cur, depth+1, fmt)...)
		return output

	case JSON_OBJECT:
		output = append(output, printObject(cur, depth+1, fmt)...)
		return output
	}

	return []byte{}
}

/**
 * Main function to print the GoJSON tree from root
 */
func GoJSONPrint(root *GoJSON) []byte {
	return printValue(root, 0, 1)
}
