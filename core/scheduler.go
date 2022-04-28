package core

import (
	"math"
)

type IScheduler interface {
	InitConsumerSequence(initValue int64) *Sequence
	InitConsumerSequences(sequences ...*Sequence)
	AddSequences(sequences ...*Sequence)
	GetConsumerSequences() []*Sequence
	GetMinConsumerSequence() int64
	GetCursor() *Sequence
	GetBufferSize() int64
	SetEntity(index int64, entity IEntity)
	GetEntity(seq int64) IEntity
}

type Scheduler struct {
	//InitialCursorValue int64
	cursor            *Sequence
	consumerSequences []*Sequence
	ringBuffer        *RingBuffer
}

func NewScheduler(ringBuffer *RingBuffer) *Scheduler {
	return &Scheduler{
		cursor:     NewSequence(),
		ringBuffer: ringBuffer,
	}
}

func (s *Scheduler) InitConsumerSequence(initValue int64) *Sequence {
	if initValue == 0 {
		initValue = -1
	}
	sequence := &Sequence{}
	sequence.Set(initValue)
	//s.consumerSequences = append(s.consumerSequences, sequence)
	return sequence
}
func (s *Scheduler) InitConsumerSequences(sequences ...*Sequence) {
	s.consumerSequences = sequences
}
func (s *Scheduler) AddSequences(sequences ...*Sequence) {
	s.consumerSequences = append(s.consumerSequences, sequences...)
}

func (s *Scheduler) GetBufferSize() int64 {
	return s.ringBuffer.Size()
}

func (s *Scheduler) GetCursor() *Sequence {
	return s.cursor
}

//func (s *Scheduler) HasAvailableCap(availableCap int) bool {
//	return s.claimStrategy.HasAvailableCap(availableCap, s.consumerSequences)
//}
func (s *Scheduler) GetConsumerSequences() []*Sequence {
	return s.consumerSequences
}
func (s *Scheduler) GetMinConsumerSequence() int64 {
	var minSeq int64 = math.MaxInt64
	for _, sequence := range s.consumerSequences {
		if sequence.Get() < minSeq {
			minSeq = sequence.Get()
		}
	}
	return minSeq
}

func (s *Scheduler) RemainingCap() int64 {
	consumed := s.GetMinConsumerSequence()
	produced := s.cursor.Get()
	return s.GetBufferSize() - (produced - consumed)
}

func (s Scheduler) SetEntity(seq int64, entity IEntity) {
	//s.cursor.Add(1)
	s.ringBuffer.SetEntity(seq, entity)
}
func (s Scheduler) GetEntity(seq int64) IEntity {
	return s.ringBuffer.GetEntity(seq)
}
