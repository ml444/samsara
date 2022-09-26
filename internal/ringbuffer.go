package internal

import (
	"golang.org/x/sys/cpu"
)

type RingBuffer struct {
	_         cpu.CacheLinePad
	indexMask int64
	size      int64
	entities  []interface{}
}

func NewRingBuffer(size int64) *RingBuffer {
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
		entities:  make([]interface{}, size),
	}
	return rb
}

func (rb *RingBuffer) Size() int64 {
	return rb.size
}

func (rb *RingBuffer) IndexMask() int64 {
	return rb.indexMask
}

func (rb *RingBuffer) GetEntity(seq int64) interface{} {
	return rb.entities[int(seq&rb.indexMask)]
}

func (rb *RingBuffer) SetEntity(seq int64, entity interface{}) {
	rb.entities[int(seq&rb.indexMask)] = entity
}
