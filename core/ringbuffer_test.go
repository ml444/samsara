package core

import (
	"golang.org/x/sys/cpu"
	"reflect"
	"testing"
)

func TestSliceMake(t *testing.T) {
	var l1 []IEntity
	t.Log("===>", len(l1))
	l2 := make([]IEntity, 1024)
	t.Log("===>", len(l2))
	t.Log(l2[0])
}

func TestNewRingBuffer(t *testing.T) {
	type args struct {
		size         int64
		eventFactory IEntityFactory
	}
	tests := []struct {
		name string
		args args
		want *RingBuffer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRingBuffer(tt.args.size, tt.args.eventFactory); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRingBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingBuffer_GetEntity(t *testing.T) {
	type fields struct {
		indexMask int64
		size      int64
		entities  []IEntity
	}
	type args struct {
		seq int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   IEntity
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &RingBuffer{
				indexMask: tt.fields.indexMask,
				size:      tt.fields.size,
				entities:  tt.fields.entities,
			}
			if got := rb.GetEntity(tt.args.seq); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingBuffer_SetEntity(t *testing.T) {
	type fields struct {
		indexMask int64
		size      int64
		entities  []IEntity
	}
	type args struct {
		seq    int64
		entity IEntity
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &RingBuffer{
				indexMask: tt.fields.indexMask,
				size:      tt.fields.size,
				entities:  tt.fields.entities,
			}
			rb.SetEntity(tt.args.seq, tt.args.entity)
		})
	}
}

func TestRingBuffer_Size(t *testing.T) {
	type fields struct {
		indexMask int64
		size      int64
		entities  []IEntity
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := &RingBuffer{
				indexMask: tt.fields.indexMask,
				size:      tt.fields.size,
				entities:  tt.fields.entities,
			}
			if got := rb.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingBuffer_fill(t *testing.T) {
	type fields struct {
		_         cpu.CacheLinePad
		indexMask int64
		size      int64
		entities  []IEntity
	}
	type args struct {
		eventFactory IEntityFactory
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//rb := &RingBuffer{
			//	indexMask: tt.fields.indexMask,
			//	size:      tt.fields.size,
			//	entities:  tt.fields.entities,
			//}
			//rb.fill()
		})
	}
}
