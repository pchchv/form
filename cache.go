package form

type cachedField struct {
	idx         int
	name        string
	isAnonymous bool
	isOmitEmpty bool
}

type cacheFields []cachedField
