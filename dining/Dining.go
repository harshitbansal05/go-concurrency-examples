package dining

import "fmt"

type Dining struct {
	forks []chan struct{}
}

func New() *Dining {
	d := Dining{
		forks: make([]chan struct{}, 0),
	}
	for i := 0; i < 5; i++ {
		d.forks = append(d.forks, make(chan struct{}, 1))
		d.forks[i] <- struct{}{}
	}
	return &d
}

func (d *Dining) Eat(id int) {
	a, b := (id+4)%5, id
	for {
		select {
		case <-d.forks[a]:
			select {
			case <-d.forks[b]:
				fmt.Printf("Philosopher %d started eating.\n", id)
				return
			default:
				d.forks[a] <- struct{}{}
			}
		case <-d.forks[b]:
			select {
			case <-d.forks[a]:
				fmt.Printf("Philosopher %d started eating.\n", id)
				return
			default:
				d.forks[b] <- struct{}{}
			}
		}
	}
}

func (d *Dining) Contemplate(id int) {
	a, b := (id+4)%5, id
	d.forks[a] <- struct{}{}
	d.forks[b] <- struct{}{}
	fmt.Printf("Philosopher %d stopped eating.\n", id)
}
