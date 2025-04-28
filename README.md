# form [![CI](https://github.com/pchchv/form/workflows/CI/badge.svg)](https://github.com/pchchv/form/actions?query=workflow%3ACI+event%3Apush) [![Go Report Card](https://goreportcard.com/badge/github.com/pchchv/form)](https://goreportcard.com/report/github.com/pchchv/form) [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/form)](https://pkg.go.dev/github.com/pchchv/form)

Package `form` decodes url.Values into Go-values and encodes Go-values into url.Values.

## Features
- Allows for Custom Type registration.
- Supports map of almost all types.
- Supports both Numbered and Normal arrays, such as `“Array[0]”` and just `“Array”` with multiple values passed in.
- Supports Encoding & Decoding of almost all Go types, For example, it can Decode to struct, array, map, int... and Encode to struct, array, map, int....
- Slice honours the specified index. For example, if `"Slice[2]” - is the only Slice value passed, it will be placed in index 2, if slice is not large enough, it will be expanded.
- Array honors the specified index. For example, if “Array[2]” - is the only Array value passed, it will be put in index 2, if the array is not large enough, a warning will be printed and the value will be ignored.
- Creates objects only as needed; for example, if no `array` or `map` values are passed, then `array` and `map` are left as their default values in the struct.
- Handles time.Time using RFC3339 time format by default, but can be easily changed by registering a custom type.

## Supported Types (out of the box)
* `string`
* `bool`
* `int`, `int8`, `int16`, `int32`, `int64`
* `uint`, `uint8`, `uint16`, `uint32`, `uint64`
* `float32`, `float64`
* `struct` and `anonymous struct`
* `interface{}`
* `time.Time` - by default using RFC3339
* a `pointer` to one of the above types
* `slice`, `array`
* `map`
* `custom types` can override any of the above types
* many other types may be supported inherently

**NOTE**: `map`, `struct` and `slice` nesting are ad infinitum.

## Installation

```sh
	go get github.com/pchchv/form
```

## Import

```go
	import "github.com/pchchv/form"
```

## Usage

- Use symbol `.` for separating fields/structs. (eg. `structfield.field`)
- Use `[index or key]` for access to index of a slice/array or key for map. (eg. `arrayfield[0]`, `mapfield[keyvalue]`)

```html
<form method="POST">
  <input type="text" name="Name" value="pchchv"/>
  <input type="text" name="Age" value="3"/>
  <input type="text" name="Gender" value="Male"/>
  <input type="text" name="Address[0].Name" value="29 Any street"/>
  <input type="text" name="Address[0].Phone" value="9(999)999-9999"/>
  <input type="text" name="Address[1].Name" value="2 Some Blvd."/>
  <input type="text" name="Address[1].Phone" value="1(111)111-1111"/>
  <input type="text" name="active" value="true"/>
  <input type="text" name="MapExample[key]" value="value"/>
  <input type="text" name="NestedMap[key][key]" value="value"/>
  <input type="text" name="NestedArray[0][0]" value="value"/>
  <input type="submit"/>
</form>
```

## Examples

### Decoding
```go
package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/pchchv/form"
)

// Address contains address information.
type Address struct {
	Name  string
	Phone string
}

// User contains user information.
type User struct {
	Name        string
	Age         uint8
	Gender      string
	Address     []Address
	Active      bool `form:"active"`
	MapExample  map[string]string
	NestedMap   map[string]map[string]string
	NestedArray [][]string
}

// use a single instance of Decoder, it caches struct info
var decoder *form.Decoder

func main() {
	var user User
	decoder = form.NewDecoder()
	values := parseForm() // this simulates the results of http.Request's ParseForm() function
	// must pass a pointer
	if err := decoder.Decode(&user, values); err != nil {
		log.Panic(err)
	}

	fmt.Printf("%#v\n", user)
}

// parseForm simulates the results of http.Request's ParseForm() function.
func parseForm() url.Values {
	return url.Values{
		"Name":                []string{"pchchv"},
		"Age":                 []string{"3"},
		"Gender":              []string{"Male"},
		"Address[0].Name":     []string{"26 Here Blvd."},
		"Address[0].Phone":    []string{"9(999)999-9999"},
		"Address[1].Name":     []string{"26 There Blvd."},
		"Address[1].Phone":    []string{"1(111)111-1111"},
		"active":              []string{"true"},
		"MapExample[key]":     []string{"value"},
		"NestedMap[key][key]": []string{"value"},
		"NestedArray[0][0]":   []string{"value"},
	}
}
```

### Encoding

```go
package main

import (
	"fmt"
	"log"

	"github.com/pchchv/form"
)

// Address contains address information.
type Address struct {
	Name  string
	Phone string
}

// User contains user information.
type User struct {
	Name        string
	Age         uint8
	Gender      string
	Address     []Address
	Active      bool `form:"active"`
	MapExample  map[string]string
	NestedMap   map[string]map[string]string
	NestedArray [][]string
}

// use a single instance of Encoder, it caches struct info
var encoder *form.Encoder

func main() {
	encoder = form.NewEncoder()
	user := User{
		Name:   "pchchv",
		Age:    3,
		Gender: "Male",
		Address: []Address{
			{Name: "26 Here Blvd.", Phone: "9(999)999-9999"},
			{Name: "26 There Blvd.", Phone: "1(111)111-1111"},
		},
		Active:      true,
		MapExample:  map[string]string{"key": "value"},
		NestedMap:   map[string]map[string]string{"key": {"key": "value"}},
		NestedArray: [][]string{{"value"}},
	}

	// must pass a pointer
	values, err := encoder.Encode(&user)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%#v\n", values)
}
```

### Registering Custom Types

#### Decoder

```go
decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
	return time.Parse("2006-01-02", vals[0])
}, time.Time{})
```
If a struct type is registered, the function will only be called if a url.Value exists for the struct and not just the struct fields eg. url.Values{"User":"Name%pchchv"} will call the custom type function with 'User' as the type, however url.Values{"User.Name":"pchchv"} will not.


#### Encoder
```go
encoder.RegisterCustomTypeFunc(func(x interface{}) ([]string, error) {
	return []string{x.(time.Time).Format("2006-01-02")}, nil
}, time.Time{})
```

## Ignoring Fields

It is possible to tell the form to ignore fields by using `-` in the tag.

```go
type MyStruct struct {
	Field string `form:"-"`
}
```

## Omitempty

It is possible to form to omit empty fields using `,omitempty` or `FieldName,omitempty` in the tag.

```go
type MyStruct struct {
	Field  string `form:",omitempty"`
	Field2 string `form:"CustomFieldName,omitempty"`
}
```

## Compatibility

To maximize compatibility with other systems the Encoder attempts to avoid using array indexes in url.Values if at all possible.

```go
// Struct field of
Field []string{"1", "2", "3"}

// Will be output a url.Value as
"Field": []string{"1", "2", "3"}

and not
"Field[0]": []string{"1"}
"Field[1]": []string{"2"}
"Field[2]": []string{"3"}

// however there are times where it is unavoidable, like with pointers
i := int(1)
Field []*string{nil, nil, &i}

// to avoid index 1 and 2 must use index
"Field[2]": []string{"1"}
```