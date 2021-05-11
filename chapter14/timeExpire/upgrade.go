package main

import (
	"fmt"
	"time"
)

func timeCost(start time.Time) {
	tc := time.Since(start)
	fmt.Printf("time cost = %v\n", tc)
}

func sum1(n int) int {
	defer timeCost(time.Now()) // 增加了一次timeCost调用,耗时增加
	total:=0
	for i := 0; i <= n; i++ {
		total += i
	}
	return total
}

func main()  {
	count := sum1(100)
	fmt.Printf("count = %v\n", count)
}
