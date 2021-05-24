package main

import (
	"context"
	"fmt"
	"time"
)

func AsyncCall1() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*800)
	defer cancel()
	timer := time.NewTimer(time.Microsecond * 900)
	go func() {
		// 发送http请求
	}()

	select {
	case <-ctx.Done():
		timer.Stop()
		timer.Reset(time.Second)
		fmt.Println("call successfully!!!")
		return
	case <-timer.C:
		fmt.Println("timeout!!!")
		return
	}
}
