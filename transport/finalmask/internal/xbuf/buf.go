// Package xbuf is the subset of Xray's common/buf used by FinalMask.
package xbuf

import "sync"

// Size of a regular buffer, as in Xray.
const Size = 8192

var pool = sync.Pool{New: func() any { return make([]byte, Size) }}

type Buffer struct {
	v     []byte
	start int32
	end   int32
}

func New() *Buffer {
	return &Buffer{v: pool.Get().([]byte)}
}

func (b *Buffer) Release() {
	if b == nil || b.v == nil {
		return
	}
	pool.Put(b.v)
	b.v = nil
}

func (b *Buffer) Bytes() []byte {
	return b.v[b.start:b.end]
}

func (b *Buffer) Len() int32 {
	return b.end - b.start
}

// Resize cuts the buffer at the given position, as in Xray.
func (b *Buffer) Resize(from, to int32) {
	oldEnd := b.end
	if from < 0 {
		from += b.Len()
	}
	if to < 0 {
		to += b.Len()
	}
	if to < from {
		panic("Invalid slice")
	}
	b.end = b.start + to
	b.start += from
	if b.end > oldEnd {
		clear(b.v[oldEnd:b.end])
	}
}
