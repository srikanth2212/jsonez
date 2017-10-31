package jsonez

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"unicode"
)

/*
 * JSON types
 */
const (
	JSON_BOOL = iota
	JSON_NULL
	JSON_INT
	JSON_UINT
	JSON_DOUBLE
	JSON_STRING
	JSON_ARRAY
	JSON_OBJECT
)

var errorOffset int = -1

/*
 * GoJSON structure
 */
type GoJSON struct {
	/**
	 * Pointers to walk array/object chains
	 */
	Next, Prev *GoJSON

	/**
	 * Child item of the current object
	 */
	Child *GoJSON

	/** JSON type */
	Jsontype int

	/**
	 * Valstr will be set when
	 *type is JSON_STRING
	 */
	Valstr string

	/**
	 * Valint will be set when type
	 * is JSON_INT
	 */
	Valint int64

	/**
	 * Valuint will be set when type
	 * is JSON_UINT
	 */
	Valuint uint64

	/**
	 * valuenum will be set when
	 * type is JSON_DOUBLE
	 */
	Valdouble float64

	/**
	 * Valbool will be set when
	 * type is JSON_BOOL
	 */
	Valbool bool

	/**
	 * JSON Key
	 */
	Key string
}

/**
 * Functions related to parsing a JSON string
 */

/**
 * nextToken finds the next token
 */
func nextToken(input []byte) []byte {
	for i, c := range input {
		if unicode.IsSpace(rune(c)) || c == '\'' {
			continue
		} else {
			return input[i:]
		}
	}

	return []byte{}
}

/**
 * Function to parse a string
 */
func parseString(cur *GoJSON, input []byte) ([]byte, error) {
	var offset int

	if len(input) == 0 {
		errorStr := fmt.Sprintf("%s: Byte slice is empty", funcName())
		return []byte{}, errors.New(errorStr)
	}

	if input[0] != '"' {
		errorStr := fmt.Sprintf("%s: Not a valid string in input %s", funcName(), string(input))
		return nil, errors.New(errorStr)
	}

	for i := range input {
		if i > 0 && input[i] == '"' {
			break
		} else {
			offset++
		}
	}

	cur.Jsontype = JSON_STRING
	cur.Valstr = strings.Trim(string(input[1:offset]), " ")

	return input[offset+1:], nil
}

/**
 * Function to parse a number
 */
func parseNumber(cur *GoJSON, input []byte) ([]byte, error) {
	var n, sign, scale float64
	var subscale, signsubscale, offset int
	var isDouble bool = false

	sign = 1
	subscale = 0
	signsubscale = 1

	if len(input) == 0 {
		errorStr := fmt.Sprintf("%s: Byte slice is empty", funcName())
		return []byte{}, errors.New(errorStr)
	}

	if input[offset] == '-' {
		sign = -1
		offset++
	}

	if unicode.IsDigit(rune(input[offset])) == false {
		errorStr := fmt.Sprintf("%s: Not a valid number in input %s", funcName(), string(input))
		return nil, errors.New(errorStr)
	}

	for {
		if input[offset] == '0' {
			offset++
		} else {
			break
		}
	}

	for {
		if unicode.IsDigit(rune(input[offset])) == true {
			n = (n * 10.0) + float64(input[offset]-'0')
			offset++
		} else {
			break
		}
	}

	if input[offset] == '.' && unicode.IsDigit(rune(input[offset+1])) == true {
		offset++
		isDouble = true

		for {
			if unicode.IsDigit(rune(input[offset])) == true {
				n = (n * 10.0) + float64(input[offset]-'0')
				offset++
				scale--
			} else {
				break
			}
		}
	}

	if input[offset] == 'e' || input[offset] == 'E' {
		offset++
		isDouble = true

		if input[offset] == '-' {
			signsubscale = -1
		}
		offset++

		for {
			if unicode.IsDigit(rune(input[offset])) == true {
				subscale = (subscale * 10) + int(input[offset]-'0')
				offset++
			} else {
				break
			}
		}
	}

	n = sign * n * math.Pow(10.0, scale+float64(subscale)*float64(signsubscale))

	if isDouble == true {
		cur.Valdouble = n
		cur.Jsontype = JSON_DOUBLE
	} else if sign == -1 {
		cur.Valint = int64(n)
		cur.Jsontype = JSON_INT
	} else {
		cur.Valuint = uint64(n)
		cur.Jsontype = JSON_UINT
	}

	return input[offset:], nil
}

/**
 * Function to parse an array
 */
func parseArray(cur *GoJSON, input []byte) ([]byte, error) {
	var child, sibling *GoJSON

	if len(input) == 0 {
		errorStr := fmt.Sprintf("%s: Byte slice is empty", funcName())
		return []byte{}, errors.New(errorStr)
	}

	if input[0] != '[' {
		errorStr := fmt.Sprintf("%s: Not a valid array in input %s", funcName(), string(input))
		return []byte{}, errors.New(errorStr)
	}

	cur.Jsontype = JSON_ARRAY
	input = nextToken(input[1:])

	if len(input) == 0 {
		errorStr := fmt.Sprintf("%s: Could not find any valid JSON", funcName())
		return []byte{}, errors.New(errorStr)
	}

	/*
	 * Check if the array is empty
	 */
	if input[0] == ']' {
		return input[1:], nil
	}

	/*
	 * Allocate memory for the child to
	 * continue processing
	 */
	cur.Child = new(GoJSON)
	child = cur.Child
	input, err := parseValue(child, nextToken(input))

	if err != nil {
		return []byte{}, err
	}

	input = nextToken(input)

	if len(input) == 0 {
		return []byte{}, nil
	}

	/*
	 * Continue processing the array and add the
	 * child entries to this parent GoJSON object
	 */
	for {
		if input[0] == ',' {
			sibling = new(GoJSON)
			child.Next = sibling
			sibling.Prev = child
			child = sibling

			input, err = parseValue(child, nextToken(input[1:]))

			if err != nil {
				return []byte{}, err
			}

			input = nextToken(input)

			if len(input) == 0 {
				return []byte{}, nil
			}
		} else {
			break
		}
	}

	if input[0] == ']' {
		return input[1:], nil
	} else {
		errorStr := fmt.Sprintf("%s: Incomplete/Malformed array in input %s", funcName(), string(input))
		return []byte{}, errors.New(errorStr)
	}
}

/**
 * Function to parse an object
 */
func parseObject(cur *GoJSON, input []byte) ([]byte, error) {
	var child, sibling *GoJSON

	if len(input) == 0 {
		errorStr := fmt.Sprintf("%s: Byte slice is empty", funcName())
		return []byte{}, errors.New(errorStr)
	}

	if input[0] != '{' {
		errorStr := fmt.Sprintf("%s: Not a valid object", funcName())
		return []byte{}, errors.New(errorStr)
	}

	cur.Jsontype = JSON_OBJECT
	input = nextToken(input[1:])

	if len(input) == 0 {
		errorStr := fmt.Sprintf("%s: Could not find any valid JSON in input %s", funcName(), string(input))
		return []byte{}, errors.New(errorStr)
	}

	/*
	 * Check if the object is empty
	 */
	if input[0] == '}' {
		return input[1:], nil
	}

	/*
	 * Allocate memory for the child to
	 * continue processing
	 */
	cur.Child = new(GoJSON)
	child = cur.Child
	input, err := parseString(child, nextToken(input))

	if err != nil {
		return []byte{}, err
	}

	child.Key = child.Valstr
	child.Valstr = ""

	/*
	 * Fetch the location of ':' after the object key
	 */
	input = nextToken(input)
	if len(input) == 0 {
		return []byte{}, nil
	}

	if input[0] != ':' {
		errorStr := fmt.Sprintf("%s: Malformed object in input %s", funcName(), string(input))
		return []byte{}, errors.New(errorStr)
	}

	input, err = parseValue(child, nextToken(input[1:]))

	if err != nil {
		return []byte{}, err
	}

	if input[0] != ',' {
		input = nextToken(input)

		if len(input) == 0 {
			return []byte{}, nil
		}
	}

	/*
	 * Continue processing the object and add the
	 * child entries to this parent GoJSON object
	 */
	for {
		if input[0] == ',' {
			sibling = new(GoJSON)
			child.Next = sibling
			sibling.Prev = child
			child = sibling

			input, err = parseString(child, nextToken(input[1:]))

			if err != nil {
				return []byte{}, err
			}

			child.Key = child.Valstr
			child.Valstr = ""

			/*
			 * Fetch the location of ':' after the object key
			 */
			input = nextToken(input)
			if len(input) == 0 {
				return []byte{}, nil
			}

			if input[0] != ':' {
				errorStr := fmt.Sprintf("%s: Malformed object in input %s", funcName(), string(input))
				return []byte{}, errors.New(errorStr)
			}

			input, err = parseValue(child, nextToken(input[1:]))

			if err != nil {
				return []byte{}, err
			}

			if input[0] != ',' {
				input = nextToken(input)

				if len(input) == 0 {
					return []byte{}, nil
				}
			}

		} else {
			break
		}
	}

	if input[0] == '}' {
		return input[1:], nil
	} else {
		errorStr := fmt.Sprintf("%s: Incomplete/Malformed object in input %s", funcName(), string(input))
		return []byte{}, errors.New(errorStr)
	}
}

/**
 * Function to parse the current token
 */
func parseValue(cur *GoJSON, input []byte) ([]byte, error) {
	input = nextToken(input)
        length := len(input)

	if length == 0 {
            cur.Jsontype = JSON_NULL
            return input[:], nil
	}

	if length >= 4 &&
        strings.Compare(string(input[:4]), "null") == 0 {
		cur.Jsontype = JSON_NULL
		return input[4:], nil
	}

	if length >= 5 && strings.Compare(string(input[:5]), "false") == 0 {
		cur.Jsontype = JSON_BOOL
		cur.Valbool = false
		return input[5:], nil
	}

	if length >= 4 && strings.Compare(string(input[:4]), "true") == 0 {
		cur.Jsontype = JSON_BOOL
		cur.Valbool = true
		return input[4:], nil
	}

	if input[0] >= '0' && input[0] <= '9' {
		return parseNumber(cur, input)
	}

	switch input[0] {
	case '"':
		return parseString(cur, input)

	case '-':
		return parseNumber(cur, input)

	case '[':
		return parseArray(cur, input)

	case '{':
		return parseObject(cur, input)

	}

	errorStr := fmt.Sprintf("%s: Parsing Error", funcName())
	return []byte{}, errors.New(errorStr)
}

/*
 * Function to begin processing the string input
 * as a byte sequence
 */
func GoJSONParse(input []byte) (*GoJSON, error) {
	var g *GoJSON

	g = new(GoJSON)

	input = nextToken(input)

	if len(input) == 0 {
		errorStr := fmt.Sprintf("%s: Could not find any valid JSON in input %s", funcName(), string(input))
		return nil, errors.New(errorStr)
	}

	_, err := parseValue(g, input)

	if err != nil {
		errorStr := fmt.Sprintf("%s: JSON Parse failed with error %s ", funcName(), err)
		return nil, errors.New(errorStr)
	}

	return g, nil
}
