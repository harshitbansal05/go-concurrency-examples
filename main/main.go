package main

import (
	"fmt"
	"go-lru/barber"
	"go-lru/bathroom"
	"go-lru/blocking_queue_ch"
	"go-lru/blocking_queue_cond"
	"go-lru/blocking_queue_mutex"
	"go-lru/deferred_callback"
	"go-lru/dining"
	"go-lru/mergesort"
	"go-lru/rwlock"
	"go-lru/semaphore"
	"go-lru/sudoku"
	"go-lru/token_bucket_filter"
	"go-lru/uber_ride"
	"math/rand"
	"sync"
	"time"
)

func main() {
	runMergesort()
}

func runMergesort() {
	randomArray := make([]int, 25)
	for i := 0; i < 25; i++ {
		randomArray[i] = rand.Intn(100) + 1
	}
	fmt.Println("Random Array of Integers: ", randomArray)
	an := mergesort.Sort(randomArray)
	fmt.Println("Sorted Array of Integers: ", an)
}

func runBarber() {
	b := barber.New(3)
	var wg sync.WaitGroup
	wg.Add(15)

	for i := 0; i < 10; i++ {
		go func(j int) {
			err := b.Add(j)
			if err != nil {
				fmt.Println(j, ": ", err)
			}
			wg.Done()
		}(i)
	}

	time.Sleep(2 * time.Second)

	for i := 10; i < 15; i++ {
		go func(j int) {
			err := b.Add(j)
			if err != nil {
				fmt.Println(j, ": ", err)
			}
			wg.Done()
		}(i)
		time.Sleep(5 * time.Millisecond)
	}
	wg.Wait()
}

func runDining() {
	d := dining.New()
	var wg sync.WaitGroup
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func(j int) {
			for {
				d.Eat(j)
				time.Sleep(time.Millisecond * 1000)
				d.Contemplate(j)
				time.Sleep(time.Millisecond * 1000)
			}
		}(i)
	}
	wg.Wait()

}

func runUber() {
	u := uber_ride.New()
	var wg sync.WaitGroup
	wg.Add(24)

	for i := 0; i < 10; i++ {
		go func(j int) {
			defer wg.Done()
			ch := make(chan struct{})
			u.AddDemocrat(j, ch)
			select {
			case <-ch:
			}
			u.Seated(j)
			u.Drive(j)
		}(i)
		time.Sleep(time.Millisecond * 50)
	}

	for i := 20; i < 34; i++ {
		go func(j int) {
			defer wg.Done()
			ch := make(chan struct{})
			u.AddRepublic(j, ch)
			select {
			case <-ch:
			}
			u.Seated(j)
			u.Drive(j)
		}(i)
		time.Sleep(time.Millisecond * 20)
	}

	wg.Wait()
}

func runBathroom() {
	b := bathroom.New(1000)
	wg := sync.WaitGroup{}
	wg.Add(5)

	go func() {
		b.AddWomen("Lisa")
		wg.Done()
	}()
	go func() {
		b.AddMen("John")
		wg.Done()
	}()
	go func() {
		b.AddMen("Bob")
		wg.Done()
	}()
	go func() {
		b.AddMen("Anil")
		wg.Done()
	}()
	go func() {
		b.AddMen("Wentao")
		wg.Done()
	}()

	wg.Wait()
}

func runRWLock() {
	rw := rwlock.New()
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		fmt.Println("Attempting to acquire read lock in t3: ", time.Now())
		rw.AcquireRLock()
		fmt.Println("read lock acquired t3: ", time.Now())
	}()

	time.Sleep(1 * time.Second)

	go func() {
		defer wg.Done()
		fmt.Println("Attempting to acquire write lock in t1: ", time.Now())
		rw.AcquireLock()
		fmt.Println("write lock acquired t1: ", time.Now())

		for {
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(3 * time.Second)

	go func() {
		defer wg.Done()
		fmt.Println("Attempting to release read lock in t4: ", time.Now())
		rw.ReleaseRLock()
		fmt.Println("read lock released t4: ", time.Now())
	}()

	time.Sleep(1 * time.Second)

	go func() {
		defer wg.Done()
		fmt.Println("Attempting to acquire write lock in t2: ", time.Now())
		rw.AcquireLock()
		fmt.Println("write lock acquired t: ", time.Now())
	}()

	wg.Wait()
}

func runSemaphore() {
	dc := semaphore.New(1)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < 5; i++ {
			dc.Acquire()
			fmt.Println("Ping ", i)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 5; i++ {
			dc.Release()
			fmt.Println("Pong ", i)
		}
		wg.Done()
	}()

	wg.Wait()
}

func runDeferredCallback() {
	dc := deferred_callback.New()
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(j int) {
			defer wg.Done()
			params := []any{j}
			duration := time.Duration(2*(j+1)) * time.Second
			fmt.Printf("Scheduling function %d to run at time %s\n", j, time.Now().Add(duration))
			dc.Add(j, func(k int) { fmt.Println("Running function ", k) }, duration, params)
		}(i)
	}
	wg.Wait()
	time.Sleep(30 * time.Second)
}

func runTokenBucketLimiterEn() {
	limiter := token_bucket_filter.NewEn(5, 100, 2)
	time.Sleep(4 * time.Second)
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			cb := make(chan struct{})
			err := limiter.GetToken(i, cb)
			if err != nil {
				panic(err)
			} else {
				select {
				case <-cb:
					fmt.Printf("Retrieved token for go-routine %d\n", i)
				}
			}
		}()
	}
	wg.Wait()
}

func runTokenBucketLimiterFifoCond() {
	limiter := token_bucket_filter.NewCond(5, 2)
	time.Sleep(4 * time.Second)
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(time.Second * 2)
			limiter.GetToken(i)
			wg.Done()
		}()
	}
	wg.Wait()
}

func runTokenBucketLimiter2Alter2() {
	limiter := token_bucket_filter.NewAlter2(5)
	time.Sleep(6 * time.Second)
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			limiter.GetToken()
			fmt.Printf("Retrieved token for go-routine %d\n", i)
			wg.Done()
		}()
	}
	wg.Wait()
}

func runTokenBucketLimiter2Alter() {
	limiter := token_bucket_filter.NewAlter(5)
	time.Sleep(6 * time.Second)
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			limiter.GetToken()
			fmt.Printf("Retrieved token for go-routine %d\n", i)
			wg.Done()
		}()
	}
	wg.Wait()
}

func runTokenBucketLimiter2() {
	limiter := token_bucket_filter.New(5)
	time.Sleep(6 * time.Second)
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			limiter.GetToken()
			fmt.Printf("Retrieved token for go-routine %d\n", i)
			wg.Done()
		}()
	}
	wg.Wait()
}

func runTokenBucketLimiter() {
	limiter := token_bucket_filter.New(1)
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			limiter.GetToken()
			fmt.Printf("Retrieved token for go-routine %d\n", i)
			wg.Done()
		}()
	}
	wg.Wait()
}

func runSudoku() {
	board := [][]byte{
		{byte(5), byte(3), byte(0), byte(0), byte(7), byte(0), byte(0), byte(0), byte(0)},
		{byte(6), byte(0), byte(0), byte(1), byte(9), byte(5), byte(0), byte(0), byte(0)},
		{byte(0), byte(9), byte(8), byte(0), byte(0), byte(0), byte(0), byte(6), byte(0)},
		{byte(8), byte(0), byte(0), byte(0), byte(6), byte(0), byte(0), byte(0), byte(3)},
		{byte(4), byte(0), byte(0), byte(8), byte(0), byte(3), byte(0), byte(0), byte(1)},
		{byte(7), byte(0), byte(0), byte(0), byte(2), byte(0), byte(0), byte(0), byte(6)},
		{byte(0), byte(6), byte(0), byte(0), byte(0), byte(0), byte(2), byte(8), byte(0)},
		{byte(0), byte(0), byte(0), byte(4), byte(1), byte(9), byte(0), byte(0), byte(5)},
		{byte(0), byte(0), byte(0), byte(0), byte(8), byte(0), byte(0), byte(7), byte(9)},
	}
	sudoku.SolveSudoku(board)
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			fmt.Printf("%d ", board[i][j])
		}
		fmt.Println()
	}
}

func runBlockingQueueCh() {
	blockingQueue := blocking_queue_ch.New(5)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()

		for i := 1; i <= 50; i++ {
			blockingQueue.Enqueue(i)
		}
		fmt.Println("Completed Go-routine 1")
	}()

	go func() {
		defer wg.Done()

		for i := 1; i <= 25; i++ {
			blockingQueue.Dequeue()
		}
		fmt.Println("Completed Go-routine 2")
	}()

	go func() {
		defer wg.Done()

		for i := 1; i <= 25; i++ {
			blockingQueue.Dequeue()
		}
		fmt.Println("Completed Go-routine 3")
	}()

	wg.Wait()
}

func runBlockingQueueMutex() {
	blockingQueue := blocking_queue_mutex.New(5)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()

		for i := 1; i <= 50; i++ {
			blockingQueue.Enqueue(i)
		}
		fmt.Println("Completed Go-routine 1")
	}()

	go func() {
		defer wg.Done()

		for i := 1; i <= 25; i++ {
			blockingQueue.Dequeue()
		}
		fmt.Println("Completed Go-routine 2")
	}()

	go func() {
		defer wg.Done()

		for i := 1; i <= 25; i++ {
			blockingQueue.Dequeue()
		}
		fmt.Println("Completed Go-routine 3")
	}()

	wg.Wait()
}

func runBlockingQueueCond() {
	blockingQueue := blocking_queue_cond.New(5)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()

		for i := 1; i <= 50; i++ {
			blockingQueue.Enqueue(i)
		}
		fmt.Println("Completed Go-routine 1")
	}()

	go func() {
		defer wg.Done()

		for i := 1; i <= 25; i++ {
			blockingQueue.Dequeue()
		}
		fmt.Println("Completed Go-routine 2")
	}()

	go func() {
		defer wg.Done()

		for i := 1; i <= 25; i++ {
			blockingQueue.Dequeue()
		}
		fmt.Println("Completed Go-routine 3")
	}()

	wg.Wait()
}
