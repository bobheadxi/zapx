package pool

import (
	"sync"
)

// ByteBuffer is a wrapper around a byte slice, integrated with a pool
type ByteBuffer struct {
	b []byte
	p ByteBufferPool
}

// NewByteBuffer creates a new buffer with a default, fixed size (1024)
func NewByteBuffer() interface{} { return &ByteBuffer{b: make([]byte, 1024, 1024)} }

// Bytes returns the underlying byte slice
func (b *ByteBuffer) Bytes() []byte { return b.b }

// Reset clears the underlying byte slice, but does not allocate a new array
func (b *ByteBuffer) Reset() {
	b.b = b.b[:0]    // empty the slice
	b.b = b.b[:1024] // grow it back to capacity
}

// Free returns the Buffer to the associated pool for reuse
func (b *ByteBuffer) Free() { b.p.put(b) }

// ByteBufferPool wraps sync.Pool with typesafe returns for *Buffer
type ByteBufferPool struct{ p *sync.Pool }

// NewByteBufferPool instantiates a new pool
func NewByteBufferPool() ByteBufferPool {
	return ByteBufferPool{&sync.Pool{
		New: NewByteBuffer,
	}}
}

// Get retrieves an available *Buffer from the pool, or allocates a new one
func (p ByteBufferPool) Get() *ByteBuffer {
	buf := p.p.Get().(*ByteBuffer)
	buf.Reset()
	buf.p = p
	return buf
}

func (p ByteBufferPool) put(buf *ByteBuffer) { p.p.Put(buf) }
