# jsonez
jsonez provides a simple interface to parse, manipulate and emit arbitrary JSON data in golang. This implementation doesn't use encoding/json, reflection and is pretty light weight. The project is inspired from cJSON by https://github.com/DaveGamble/cJSON 

## How to install:
```bash
go get github.com/srikanth2212/jsonez
```
## Usage:

## How to parse and query:
```go
...
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
...
```
  
To fetch the json output in []byte from GoJSON object:
...
```
output := GoJSONPrint(g)
...
```

To get a child object:
```go
o, err := g.Get("outer", "val1")
...
```

Adding a new child by value:
```go
err = g.AddVal(100, "outer", "val6")

err = g.AddVal(245.67, "outer", "val7")

err = g.AddVal("hello world", "outer", "val8")

err = g.AddVal(true, "outer", "val9")

...
```

Adding an array:
```go

/*
 * The API will create the array object val10 and
 * populates 100 as its first element
 */
err = g.AddToArray(100, "outer", "val10")

/*
 * Subsequent adds will append to the same array object
 */
err = g.AddToArray(200.25, "outer", "val10")

err = g.AddToArray("hello world", "outer", "val10")

err = g.AddToArray(true, "outer", "val10")

...
```

Printing the JSON output after the above operations:
```go

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
...
```




