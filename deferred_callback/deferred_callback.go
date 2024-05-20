package deferred_callback

import (
	"container/heap"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type JobHeap []Job

type Job struct {
	at         time.Time
	id         int
	function   any
	parameters []any
}

func (h JobHeap) Len() int           { return len(h) }
func (h JobHeap) Less(i, j int) bool { return h[i].at.Before(h[j].at) }
func (h JobHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *JobHeap) Push(x any) {
	*h = append(*h, x.(Job))
}

func (h *JobHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type DeferredCallback struct {
	isPqChanged chan struct{}
	pq          *JobHeap
	nextAt      time.Time
	sync.Mutex
}

func New() *DeferredCallback {
	dc := &DeferredCallback{
		isPqChanged: make(chan struct{}),
		pq:          new(JobHeap),
	}
	heap.Init(dc.pq)
	go dc.process()
	return dc
}

func (dc *DeferredCallback) Add(id int, function any, after time.Duration, parameters []any) {
	dc.Lock()
	heap.Push(dc.pq, Job{
		at:         time.Now().Add(after),
		id:         id,
		function:   function,
		parameters: parameters,
	})
	dc.Unlock()
	dc.isPqChanged <- struct{}{}
}

func (dc *DeferredCallback) process() {
	for {
		after := time.Duration(1<<63 - 1)
		if !dc.nextAt.IsZero() {
			after = dc.nextAt.Sub(time.Now())
		}
		select {
		case <-dc.isPqChanged:
			dc.Lock()
			nextJob := heap.Pop(dc.pq).(Job)
			dc.nextAt = nextJob.at
			heap.Push(dc.pq, nextJob)
			dc.Unlock()
		case <-time.After(after):
			dc.Lock()
			go runJob(heap.Pop(dc.pq).(Job))
			if dc.pq.Len() == 0 {
				dc.nextAt = time.Time{}
			} else {
				nextJob := heap.Pop(dc.pq).(Job)
				dc.nextAt = nextJob.at
				heap.Push(dc.pq, nextJob)
			}
			dc.Unlock()
		}
	}
}

func runJob(j Job) {
	err := callJobFuncWithParams(j.function, j.parameters...)
	if err != nil {
		fmt.Printf("callJobFuncWithParams failed for job %d with error: %s\n", j.id, err)
	} else {
		fmt.Printf("callJobFuncWithParams succeeded for job %d at time %s\n", j.id, time.Now())
	}
}

func callJobFuncWithParams(jobFunc any, params ...any) error {
	if jobFunc == nil {
		return nil
	}
	f := reflect.ValueOf(jobFunc)
	if f.IsZero() {
		return nil
	}
	if len(params) != f.Type().NumIn() {
		return nil
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	returnValues := f.Call(in)
	for _, val := range returnValues {
		i := val.Interface()
		if err, ok := i.(error); ok {
			return err
		}
	}
	return nil
}
