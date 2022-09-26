package subscribe

import (
	"github.com/ml444/samsara/entity"
	"github.com/ml444/samsara/internal"
)

type Handler func(entity entity.IEntity)

type SimpleSubscriber struct {
	ISubscriber

	sequence *internal.Sequence
	barrier  internal.ISubscribeBarrier
	strategy internal.ISubscriberStrategy
	handler  Handler
	isDone   bool
}

func NewSimpleSubscriber(scheduler internal.IScheduler, strategy internal.ISubscriberStrategy, handler Handler) *SimpleSubscriber {
	seq := scheduler.InitConsumerSequence(internal.SequenceInitValue)
	return &SimpleSubscriber{
		sequence: seq,
		barrier:  internal.NewSubscribeBarrier(scheduler, strategy),
		handler:  handler,
	}
}

func (s *SimpleSubscriber) GetSequence() *internal.Sequence {
	return s.sequence
}

func (s *SimpleSubscriber) GetBarrier() {

}

func (s *SimpleSubscriber) Start() {
	// TODO log
	println("===> starting")
	for !s.isDone {
		nextSeq := s.sequence.Get() + 1
		s.barrier.WaitFor(nextSeq)
		var e entity.IEntity
		for {
			e = s.barrier.GetEntity(nextSeq)
			if e != nil {
				break
			}
		}
		s.sequence.Add(1)
		s.handler(e)
	}

}
func (s *SimpleSubscriber) Stop() {
	s.isDone = true
}
