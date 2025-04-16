package form

import (
	"reflect"
	"testing"

	"github.com/go-playground/assert/v2"
)

type Struct struct {
	Array []int
}

func TestDecoderMultipleSimultaniousParseStructRequests(t *testing.T) {
	var test Struct
	sc := newStructCacheMap()
	proceed := make(chan struct{})
	sv := reflect.ValueOf(test)
	typ := sv.Type()
	for i := 0; i < 200; i++ {
		go func() {
			<-proceed
			s := sc.parseStruct(ModeImplicit, sv, typ, "form")
			assert.NotEqual(t, s, nil)
		}()
	}

	close(proceed)
}
