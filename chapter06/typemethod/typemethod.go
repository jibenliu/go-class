package main

import (
	"fmt"
)

// MyInt 将int定义为MyInt类型
type MyInt int

// IsZero 为MyInt定义IsZero方法
func (m MyInt) IsZero() bool {
	return m == 0
}

// Add 为MyInt添加Add方法
func (m MyInt) Add(other int) int {
	return other + int(m)
}

func main() {

	var b MyInt

	fmt.Println(b.IsZero())

	b = 1

	fmt.Println(b.Add(2))

}
