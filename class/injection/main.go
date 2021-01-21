package main

import (
	"fmt"
	"reflect"
)

//var inj *Injector

type Injector struct {
	mappers map[reflect.Type]reflect.Value
}

func (inj *Injector) SetMap(value interface{}) {
	inj.mappers[reflect.TypeOf(value)] = reflect.ValueOf(value)
}

func (inj *Injector) Get(t reflect.Type) reflect.Value {
	return inj.mappers[t]
}

func (inj *Injector) Invoke(i interface{}) interface{} {
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Func {
		panic("should invoke a function")
	}
	inValues := make([]reflect.Value, t.NumIn())

	for k := 0; k < t.NumIn(); k++ {
		inValues[k] = inj.Get(t.In(k))
	}
	ret := reflect.ValueOf(i).Call(inValues)
	return ret
}

func (inj *Injector) Host(name string, f func(a int, b string) string) {
	fmt.Println("Enter Host:", name)
	fmt.Println(inj.Invoke(f))
	fmt.Println("Exit Host:", name)
}

func Dependency(a int, b string) string {
	fmt.Println("Dependency: ", a, b)
	return `injection function exec finished ...`
}

func main() {
	in := &Injector{make(map[reflect.Type]reflect.Value)}
	in.SetMap(3030)
	in.SetMap("zdd")
	d := Dependency
	in.Host("zddhub", d)
	in.SetMap(8080)
	in.SetMap("www.zddhub.com")
	in.Host("website", d)
}
