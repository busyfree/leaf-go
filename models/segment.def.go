package models

import "go.uber.org/atomic"

type Segment struct {
	Value *atomic.Int64
	Max   int64
	Step  int
}
