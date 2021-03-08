package models

import (
	"sync"

	"go.uber.org/atomic"
)

type SegmentBuffer struct {
	Key           string
	RWMutex       *sync.RWMutex
	CurrentPos    int
	NextReady     bool
	InitOk        bool
	ThreadRunning *atomic.Bool
	Step          int
	MinStep       int
	UpdatedAt     int64
}
