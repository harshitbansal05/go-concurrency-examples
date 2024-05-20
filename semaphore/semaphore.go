package semaphore

import "sync"

type Semaphore struct {
	size int
	used int
	cond sync.Cond
}

func New(size int) *Semaphore {
	return &Semaphore{
		size: size,
		cond: sync.Cond{L: &sync.Mutex{}},
	}
}

func NewWithUsed(size, used int) *Semaphore {
	return &Semaphore{
		size: size,
		used: used,
		cond: sync.Cond{L: &sync.Mutex{}},
	}
}

func (s *Semaphore) Acquire() {
	s.cond.L.Lock()
	for s.used == s.size {
		s.cond.Wait()
	}
	s.used++
	s.cond.L.Unlock()
	s.cond.Signal()
}

func (s *Semaphore) Release() {
	s.cond.L.Lock()
	for s.used == 0 {
		s.cond.Wait()
	}
	s.used--
	s.cond.L.Unlock()
	s.cond.Signal()
}
