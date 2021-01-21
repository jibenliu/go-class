package main

import (
	"sync"
	"sync/atomic"
)

type singleton struct {
}

var instance = &singleton{}

func NewSingleton() *singleton {
	if instance == nil {
		instance = &singleton{}
	}
	return instance
}

func NewSingleton() *singleton {
	l.Lock() // lock
	defer l.Unlock()
	if instance == nil { // check
		instance = &singleton{}
	}
	return instance
}

var l sync.Locker

func NewSingleton() *singleton {
	if instance == nil { // check
		l.Lock() // lock
		defer l.Unlock()
		if instance == nil { // check
			instance = &singleton{}
		}
	}
	return instance
}

var mu sync.Mutex
var initialized uint32
func NewSingleton() *singleton {
	if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}
	mu.Lock()
	defer mu.Unlock()
	if initialized == 0 {
		instance = &singleton{}
		atomic.StoreUint32(&initialized, 1)
	}
	return instance
}

var once sync.Once

func NewSingleton() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}
