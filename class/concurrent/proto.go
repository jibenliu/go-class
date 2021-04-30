package main

import "fmt"

type Cloneable interface {
	Clone() Cloneable
}

type PrototypeManager struct {
	prototypes map[string]Cloneable
}

func NewPrototypeManager() *PrototypeManager {
	return &PrototypeManager{
		prototypes: make(map[string]Cloneable),
	}
}

func (m *PrototypeManager) Get(name string) Cloneable {
	return m.prototypes[name]
}

func (m *PrototypeManager) Set(name string, prototype Cloneable) {
	m.prototypes[name] = prototype
}

// 测试
type Person struct {
	name   string
	age    int
	height int
}

func (p *Person) Clone() Cloneable {
	person := *p
	return &person
}

func main() {
	manager := NewPrototypeManager()

	person := &Person{
		name:   "zhangsan",
		age:    18,
		height: 175,
	}

	manager.Set("person", person)
	c := manager.Get("person").Clone()

	person1 := c.(*Person)

	fmt.Println("name:", person1.name)
	fmt.Println("age:", person1.age)
	fmt.Println("height:", person1.height)
}

// https://juejin.cn/post/6844903728533733389
