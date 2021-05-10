package main

import (
	"sync"
	"time"
)

// Lock try lock
type Lock struct {
	c chan struct{}
}

// NewLock generate a try lock
func NewLock() Lock {
	var l Lock
	l.c = make(chan struct{}, 1)
	l.c <- struct{}{}
	return l
}

// Lock try lock, return lock result
func (l Lock) Lock() bool {
	lockResult := false
	select {
	case <-l.c:
		lockResult = true
	default:
	}
	return lockResult
}

// Unlock , Unlock the try lock
func (l Lock) Unlock() {
	l.c <- struct{}{}
}

var newCounter int

func main() {
	var l = NewLock()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(20 * time.Millisecond) //非耗时情况下,加锁失败的概率很低
			defer wg.Done()
			if !l.Lock() {
				println("lock failed")
				return
			}
			newCounter++
			println("current counter", newCounter)
			l.Unlock()
		}()
	}
	wg.Wait()
}
