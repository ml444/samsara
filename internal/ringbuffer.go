package internal

import (
	"github.com/ml444/samsara/entity"
	"golang.org/x/sys/cpu"
)

type RingBuffer struct {
	_ cpu.CacheLinePad
	//Scheduler
	indexMask int64
	size      int64
	entities  []entity.IEntity
}

func NewRingBuffer(size int64, eventFactory entity.IEntityFactory) *RingBuffer {
	if size <= 0 {
		panic("bufferSize must be gather than 0")
	}
	// 2^n
	if size&(size-1) != 0 {
		panic("bufferSize must be nth power of 2")
	}
	var rb = &RingBuffer{
		indexMask: size - 1,
		size:      size,
		entities:  make([]entity.IEntity, size),
	}
	//rb.fill(eventFactory)
	return rb
}

func (rb *RingBuffer) Size() int64 {
	return rb.size
}

func (rb *RingBuffer) IndexMask() int64 {
	return rb.indexMask
}
func (rb *RingBuffer) fill(eventFactory entity.IEntityFactory) {
	for i := 0; i < len(rb.entities); i++ {
		rb.entities[i] = eventFactory.NewEntity()
	}
}

func (rb *RingBuffer) GetEntity(seq int64) entity.IEntity {
	return rb.entities[int(seq&rb.indexMask)]
}

func (rb *RingBuffer) SetEntity(seq int64, entity entity.IEntity) {
	rb.entities[int(seq&rb.indexMask)] = entity
}
