package main

import (
	"sync"
)

/**
以汽车组装为例，汽车生产线上有个阶段是给小汽车装4个轮子，可以把这个阶段任务交给4个人同时去做，这4个人把轮子都装完后，再把汽车移动到生产线下一个阶段。这个过程中，就有任务的分发，和任务结果的收集。其中任务分发是FAN-OUT，任务收集是FAN-IN。
FAN-OUT模式：多个goroutine从同一个通道读取数据，直到该通道关闭。OUT是一种张开的模式，所以又被称为扇出，可以用来分发任务。
FAN-IN模式：1个goroutine从多个通道读取数据，直到这些通道关闭。IN是一种收敛的模式，所以又被称为扇入，用来收集处理的结果。
*/

func producer(num int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < num; i++ {
			out <- i
		}
	}()
	return out
}

func square(inCh <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range inCh {
			out <- n * n
		}
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	collect := func(in <-chan int) {
		defer wg.Done()
		for n := range in {
			out <- n
		}
	}

	wg.Add(len(cs))
	// FAN-IN
	for _, c := range cs {
		go collect(c)
	}

	// 错误方式：直接等待是bug，死锁，因为merge写了out，main却没有读
	// wg.Wait()
	// close(out)

	// 正确方式
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	in := producer(100)

	// FAN-OUT
	c1 := square(in)
	c2 := square(in)
	c3 := square(in)

	// consumer
	//for ret := range merge(c1, c2, c3) {
	//	fmt.Printf("%3d ", ret)
	//	fmt.Println()
	//}
	for _ = range merge(c1, c2, c3) {
	}
}

// time go run fanInAndFanOut.go