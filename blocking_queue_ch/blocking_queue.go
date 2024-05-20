package blocking_queue_ch

import (
	"fmt"
	"sync"
)

type BlockingQueue struct {
	capacity  int
	queue     chan struct{}
	readQueue chan struct{}
	arr       []int
	head      int
	tail      int
	size      int
	sync.Mutex
}

func New(capacity int) *BlockingQueue {
	return &BlockingQueue{
		capacity:  capacity,
		queue:     make(chan struct{}, capacity),
		readQueue: make(chan struct{}, capacity),
		arr:       make([]int, capacity),
	}
}

func (bq *BlockingQueue) Enqueue(i int) {
	select {
	case bq.queue <- struct{}{}:
		bq.Lock()
	}

	bq.arr[bq.head] = i
	bq.head = (bq.head + 1) % bq.capacity
	bq.size++
	bq.readQueue <- struct{}{}
	fmt.Println("Enqueue ", i)
	bq.Unlock()
}

func (bq *BlockingQueue) Dequeue() (i int) {
	select {
	case <-bq.readQueue:
		bq.Lock()
	}

	i = bq.arr[bq.tail]
	bq.tail = (bq.tail + 1) % bq.capacity
	bq.size--
	<-bq.queue
	fmt.Println("Dequeue ", i)
	bq.Unlock()

	return
}
