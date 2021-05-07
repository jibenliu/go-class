package main

import "fmt"

// Wheel 车轮
type Wheel struct {
	Size int
}

// Engine 引擎
type Engine struct {
	Power int    // 功率
	Type  string // 类型
}

// Car 车
type Car struct {
	Wheel
	Engine
}

func main() {

	c := Car{

		// 初始化轮子
		Wheel: Wheel{
			Size: 18,
		},

		// 初始化引擎
		Engine: Engine{
			Type:  "1.4T",
			Power: 143,
		},
	}

	fmt.Printf("%+v\n", c)

}
