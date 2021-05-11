package main

import (
	"fmt"
	"time"
)

func timeCost1() func() {
	startT := time.Now()
	return func() {
		tc := time.Since(startT)
		fmt.Printf("time cost = %v\n", tc)
	}
}

func sum2(n int) int {
	// 注解式调用,和业务完全解耦
	defer timeCost1()() // 注意，是对 timeCost() 返回的函数进行延迟调用，因此需要加两对小括号
	total := 0
	for i := 0; i <= n; i++ {
		total += i
	}
	return total
}

func main() {
	count := sum2(100)
	fmt.Printf("count = %v\n", count)
}
