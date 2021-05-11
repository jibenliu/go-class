package main

import (
	"fmt"
	"time"
)

func sum(n int) int {
	startT := time.Now()
	total := 0
	for i := 0; i <= n; i++ {
		total += i
	}

	tc := time.Since(startT)
	fmt.Printf("time cost = %v\n", tc)
	return total
}

func main() {
	count := sum(100)
	fmt.Printf("count = %v\n", count)
}
