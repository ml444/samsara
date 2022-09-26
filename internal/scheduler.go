package internal

import (
	"math"
)

type IScheduler interface {
	InitConsumerSequence(initValue int64) *Sequence
	//InitConsumerSequences(sequences ...*Sequence)
	AddSequences(sequences ...*Sequence)
	GetConsumerSequences() []*Sequence
	GetMinConsumerSeq(defaultValue int64) int64
	GetMinConsumerSequence() int64
	IsAvailable(int64) bool
	SetAvailableArrayValue(int64)
	GetCursor() *Sequence
	GetBufferSize() int64
	GetIndexMask() int64
	SetEntity(index int64, entity interface{})
	GetEntity(seq int64) interface{}
}

type Scheduler struct {
	IScheduler
	cursor            *Sequence
	consumerSequences []*Sequence
	ringBuffer        *RingBuffer
	indexShift        int
	availableArray    []int
}

func NewScheduler(ringBuffer *RingBuffer) *Scheduler {
	bufferSize := ringBuffer.Size()
	return &Scheduler{
		cursor:         NewSequence(),
		ringBuffer:     ringBuffer,
		indexShift:     int(math.Log2(float64(bufferSize))),
		availableArray: make([]int, int(bufferSize)),
	}
}

func (s *Scheduler) InitConsumerSequence(initValue int64) *Sequence {
	sequence := &Sequence{}
	sequence.Set(initValue)
	//s.consumerSequences = append(s.consumerSequences, sequence)
	return sequence
}

//func (s *Scheduler) InitConsumerSequences(sequences ...*Sequence) {
//	s.consumerSequences = sequences
//}

func (s *Scheduler) AddSequences(sequences ...*Sequence) {
	s.consumerSequences = append(s.consumerSequences, sequences...)
}

func (s *Scheduler) GetBufferSize() int64 {
	return s.ringBuffer.Size()
}

func (s *Scheduler) GetIndexMask() int64 {
	return s.ringBuffer.IndexMask()
}
func (s *Scheduler) GetCursor() *Sequence {
	return s.cursor
}

func (s *Scheduler) RemainingCap() int64 {
	curValue := s.cursor.Get()
	minSeq := s.GetMinConsumerSeq(curValue)
	return s.GetBufferSize() - (curValue - minSeq)
}

//func (s *Scheduler) HasAvailableCap(availableCap int) bool {
//	return s.claimStrategy.HasAvailableCap(availableCap, s.consumerSequences)
//}

func (s *Scheduler) GetConsumerSequences() []*Sequence {
	return s.consumerSequences
}
func (s *Scheduler) GetMinConsumerSeq(defaultValue int64) int64 {
	minSeq := defaultValue
	for _, sequence := range s.consumerSequences {
		if sequence.Get() < minSeq {
			minSeq = sequence.Get()
		}
	}
	return minSeq
}

func (s *Scheduler) GetMinConsumerSequence() int64 {
	var minSeq int64 = math.MaxInt64
	for _, sequence := range s.consumerSequences {
		if sequence.Get() < minSeq {
			minSeq = sequence.Get()
		}
	}
	return minSeq
	//return s.GetMinConsumerSeq(math.MaxInt64)
}
func (s *Scheduler) calculateIndex(seq int64) int {
	return int(seq & s.GetIndexMask())
}

func (s *Scheduler) calculateFlag(seq int64) int {
	return int(uint64(seq) >> s.indexShift)
}

func (s *Scheduler) getAvailableArrayFlag(index int) int {
	return s.availableArray[index]
}

func (s *Scheduler) IsAvailable(seq int64) bool {
	index := s.calculateIndex(seq)
	flag := s.calculateFlag(seq)
	return s.getAvailableArrayFlag(index) == flag
}

func (s *Scheduler) SetAvailableArrayValue(seq int64) {
	index := s.calculateIndex(seq)
	flag := s.calculateFlag(seq)
	s.availableArray[index] = flag
}

func (s *Scheduler) SetEntity(seq int64, entity interface{}) {
	s.ringBuffer.SetEntity(seq, entity)
}
func (s *Scheduler) GetEntity(seq int64) interface{} {
	return s.ringBuffer.GetEntity(seq)
}
