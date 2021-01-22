package main

import (
	"context"
	"fmt"
	"sync"
)

//分离模式
/**
----------------------分割线--------------------
*/
// 一次性任务 参考$GOROOT/src/net/dial.go
type Dialer struct {
	Cancel <-chan struct{}
}

func (d *Dialer) DialContext(ctx context.Context) {
	// ...
	if oldCancel := d.Cancel; oldCancel != nil {
		subCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		go func() {
			select {
			case <-oldCancel:
				cancel()
			case <-subCtx.Done():
			}
		}()
		ctx = subCtx
	}
	// ...
}

/**
----------------------分割线--------------------
*/
// 常驻进程任务 参考$GOROOT/src/runtime/mgc.go
//func gcBgMarkStartWorkers() {
//	// Background marking is performed by per-P G's. Ensure that
//	// each P has a background GC G.
//	for _, p := range allp {
//		if p.gcBgMarkWorker == 0 {
//			go gcBgMarkWorker(p) // 每个P创建一个goroutine，以运行gcBgMarkWorker
//			notetsleepg(&work.bgMarkReady, -1)
//			noteclear(&work.bgMarkReady)
//		}
//	}
//}
//
//func gcBgMarkWorker(_p_ *p) {
//	gp := getg()
//	... ...
//	for {
//		// 处理GC事宜
//		... ...
//	}
//}

/**
----------------------分割线--------------------
*/
// join模式 创建者不仅仅要等待goroutine的退出，还要知道结束状态-
func add(a, b int) int {
	return a + b
}

func Do(f func(a, b int) int, a, b, n int) chan int {
	c := make(chan int)

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			r := f(a, b)
			fmt.Println(r)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		c <- 100
	}()

	go func() {

	}()

	return c
}

func main() {
	c := Do(add, 1, 5, 5)
	fmt.Println(<-c)
}
