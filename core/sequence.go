package core

import (
	"github.com/ml444/samsara/utils"
	"strconv"
)

type Sequence struct {
	value utils.AtomicInt64
}

func NewSequence() *Sequence {
	s := &Sequence{}
	s.Set(-1)
	return s
}

func (s *Sequence) Init(initValue int64) {
	if initValue == 0 {
		s.Set(-1)
	} else {
		s.Set(initValue)
	}
}
func (s *Sequence) Get() int64 {
	return s.value.Get()
}

func (s *Sequence) Set(value int64) {
	s.value.Set(value)
}
func (s *Sequence) Add(value int64) int64 {
	return s.value.Add(value)
}

func (s *Sequence) CompareAndSwap(oldValue int64, newValue int64) bool {
	return s.value.CompareAndSwap(oldValue, newValue)
}

func (s *Sequence) String() string {
	return strconv.FormatInt(s.Get(), 10)
}

func (s *Sequence) IncrementAndGet() int64 {
	return s.Add(1)
}
