package syncDemo

import (
	"fmt"
	"testing"
	"time"
)

func hello1() {
	fmt.Println("123")
}

func hello2(i string) {
	fmt.Println(i)
}

func Test_janitor_RunAsyncCheck(t *testing.T) {
	j := &janitor{
		interval: time.Second,
		overtime: 5 * time.Second,
	}
	j.RunAsyncCheck(hello1)
	j.RunAsyncCheck(hello2, "23")
}
