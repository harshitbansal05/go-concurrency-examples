package blocking_queue_mutex

import (
	"fmt"
	"sync"
)

type BlockingQueue struct {
	capacity int
	arr      []int
	head     int
	tail     int
	size     int
	sync.Mutex
}

func New(capacity int) *BlockingQueue {
	return &BlockingQueue{
		capacity: capacity,
		arr:      make([]int, capacity),
	}
}

func (bq *BlockingQueue) Enqueue(i int) {
	bq.Lock()
	for bq.size == bq.capacity {
		bq.Unlock()

		bq.Lock()
	}

	bq.arr[bq.head] = i
	bq.head = (bq.head + 1) % bq.capacity
	bq.size++
	fmt.Println("Enqueue ", i)
	bq.Unlock()
}

func (bq *BlockingQueue) Dequeue() (i int) {
	bq.Lock()
	for bq.size == 0 {
		bq.Unlock()

		bq.Lock()
	}

	i = bq.arr[bq.tail]
	bq.tail = (bq.tail + 1) % bq.capacity
	bq.size--
	fmt.Println("Dequeue ", i)
	bq.Unlock()
	return
}
