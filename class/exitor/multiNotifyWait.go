package main

import (
	"fmt"
	"sync"
	"time"
)

//通知并等待多个goroutine退出
func worker(x int) {
	fmt.Printf("worker %d sleeping...", x)
	time.Sleep(time.Second * time.Duration(x))
}

func Do(f func(a int), n int) chan int {
	quit := make(chan int)
	job := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("worker-%d", i)
			for {
				j, ok := <-job
				if !ok {
					fmt.Println(name, "done")
					return
				}
				f(j)
			}
		}(i)
	}

	go func() {
		<-quit
		close(job)
		wg.Wait()
		quit <- 200
	}()

	return quit
}

func main() {
	quit := Do(worker, 5)
	fmt.Println("func Work...")
	quit <- 1
	timer := time.NewTimer(time.Second * 100)
	defer timer.Stop()
	select {
	case status := <-quit:
		fmt.Println(status)
	case <-timer.C:
		fmt.Println("等待...")
	}
}
