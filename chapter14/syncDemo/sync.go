package syncDemo

import (
	"fmt"
	"reflect"
	"time"
)

type janitor struct {
	interval time.Duration
	overtime time.Duration
}

func (j *janitor) RunAsyncCheck(f interface{}, params ...interface{}) {
	// 创建周期断续器
	ticker := time.NewTicker(j.interval)
	// 创建定时器
	timer := time.NewTimer(j.overtime)
loop:
	for {
		select {
		case <-timer.C: //当Timer每次到达设置的时间时就会向管道发送消息，此时超时退出
			print("超时\n")
			break loop
		case <-ticker.C: //当Ticker每次到达设置的时间时就会向管道发送消息，此时进行异步check操作
			print("异步check\n")
			fv := reflect.ValueOf(f)
			realParams := make([]reflect.Value, len(params))
			for i, item := range params {
				realParams[i] = reflect.ValueOf(item)
			}
			fv.Call(realParams)
		}
	}
	fmt.Println("Break")
}
