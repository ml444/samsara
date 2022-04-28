package subscribe

import (
	"fmt"
	"github.com/ml444/samsara/core"
	"os"
)

type SimpleSubscriber struct {
	sequence *core.Sequence
	barrier  core.ISubscribeBarrier
	strategy core.ISubscriberStrategy
	isDone   bool
}

func NewSimpleSubscriber(scheduler core.IScheduler, strategy core.ISubscriberStrategy) *SimpleSubscriber {
	seq := scheduler.InitConsumerSequence(0)
	return &SimpleSubscriber{
		sequence: seq,
		barrier:  core.NewSubscribeBarrier(scheduler, strategy),
	}
}

func (s *SimpleSubscriber) GetSequence() *core.Sequence {
	return s.sequence
}

func (s *SimpleSubscriber) GetBarrier() {

}

func (s *SimpleSubscriber) Start() {
	fmt.Println("===> starting")
	for !s.isDone {
		nextSeq := s.sequence.Get() + 1
		s.barrier.WaitFor(nextSeq)
		seq := s.sequence.IncrementAndGet()
		var e core.IEntity
		for {
			e = s.barrier.GetEntity(seq)
			if e != nil {
				break
			}
		}
		s.Handler(e)
	}

}
func (s *SimpleSubscriber) Stop() {
	s.isDone = true
}

func (s *SimpleSubscriber) Handler(entity core.IEntity) {
	buf = append(buf, entity.DataByte()...)
	if len(buf) > 2*1024 {
		_, _ = DestFile.Write(buf)
		buf = []byte{}
	}
}

func FileFlush() {
	DestFile.Write(buf)
}

var buf []byte
var DestFile *os.File

func init() {
	var err error
	path := "./log_test.log"
	DestFile, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("open file fail, path %s, err %s", path, err)
		return
	}
}
