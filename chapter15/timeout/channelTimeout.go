package main

import (
	"context"
	"fmt"
	"time"
)

func AsyncCall2() {
	ctx := context.Background()
	done := make(chan struct{}, 1)
	go func(ctx context.Context) {

		// 发送HTTP请求后推送信号
		//@TODO warning,如果http请求发生异常，则很可能channel推送不会执行
		done <- struct{}{}
	}(ctx)

	select {
	case <-done:
		fmt.Println("call successfully!!!")
		return
	case <-time.After(800 * time.Millisecond):
		fmt.Println("timeout!!!")
		return
	}
}
