package main

import "sync"

var (
	wg sync.WaitGroup
	l sync.Mutex
	counter int
)

func main() {
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Lock()
			counter++
			l.Unlock()
		}()
	}

	wg.Wait()
	println(counter)
}
