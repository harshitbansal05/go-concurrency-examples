package bathroom

import (
	"fmt"
	"sync"
	"time"
)

type Bathroom struct {
	m    int
	w    int
	max  int
	cond sync.Cond
}

func New(max int) *Bathroom {
	return &Bathroom{max: max, cond: sync.Cond{L: &sync.Mutex{}}}
}

func (b *Bathroom) AddMen(name string) {
	b.cond.L.Lock()

	for b.w > 0 || b.w+b.m == b.max {
		b.cond.Wait()
	}
	b.m++
	b.cond.L.Unlock()

	fmt.Printf("%s is using bathroom. Current employees in bathroom = %d\n", name, b.m)
	time.Sleep(time.Second * 10)
	fmt.Println(name + " is done using bathroom")

	b.cond.L.Lock()
	b.m--
	b.cond.L.Unlock()
	b.cond.Broadcast()
}

func (b *Bathroom) AddWomen(name string) {
	b.cond.L.Lock()

	for b.m > 0 || b.w+b.m == b.max {
		b.cond.Wait()
	}
	b.w++
	b.cond.L.Unlock()

	fmt.Printf("%s is using bathroom. Current employees in bathroom = %d\n", name, b.w)
	time.Sleep(time.Second * 10)
	fmt.Println(name + " is done using bathroom")

	b.cond.L.Lock()
	b.w--
	b.cond.L.Unlock()
	b.cond.Broadcast()
}
