package models

type (
	Status        int
	ExceptionCode int
)

const (
	SUCCESS   Status = 0
	EXCEPTION Status = 1
)

const (
	EXCEPTION_ID_IDCACHE_INIT_FALSE    ExceptionCode = -1
	EXCEPTION_ID_KEY_NOT_EXISTS                      = -2
	EXCEPTION_ID_TWO_SEGMENTS_ARE_NULL               = -3
)
