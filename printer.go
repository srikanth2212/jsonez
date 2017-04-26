package jsonez

import "strconv"

/**
 * Functions to print the contents
 */

/**
 * Function to print a number
 */
func printNumber(cur *GoJSON) []byte {

	switch cur.jsontype {
	case JSON_DOUBLE:
		return []byte(strconv.FormatFloat(cur.valuedouble, 'E', -1, 64))

	case JSON_INT:
		return []byte(strconv.FormatInt(int64(cur.valueint), 10))
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
	for child = cur.child; child != nil; child = child.next {
		entryCount++
	}

	if entryCount == 0 {
		return []byte("[]")
	}

	/*
	 * Print the child entries
	 */
	child = cur.child
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

		child = child.next
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
	for child = cur.child; child != nil; child = child.next {
		entryCount++
	}

	output = append(output, '{')
	output = append(output, '\n')

	if entryCount != 0 {
		/*
		 * Walk the child entries
		 */
		child = cur.child

		for i := 0; i < entryCount && child != nil; i++ {
			if fmt != 0 {
				for j := 0; j < depth; j++ {
					output = append(output, '\t')
				}
			}

			output = append(output, '"')
			output = append(output, []byte(child.key)...)
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

			child = child.next
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

	switch cur.jsontype {
	case JSON_NULL:
		output = append(output, []byte("null")...)
		return output

	case JSON_BOOL:
		if cur.valuebool == false {
			output = append(output, []byte("false")...)
		} else {
			output = append(output, []byte("true")...)
		}
		return output

	case JSON_INT:
		fallthrough
	case JSON_DOUBLE:
		output = append(output, printNumber(cur)...)
		return output

	case JSON_STRING:
		output = append(output, '"')
		output = append(output, []byte(cur.valuestring)...)
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
