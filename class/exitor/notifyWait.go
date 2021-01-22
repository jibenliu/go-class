package main

import (
	"fmt"
	"time"
)

func add(a, b int) int {
	return a + b
}

func Do(f func(a, b int) int, a, b int) chan int {
	quit := make(chan int)

	go func() {
		var job chan string
		for {
			select {
			case x := <-job:
				f(a, b)
				fmt.Println(x)
			case y := <-quit:
				quit <- y
			}
		}
	}()

	return quit
}

// 有些场景，goroutine的创建者需要主动通知那些新goroutine退出.
func main() {
	c := Do(add, 1, 5)
	fmt.Println("开始干活")
	time.Sleep(1 * time.Second)
	c <- 0
	timer := time.NewTimer(time.Second * 10)
	defer timer.Stop()
	select {
	case status := <-c:
		fmt.Println(status)
	case <-timer.C:
		fmt.Println("等待...")
	}
}
