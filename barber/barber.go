package barber

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Barber struct {
	n    int
	ch   []chan struct{}
	c    int
	h    int
	t    int
	cond sync.Cond
}

func New(n int) *Barber {
	b := &Barber{
		n:    n,
		ch:   make([]chan struct{}, n),
		cond: sync.Cond{L: &sync.Mutex{}},
	}
	go b.process()
	return b
}

func (b *Barber) Add(id int) error {
	b.cond.L.Lock()

	if b.c == b.n {
		b.cond.L.Unlock()
		return errors.New("barber is full")
	}
	cb := make(chan struct{})
	b.ch[b.t] = cb
	b.t = (b.t + 1) % b.n
	b.c++
	b.cond.L.Unlock()
	b.cond.Signal()
	fmt.Println("User ", id, " seated at barber at time ", time.Now())

	select {
	case <-cb:
		fmt.Println("Barbing finished for user ", id, " at time ", time.Now())
		return nil
	}
}

func (b *Barber) process() {
	for {
		b.cond.L.Lock()
		for b.c == 0 {
			b.cond.Wait()
		}
		cb := b.ch[b.h]
		b.h = (b.h + 1) % b.n
		b.c--
		b.cond.L.Unlock()

		time.Sleep(time.Second * 1)
		cb <- struct{}{}
	}
}
