# jsonez
jsonez provides a simple interface to parse, manipulate and emit arbitrary JSON data in go. This implementation doesn't use encoding/json, reflection and is pretty light weight. The project is inspired from cJSON by https://github.com/DaveGamble/cJSON 

## Usage:
```
The following JSON types are defined:
/*
 * JSON types
 */
const (
	JSON_BOOL = iota
	JSON_NULL
	JSON_INT
	JSON_DOUBLE
	JSON_STRING
	JSON_ARRAY
	JSON_OBJECT
)

input := []byte(`{
   "outer": {
	"val1":	"foo",
	"val2":	"bar",
	"val3":	1234,
	"val4":	225.1245,
	"val5":	[
	 	1,
		2,
		3,
		4,
		5
	]
  }
}`)

g, err := GoJSONParse(input)

```
  
To fetch the json output as []byte from the root GoJSON object:

```
output := GoJSONPrint(g)

```

To get a child object:
```
o, err := g.Get("outer", "val1")

```

Adding a new child by value:
```
err = g.AddVal(100, "outer", "val6")

err = g.AddVal(245.67, "outer", "val7")

err = g.AddVal("hello world", "outer", "val8")

err = g.AddVal(true, "outer", "val9")

```

Adding an array:
```
/*
 * The API will create the array object outer and
 * populates 100 as its first element
 */
err = g.AddToArray(100, "outer", "val10")

/*
 * Subsequent adds will append to the same array object
 */
err = g.AddToArray(200.25, "outer", "val10")

err = g.AddToArray("hello world", "outer", "val10")

err = g.AddToArray(true, "outer", "val10")
```

Printing the JSON output after the above operations:
```
fmt.Println(string(GoJSONPrint(g)))

will produce the follwing output:

{
	"outer": {
		"val1": "foo",
		"val2": "bar",
		"val3": 1234,
		"val4": 2.251245E+02,
		"val5": [
			1,
			2,
			3,
			4,
			5
		],
		"val6": 100,
		"val7": 2.4567E+02,
		"val8": "hello world",
		"val9": true,
		"val10": [
			100,
			2.0025E+02,
			"hello world",
			true
		]
	}
}
```
Query for child object:
```
child, err = g.Get("outer", "val9")

```
Getting the value of a child based on type:
```
i, err = g.GetIntVal("outer", "val6")

d, err := g.GetDoubleVal("outer", "val7")

d, err := g.GetStringVal("outer", "val8")

d, err := g.GetbooleVal("outer", "val9")

```

Getting the child object at a specific array index:
```
entry, err := arr.GetArrayElemByIndex(1)

```
Getting the child object based on value (works only for types int, double, string and bool):
```
entry, err := arr.GetArrayEntry(100, JSON_INT)

```

Deleting an Array element based on Index:
```
arr, err := g.Get("outer", "val10")

err = arr.DelIndexFromArray(3)


```
The resultant JSON is:
```
{
	"outer": {
		"val1": "foo",
		"val2": "bar",
		"val3": 1234,
		"val4": 2.251245E+02,
		"val5": [
			1,
			2,
			3,
			4,
			5
		],
		"val6": 100,
		"val7": 2.4567E+02,
		"val8": "hello world",
		"val9": true,
		"val10": [
			100,
			2.0025E+02,
			"hello world"
		]
	}
}

```

Deleting an array element based on value (works only for types int, double, string and bool):
```
err = g.DelFromArray(3, "outer", "val5")

```
The resultant JSON is:
```
{
	"outer": {
		"val1": "foo",
		"val2": "bar",
		"val3": 1234,
		"val4": 2.251245E+02,
		"val5": [
			1,
			2,
			4,
			5
		],
		"val6": 100,
		"val7": 2.4567E+02,
		"val8": "hello world",
		"val9": true,
		"val10": [
			100,
			2.0025E+02,
			"hello world"
		]
	}
}

...
```


Deleting an object based on path:
```go
err = g.DelVal("outer", "val10")

...
```
The resultant JSON is:
```go
fmt.Println(string(GoJSONPrint(g)))
{
	"outer": {
		"val1": "foo",
		"val2": "bar",
		"val3": 1234,
		"val4": 2.251245E+02,
		"val5": [
			1,
			2,
			3,
			4,
			5
		],
		"val6": 100,
		"val7": 2.4567E+02,
		"val8": "hello world",
		"val9": true
	}
}
...
```

