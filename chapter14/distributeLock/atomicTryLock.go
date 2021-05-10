package main

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var counter int32
	var totalCounter int32
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		//rand.Seed(time.Now().Unix())
		go func() {
			defer wg.Done()
			t := rand.Int63n(100)
			time.Sleep(time.Millisecond * time.Duration(t)) //非耗时情况下,加锁失败的概率很低,随机休眠时间
			// method 1
			//newCounter := atomic.AddInt32(&counter, 1)
			//println("increase success, current counter", newCounter)

			// method 2
			for {
				old := atomic.LoadInt32(&counter)
				ok := atomic.CompareAndSwapInt32(&counter, old, old+1)
				atomic.AddInt32(&totalCounter, 1)
				if !ok {
					println("increase failed, current counter ", atomic.LoadInt32(&counter))
					continue
				} else {
					println("increase success, current counter", atomic.LoadInt32(&counter))
					return
				}
			}
		}()
	}

	wg.Wait()
	println("total counter is ", totalCounter)
}
