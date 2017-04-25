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

To fetch a child object:
...
```
o, err := g.Get("outer", "val1")
...
```



