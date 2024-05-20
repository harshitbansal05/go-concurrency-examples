package blocking_queue_semaphore

import (
	"fmt"
	"go-lru/semaphore"
	"sync"
)

type BlockingQueue struct {
	capacity int
	arr      []int
	head     int
	tail     int
	size     int
	cm       *semaphore.Semaphore
	dm       *semaphore.Semaphore
	sync.Mutex
}

func New(capacity int) *BlockingQueue {
	return &BlockingQueue{
		capacity: capacity,
		arr:      make([]int, capacity),
		cm:       semaphore.New(capacity),
		dm:       semaphore.NewWithUsed(capacity, capacity),
	}
}

func (bq *BlockingQueue) Enqueue(i int) {
	bq.cm.Acquire()
	bq.Lock()

	bq.arr[bq.head] = i
	bq.head = (bq.head + 1) % bq.capacity
	bq.size++
	fmt.Println("Enqueue ", i)
	bq.Unlock()
	bq.dm.Release()
}

func (bq *BlockingQueue) Dequeue() (i int) {
	bq.dm.Acquire()
	bq.Lock()

	i = bq.arr[bq.tail]
	bq.tail = (bq.tail + 1) % bq.capacity
	bq.size--
	fmt.Println("Dequeue ", i)
	bq.Unlock()
	bq.cm.Release()
	return
}
