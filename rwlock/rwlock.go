package rwlock

import "sync"

type RWLock struct {
	readers int
	isWrite bool
	cond    *sync.Cond
}

func New() *RWLock {
	return &RWLock{
		readers: 0,
		isWrite: false,
		cond:    sync.NewCond(new(sync.Mutex)),
	}
}

func (rl *RWLock) AcquireRLock() {
	rl.cond.L.Lock()
	defer rl.cond.L.Unlock()

	for rl.isWrite {
		rl.cond.Wait()
	}
	rl.readers++
}

func (rl *RWLock) ReleaseRLock() {
	rl.cond.L.Lock()
	defer rl.cond.L.Unlock()

	rl.readers--
	rl.cond.Signal()
}

func (rl *RWLock) AcquireLock() {
	rl.cond.L.Lock()
	defer rl.cond.L.Unlock()

	for rl.isWrite || rl.readers != 0 {
		rl.cond.Wait()
	}
	rl.isWrite = true
}

func (rl *RWLock) ReleaseLock() {
	rl.cond.L.Lock()
	defer rl.cond.L.Unlock()

	rl.isWrite = true
	rl.cond.Signal()
}
