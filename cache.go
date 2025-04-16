package form

import (
	"reflect"
	"sync"
	"sync/atomic"
)

type cachedField struct {
	idx         int
	name        string
	isAnonymous bool
	isOmitEmpty bool
}

type cacheFields []cachedField

func (s cacheFields) Len() int {
	return len(s)
}

func (s cacheFields) Less(i, j int) bool {
	return !s[i].isAnonymous
}

func (s cacheFields) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type cachedStruct struct {
	fields cacheFields
}

// TagNameFunc allows for adding of a custom tag name parser
type TagNameFunc func(field reflect.StructField) string

func newStructCacheMap() *structCacheMap {
	sc := new(structCacheMap)
	sc.m.Store(make(map[reflect.Type]*cachedStruct))
	return sc
}

type structCacheMap struct {
	m     atomic.Value // map[reflect.Type]*cachedStruct
	lock  sync.Mutex
	tagFn TagNameFunc
}
