package uber_ride

import (
	"fmt"
	"sync"
	"time"
)

type Uber struct {
	r        int
	d        int
	pendingR []*User
	pendingD []*User
	userRide map[int]*Ride
	cond     sync.Cond
}

type User struct {
	id int
	ch chan struct{}
	at time.Time
}

type Ride struct {
	riders map[int]int
	status string
}

func New() *Uber {
	uber := &Uber{
		cond:     sync.Cond{L: &sync.Mutex{}},
		pendingR: make([]*User, 0),
		pendingD: make([]*User, 0),
		userRide: make(map[int]*Ride),
	}
	go uber.process()
	return uber
}

func (u *Uber) AddDemocrat(id int, ch chan struct{}) {
	u.cond.L.Lock()
	defer u.cond.L.Unlock()

	u.d++
	u.pendingD = append(u.pendingD, &User{id: id, ch: ch, at: time.Now()})
	u.cond.Signal()
}

func (u *Uber) AddRepublic(id int, ch chan struct{}) {
	u.cond.L.Lock()
	defer u.cond.L.Unlock()

	u.r++
	u.pendingR = append(u.pendingR, &User{id: id, ch: ch, at: time.Now()})
	u.cond.Signal()
}

func (u *Uber) Seated(id int) {
	u.cond.L.Lock()
	defer u.cond.L.Unlock()

	u.userRide[id].riders[id] = 1
}

func (u *Uber) Drive(id int) {
	u.cond.L.Lock()

	ride := u.userRide[id]
	if len(ride.riders) < 4 {
		u.cond.L.Unlock()
		return
	}
	if ride.status != "PENDING" {
		u.cond.L.Unlock()
		return
	}

	ride.status = "STARTED"
	u.cond.L.Unlock()

	fmt.Println("Ride started for users: ", ride.riders, "by user: ", id)
	time.Sleep(10 * time.Second)
	fmt.Println("Ride ended for users: ", ride.riders)
}

func (u *Uber) process() {
	u.cond.L.Lock()
	defer u.cond.L.Unlock()

	for {
		for u.r < 4 && u.d < 4 && (u.d < 2 || u.r < 2) {
			u.cond.Wait()
		}

		rUsers := make([]*User, 0)
		if u.r >= 4 {
			rUsers = u.pendingR[0:4]
		}

		dUsers := make([]*User, 0)
		if u.d >= 4 {
			dUsers = u.pendingD[0:4]
		}

		rdUsers := make([]*User, 0)
		if u.d >= 2 && u.r >= 2 {
			rdUsers = append(u.pendingR[0:2], u.pendingD[0:2]...)
		}

		rWaitTime := meanAtTime(rUsers)
		dWaitTime := meanAtTime(dUsers)
		rdWaitTime := meanAtTime(rdUsers)

		if rWaitTime.After(dWaitTime) {
			if rWaitTime.After(rdWaitTime) || (rWaitTime == rdWaitTime && len(rUsers) > 0) {
				u.createRide(rUsers)
				u.pendingR = u.pendingR[4:]
				u.r -= 4
			} else {
				u.createRide(rdUsers)
				u.pendingR = u.pendingR[2:]
				u.pendingD = u.pendingD[2:]
				u.r -= 2
				u.d -= 2
			}
		} else {
			if dWaitTime.After(rdWaitTime) || (dWaitTime == rdWaitTime && len(dUsers) > 0) {
				u.createRide(dUsers)
				u.pendingD = u.pendingD[4:]
				u.d -= 4
			} else {
				u.createRide(rdUsers)
				u.pendingR = u.pendingR[2:]
				u.pendingD = u.pendingD[2:]
				u.r -= 2
				u.d -= 2
			}
		}
	}
}

func (u *Uber) createRide(users []*User) {
	ride := &Ride{
		riders: make(map[int]int),
		status: "PENDING",
	}

	for _, user := range users {
		select {
		case user.ch <- struct{}{}:
		default:
		}
		u.userRide[user.id] = ride
	}
}

func meanAtTime(users []*User) time.Time {
	var d = time.Time{}
	for _, user := range users {
		d.Add(time.Now().Sub(user.at))
	}
	return d
}
