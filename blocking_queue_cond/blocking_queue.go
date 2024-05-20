package blocking_queue_cond

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
	cond     *sync.Cond
}

func New(capacity int) *BlockingQueue {
	return &BlockingQueue{
		capacity: capacity,
		arr:      make([]int, capacity),
		cond:     sync.NewCond(&sync.Mutex{}),
	}
}

func (bq *BlockingQueue) Enqueue(i int) {
	bq.cond.L.Lock()
	for bq.size == bq.capacity {
		bq.cond.Wait()
	}

	bq.arr[bq.head] = i
	bq.head = (bq.head + 1) % bq.capacity
	bq.size++
	fmt.Println("Enqueue ", i)
	bq.cond.L.Unlock()
	bq.cond.Broadcast()
}

func (bq *BlockingQueue) Dequeue() (i int) {
	bq.cond.L.Lock()
	for bq.size == 0 {
		bq.cond.Wait()
	}

	i = bq.arr[bq.tail]
	bq.tail = (bq.tail + 1) % bq.capacity
	bq.size--
	fmt.Println("Dequeue ", i)
	bq.cond.L.Unlock()
	bq.cond.Broadcast()
	return
}
