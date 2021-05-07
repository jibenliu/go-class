package goroutinePool

import (
	"fmt"
	"testing"
)

func TestPool_GetWorker(t *testing.T) {
	var PrintTest = func(s string) error {
		fmt.Println(s)
		return nil
	}
	Pool, err := NewPool(20, 5)
	if err != nil {
		return
	}
	i := 0
	for i < 50 {
		err = Pool.Submit(PrintTest, "并发测试！")
		if err != nil {
			fmt.Println(err)
		}
		i++
	}
}
