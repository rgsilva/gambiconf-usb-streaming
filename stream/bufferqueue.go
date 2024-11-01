package stream

import (
	"fmt"
	"sync"
	"time"
)

type BufferQueue struct {
	buffer []byte
	lock   sync.Mutex
}

// NewBufferQueue initializes and returns a new BufferQueue.
func NewBufferQueue() *BufferQueue {
	return &BufferQueue{
		buffer: make([]byte, 0),
	}
}

func (bq *BufferQueue) Has(size int) bool {
	bq.lock.Lock()
	defer bq.lock.Unlock()

	return len(bq.buffer) > size
}

// Put adds bytes to the buffer in a thread-safe manner.
func (bq *BufferQueue) Put(bytes []byte) {
	bq.lock.Lock()
	defer bq.lock.Unlock()

	bq.buffer = append(bq.buffer, bytes...)
	fmt.Printf("[BQ] Buffer size: %d\n", len(bq.buffer))
}

// Get retrieves `size` bytes from the buffer in a thread-safe manner.
func (bq *BufferQueue) Get(size int) []byte {
	bq.lock.Lock()
	defer bq.lock.Unlock()

	before := len(bq.buffer)
	if size > len(bq.buffer) {
		size = len(bq.buffer)
	}
	buf := bq.buffer[:size]
	bq.buffer = bq.buffer[size:]
	after := len(bq.buffer)

	fmt.Printf("[BQ] New buffer size: %d -> %d\n", before, after)
	return buf
}

// Empty checks if the buffer is empty in a thread-safe manner.
func (bq *BufferQueue) Empty() bool {
	bq.lock.Lock()
	defer bq.lock.Unlock()

	return len(bq.buffer) == 0
}

// Size returns the current size of the buffer in a thread-safe manner.
func (bq *BufferQueue) Size() int {
	bq.lock.Lock()
	defer bq.lock.Unlock()

	return len(bq.buffer)
}

// WaitFor waits until the buffer has at least `size` bytes and then signals via a channel.
func (bq *BufferQueue) WaitFor(size int) <-chan bool {
	done := make(chan bool)

	go func() {
		for {
			if bq.Size() >= size {
				done <- true
				return
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	return done
}
