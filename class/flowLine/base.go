package main

/**
这是一种原始的流水线模型
1.每个阶段把数据通过channel传递给下一个阶段。
2.每个阶段要创建1个goroutine和1个通道，这个goroutine向里面写数据，函数要返回这个通道。
3.有1个函数来组织流水线，我们例子中是main函数。
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

func main() {
	in := producer(100)
	ch := square(in)

	//for ret := range ch {
	//	fmt.Printf("%3d", ret)
	//	fmt.Println()
	//}
	for _ = range ch {
	}
}



// time go run base.go