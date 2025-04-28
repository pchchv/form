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