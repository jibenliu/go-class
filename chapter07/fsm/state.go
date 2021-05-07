package main

import (
	"reflect"
)

// State 状态接口
type State interface {

	// Name 获取状态名字
	Name() string

	// EnableSameTransit 该状态是否允许同状态转移
	EnableSameTransit() bool

	// OnBegin 响应状态开始时
	OnBegin()

	// OnEnd 响应状态结束时
	OnEnd()

	// CanTransitTo 判断能否转移到某状态
	CanTransitTo(name string) bool
}

// StateName 从状态实例获取状态名
func StateName(s State) string {
	if s == nil {
		return "none"
	}

	// 使用反射获取状态的名称
	return reflect.TypeOf(s).Elem().Name()
}
