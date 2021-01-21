package main

type Object struct {
	item interface{}
}

func (o *Object) Get(item string) interface{} {
	return o.item
}

func (o *Object) Set(item string) bool {
	o.item = item
	return true
}
