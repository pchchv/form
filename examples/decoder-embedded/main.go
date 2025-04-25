package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/pchchv/form"
)

// A ...
type A struct {
	Field string
}

// B ...
type B struct {
	A
	Field string
}

// use a single instance of Decoder, it caches struct info
var decoder *form.Decoder

// parseFormB simulates the results of http.Request's ParseForm() function.
func parseFormB() url.Values {
	return url.Values{
		"Field": []string{"B FieldVal"},
	}
}

// parseFormAB simulates the results of http.Request's ParseForm() function.
func parseFormAB() url.Values {
	return url.Values{
		"Field":   []string{"B FieldVal"},
		"A.Field": []string{"A FieldVal"},
	}
}

func main() {
	var b B
	decoder = form.NewDecoder()
	// this simulates the results of http.Request's ParseForm() function
	values := parseFormB()
	// must pass a pointer
	if err := decoder.Decode(&b, values); err != nil {
		log.Panic(err)
	}

	fmt.Printf("%#v\n", b)

	values = parseFormAB()
	// must pass a pointer
	if err := decoder.Decode(&b, values); err != nil {
		log.Panic(err)
	}

	fmt.Printf("%#v\n", b)
}
