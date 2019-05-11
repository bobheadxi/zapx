package internal

import (
	"sync"
)

// Buffer is a wrapper around a byte slice, integrated with a pool
type Buffer struct {
	b []byte
	p Pool
}

// NewBuffer creates a new buffer with a default, fixed size (1024)
func NewBuffer() interface{} { return &Buffer{b: make([]byte, 1024, 1024)} }

// Bytes returns the underlying byte slice
func (b *Buffer) Bytes() []byte { return b.b }

// Reset clears the underlying byte slice, but does not allocate a new array
func (b *Buffer) Reset() {
	b.b = b.b[:0]    // empty the slice
	b.b = b.b[:1024] // grow it back to capacity
}

// Free returns the Buffer to the associated pool for reuse
func (b *Buffer) Free() { b.p.put(b) }

// Pool wraps sync.Pool with typesafe returns for *Buffer
type Pool struct {
	p *sync.Pool
}

// NewPool instantiates a new pool
func NewPool() Pool {
	return Pool{p: &sync.Pool{
		New: NewBuffer,
	}}
}

// Get retrieves an available *Buffer from the pool, or allocates a new one
func (p Pool) Get() *Buffer {
	buf := p.p.Get().(*Buffer)
	buf.Reset()
	buf.p = p
	return buf
}

func (p Pool) put(buf *Buffer) { p.p.Put(buf) }
