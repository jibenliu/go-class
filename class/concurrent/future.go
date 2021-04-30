package main

import "fmt"

type Function func(string) (string, error)

type Future interface {
	SuccessCallback() error
	FailCallback() error
	Execute(Function) (bool, chan struct{})
}

type AccountCache struct {
	Name string
}

func (a *AccountCache) SuccessCallback() error {
	fmt.Println("It's success")
	return nil
}

func (a *AccountCache) FailCallback() error {
	fmt.Println("It's fail")
	return nil
}

func (a *AccountCache) Execute(f Function) (bool, chan struct{}) {
	done := make(chan struct{})
	go func(a *AccountCache) {
		_, err := f(a.Name)
		if err != nil {
			_ = a.FailCallback()
		} else {
			_ = a.SuccessCallback()
		}
		done <- struct{}{} //为什么使用 struct 类型作为 channel 的通知 因为空 struct 在 Go 中占的内存是最少的
	}(a)
	return true, done
}

func NewAccountCache(name string) *AccountCache {
	return &AccountCache{
		name,
	}
}

func main() {
	var future Future
	future = NewAccountCache("Tom")
	updateFunc := func(name string) (string, error) {
		fmt.Println("cache update:", name)
		return name, nil
	}

	_, done := future.Execute(updateFunc)
	defer func() {
		<-done
	}()
}
