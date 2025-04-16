package form

import (
	"reflect"
	"sort"
	"strings"
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

func (s *structCacheMap) Get(key reflect.Type) (value *cachedStruct, ok bool) {
	value, ok = s.m.Load().(map[reflect.Type]*cachedStruct)[key]
	return
}

func (s *structCacheMap) Set(key reflect.Type, value *cachedStruct) {
	m := s.m.Load().(map[reflect.Type]*cachedStruct)
	nm := make(map[reflect.Type]*cachedStruct, len(m)+1)
	for k, v := range m {
		nm[k] = v
	}

	nm[key] = value
	s.m.Store(nm)
}

func (s *structCacheMap) parseStruct(mode Mode, current reflect.Value, key reflect.Type, tagName string) *cachedStruct {
	s.lock.Lock()
	// could have been multiple trying to access,
	// but once first is done this ensures struct isn't parsed again
	cs, ok := s.Get(key)
	if ok {
		s.lock.Unlock()
		return cs
	}

	var idx int
	var name string
	var isOmitEmpty bool
	var fld reflect.StructField
	typ := current.Type()
	numFields := current.NumField()
	cs = &cachedStruct{fields: make([]cachedField, 0, 4)} // init 4, betting most structs decoding into have at aleast 4 fields
	for i := 0; i < numFields; i++ {
		isOmitEmpty = false
		fld = typ.Field(i)
		if fld.PkgPath != blank && !fld.Anonymous {
			continue
		}

		if s.tagFn != nil {
			name = s.tagFn(fld)
		} else {
			name = fld.Tag.Get(tagName)
		}

		if name == ignore || (mode == ModeExplicit && len(name) == 0) {
			continue
		}

		// check for omitempty
		if idx = strings.LastIndexByte(name, ','); idx != -1 {
			isOmitEmpty = name[idx+1:] == "omitempty"
			name = name[:idx]
		}

		if len(name) == 0 {
			name = fld.Name
		}

		cs.fields = append(cs.fields, cachedField{idx: i, name: name, isAnonymous: fld.Anonymous, isOmitEmpty: isOmitEmpty})
	}

	sort.Sort(cs.fields)
	s.Set(typ, cs)
	s.lock.Unlock()

	return cs
}
