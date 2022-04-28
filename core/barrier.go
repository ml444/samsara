package core

import (
	"errors"
	"fmt"
)

type IPublishBarrier interface {
	Next() int64
	NextN(n int64) int64
	TryNext() (int64, error)
	TryNextN(n int64) (int64, error)
	Commit(seq int64, entity IEntity)
}
type ISubscribeBarrier interface {
	WaitFor(sequence int64) int64
	GetEntity(sequence int64) IEntity
}

/* single publisher call */

type SinglePublishBarrier struct {
	scheduler   IScheduler
	pubStrategy IPublisherStrategy
	//subStrategy ISubscriberStrategy
}

func NewSinglePublishBarrier(scheduler IScheduler, strategy IPublisherStrategy) *SinglePublishBarrier {
	return &SinglePublishBarrier{
		scheduler:   scheduler,
		pubStrategy: strategy,
	}
}

func (b *SinglePublishBarrier) Next() int64 {
	if b.scheduler.GetConsumerSequences() == nil {
		panic("consumerSequences must be set before claiming sequences")
	}
	cursor := b.scheduler.GetCursor()
	for {
		minSeq := b.scheduler.GetMinConsumerSequence()
		if v := cursor.Get(); v-b.scheduler.GetBufferSize() < minSeq {
			ok := cursor.CompareAndSwap(v, v+1)
			if ok {
				return v + 1
			} else {
				b.pubStrategy.Wait()
			}

		} else {
			b.pubStrategy.Wait()
		}
	}

	//return cursor.Add(1)
}

func (b *SinglePublishBarrier) NextN(n int64) int64 {
	if b.scheduler.GetConsumerSequences() == nil {
		panic("consumerSequences must be set before claiming sequences")
	}
	if n < 1 {
		panic("Required capacity must be greater than 0")
	}

	for {
		minSeq := b.scheduler.GetMinConsumerSequence()
		if b.scheduler.GetCursor().Get()-b.scheduler.GetBufferSize()-n < minSeq {
			break
		} else {
			b.pubStrategy.Wait()
		}
	}
	return b.scheduler.GetCursor().Add(n) - n
}

func (b *SinglePublishBarrier) TryNext() (int64, error) {
	return b.TryNextN(1)
}

func (b *SinglePublishBarrier) TryNextN(n int64) (int64, error) {
	return 0, nil
}

func (b *SinglePublishBarrier) Commit(seq int64, entity IEntity) {
	//fmt.Println("===>", seq)
	//index := seq % b.scheduler.GetBufferSize()
	b.scheduler.SetEntity(seq, entity)
}

/* multi publisher call */

type MultiPublishBarrier struct {
	scheduler        IScheduler
	pubStrategy      IPublisherStrategy
	minSequenceCache *Sequence
	bufferSize       int64
	//subStrategy ISubscriberStrategy
}

func NewMultiPublishBarrier(scheduler IScheduler, strategy IPublisherStrategy) *MultiPublishBarrier {

	return &MultiPublishBarrier{
		scheduler:        scheduler,
		pubStrategy:      strategy,
		bufferSize:       scheduler.GetBufferSize(),
		minSequenceCache: NewSequence(),
	}
}

func (b *MultiPublishBarrier) Next() int64 {
	return b.NextN(1)
}

func (b *MultiPublishBarrier) NextN(n int64) int64 {
	if b.scheduler.GetConsumerSequences() == nil {
		panic("consumerSequences must be set before claiming sequences")
	}
	if n < 1 || n > b.bufferSize {
		panic(fmt.Sprintf("must be 0 < n < %d", b.bufferSize))
	}

	cursor := b.scheduler.GetCursor()
	cursorValue := cursor.Add(n) - n
	nextSeq := cursorValue + n
	wrapPoint := nextSeq - b.bufferSize
	minSeq := b.minSequenceCache.Get()
	if wrapPoint > minSeq || minSeq > cursorValue {
		var newMinSeq int64
		for newMinSeq = b.scheduler.GetMinConsumerSequence(); wrapPoint > newMinSeq; {
			b.pubStrategy.Wait()
		}
		b.minSequenceCache.Set(newMinSeq)
	}
	return nextSeq
	//return b.scheduler.GetCursor().Add(n) - n
}

func (b *MultiPublishBarrier) TryNext() (int64, error) {
	return b.TryNextN(1)
}

func (b *MultiPublishBarrier) TryNextN(n int64) (int64, error) {
	if b.scheduler.GetConsumerSequences() == nil {
		panic("consumerSequences must be set before claiming sequences")
	}
	if n < 1 || n > b.bufferSize {
		panic(fmt.Sprintf("must be 0 < n < %d", b.bufferSize))
	}
	for {
		cursor := b.scheduler.GetCursor()
		cursorValue := cursor.Get()
		nextSeq := cursorValue + n
		wrapPoint := nextSeq - b.bufferSize
		minSeq := b.minSequenceCache.Get()
		if wrapPoint > minSeq || minSeq > cursorValue {
			newMinSeq := b.scheduler.GetMinConsumerSequence()
			b.minSequenceCache.Set(newMinSeq)
			if wrapPoint > newMinSeq {
				return 0, errors.New("insufficient capacity")
			}
			if cursor.CompareAndSwap(cursorValue, nextSeq) {
				return nextSeq, nil
			}
		}
	}
}

func (b *MultiPublishBarrier) Commit(seq int64, entity IEntity) {
	//fmt.Println("===>", seq)
	//index := seq % b.scheduler.GetBufferSize()
	b.scheduler.SetEntity(seq, entity)
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

// WaitFor
func (b *SubscribeBarrier) WaitFor(seq int64) int64 {
	for {
		availableSeq := b.scheduler.GetCursor().Get()
		if availableSeq < seq {
			b.subStrategy.Wait()
		} else {
			break
		}
	}
	return seq
}

func (b *SubscribeBarrier) GetEntity(seq int64) IEntity {
	return b.scheduler.GetEntity(seq)
}
