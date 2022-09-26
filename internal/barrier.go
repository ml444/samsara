package internal

import (
	"errors"
	"fmt"
)

type IPublishBarrier interface {
	Next() int64
	NextN(n int64) int64
	TryNext() (int64, error)
	TryNextN(n int64) (int64, error)
	Commit(seq int64, entity interface{})
}
type ISubscribeBarrier interface {
	WaitFor(sequence int64) int64
	GetEntity(sequence int64) interface{}
}

type BasePublishBarrier struct {
	bufferSize  int64
	minSeqCache int64
	scheduler   IScheduler
	pubStrategy IPublisherStrategy
}

func (b *BasePublishBarrier) Next() int64 {
	return b.NextN(1)
}

func (b *BasePublishBarrier) NextN(n int64) int64 {
	if n < 1 || n > b.bufferSize {
		panic(fmt.Sprintf("must be 0 < n < %d", b.bufferSize))
	}

	cursor := b.scheduler.GetCursor()
DO:
	cursorValue := cursor.Get()
	nextSeq := cursorValue + n
	wrapPoint := nextSeq - b.bufferSize
	minSeq := b.minSeqCache
	if wrapPoint > minSeq || minSeq > cursorValue {
		var newMinSeq int64
		for wrapPoint > newMinSeq {
			b.pubStrategy.Wait()
			newMinSeq = b.scheduler.GetMinConsumerSequence()
			//newMinSeq = b.scheduler.GetMinConsumerSeq(9223372036854775797)
		}
		b.minSeqCache = newMinSeq
	}
	if !cursor.CompareAndSwap(cursorValue, nextSeq) {
		goto DO
	}

	return nextSeq
}

func (b *BasePublishBarrier) TryNext() (int64, error) {
	return b.TryNextN(1)
}

func (b *BasePublishBarrier) TryNextN(n int64) (int64, error) {
	if n < 1 || n > b.bufferSize {
		panic(fmt.Sprintf("must be 0 < n < %d", b.bufferSize))
	}
	cursor := b.scheduler.GetCursor()
DO:
	cursorValue := cursor.Get()
	nextSeq := cursorValue + n
	if !b.hasAvailableCap(n, cursorValue) {
		return 0, errors.New("insufficient capacity")
	}
	if !cursor.CompareAndSwap(cursorValue, nextSeq) {
		goto DO
	}
	return nextSeq, nil

}

func (b *BasePublishBarrier) hasAvailableCap(n int64, curValue int64) bool {
	nextSeq := curValue + n
	wrapPoint := nextSeq - b.bufferSize
	minSeq := b.minSeqCache
	if wrapPoint > minSeq || minSeq > curValue {
		newMinSeq := b.scheduler.GetMinConsumerSeq(curValue)
		b.minSeqCache = newMinSeq
		if wrapPoint > newMinSeq {
			return false
		}
	}
	return true
}

func (b *BasePublishBarrier) Commit(seq int64, entity interface{}) {
	b.scheduler.SetEntity(seq, entity)
}

//func NewBasePublishBarrier(scheduler IScheduler, strategy IPublisherStrategy) *BasePublishBarrier {
//	return &BasePublishBarrier{
//		scheduler:        scheduler,
//		pubStrategy:      strategy,
//		bufferSize:       scheduler.GetBufferSize(),
//	}
//}

/*
========== single publisher ===============
*/

type SinglePublishBarrier struct {
	BasePublishBarrier
}

func NewSinglePublishBarrier(scheduler IScheduler, strategy IPublisherStrategy) *SinglePublishBarrier {
	return &SinglePublishBarrier{
		BasePublishBarrier{
			bufferSize:  scheduler.GetBufferSize(),
			minSeqCache: 0,
			scheduler:   scheduler,
			pubStrategy: strategy,
		},
	}
}

/*
=========== multi publisher ===========
*/

type MultiPublishBarrier struct {
	BasePublishBarrier
}

func NewMultiPublishBarrier(scheduler IScheduler, strategy IPublisherStrategy) *MultiPublishBarrier {
	return &MultiPublishBarrier{
		BasePublishBarrier{
			bufferSize:  scheduler.GetBufferSize(),
			minSeqCache: 0,
			scheduler:   scheduler,
			pubStrategy: strategy,
		},
	}
}

func (b *MultiPublishBarrier) Commit(seq int64, entity interface{}) {
	b.scheduler.SetEntity(seq, entity)
	b.scheduler.SetAvailableArrayValue(seq)
}

/* subscriber call */

type SubscribeBarrier struct {
	scheduler   IScheduler
	subStrategy ISubscriberStrategy
}

func NewSubscribeBarrier(scheduler IScheduler, strategy ISubscriberStrategy) ISubscribeBarrier {
	return &SubscribeBarrier{
		scheduler:   scheduler,
		subStrategy: strategy,
	}
}

// WaitFor waiting sequence is available
func (b *SubscribeBarrier) WaitFor(seq int64) int64 {
	cursor := b.scheduler.GetCursor()
	availableSeq := cursor.Get()
	for availableSeq < seq || !b.scheduler.IsAvailable(seq) {
		b.subStrategy.Wait()
		availableSeq = cursor.Get()
	}
	return seq
}

func (b *SubscribeBarrier) GetEntity(seq int64) interface{} {
	return b.scheduler.GetEntity(seq)
}
