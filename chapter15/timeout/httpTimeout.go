package main

import (
	"context"
	"fmt"
	"time"
)

func AsyncCall() {
	// ctx是从上游一直传递过来的，对于上游传递过来的context还剩多少时间无法确定，所以需要设置一个自己预期的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*800)
	defer cancel() //为了实现没到定时时间任务执行完了回收context
	go func() {
		// 发送HTTP请求
	}()
	select {
	case <-ctx.Done():
		fmt.Println("call successfully!!!")
		return
	case <-time.After(time.Duration(time.Millisecond * 900)):
		fmt.Println("timeout!!!")
		return
	}
}
