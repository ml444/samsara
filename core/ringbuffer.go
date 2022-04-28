package core

import (
	"golang.org/x/sys/cpu"
)

type RingBuffer struct {
	_ cpu.CacheLinePad
	//Scheduler
	indexMask int64
	size      int64
	entities  []IEntity
}

func NewRingBuffer(size int64, eventFactory IEntityFactory) *RingBuffer {
	if size <= 0 {
		panic("bufferSize must be gather than 0")
	}
	if size%2 != 0 {
		panic("bufferSize must be a power of 2")
	}
	var rb = &RingBuffer{
		indexMask: size - 1,
		size:      size,
		entities:  make([]IEntity, size),
	}
	//rb.fill(eventFactory)
	return rb
}

func (rb *RingBuffer) Size() int64 {
	return rb.size
}

func (rb *RingBuffer) fill(eventFactory IEntityFactory) {
	for i := 0; i < len(rb.entities); i++ {
		rb.entities[i] = eventFactory.NewEntity()
	}
}

func (rb *RingBuffer) GetEntity(seq int64) IEntity {
	return rb.entities[int(seq&rb.indexMask)]
}

func (rb *RingBuffer) SetEntity(seq int64, entity IEntity) {
	rb.entities[int(seq&rb.indexMask)] = entity
}

