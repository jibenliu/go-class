package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type TaskHandler func(interface{})

type Task struct {
	Param   interface{}
	Handler TaskHandler
}

type WorkerPoolImpl interface {
	AddWorker()
	SendTask(task Task)
	Release()
}

type WorkPool struct {
	wg   sync.WaitGroup
	inCh chan Task
}

func (d *WorkPool) AddWorker() {
	d.wg.Add(1)
	go func() {
		for task := range d.inCh {
			task.Handler(task.Param)
		}
		d.wg.Done()
	}()
}

func (d *WorkPool) Release() {
	close(d.inCh)
	d.wg.Wait()
}

func (d *WorkPool) SendTask(t Task) {
	d.inCh <- t
}

func NewWorkerPool(buffer int) WorkerPoolImpl {
	return &WorkPool{
		inCh: make(chan Task, buffer),
	}
}

func main() {
	bufferSize := 100
	var workerPool = NewWorkerPool(bufferSize)
	workers := 4
	for i := 0; i < workers; i++ {
		workerPool.AddWorker()
	}
	var sum int32
	testFunc := func(i interface{}) {
		n := i.(int32)
		atomic.AddInt32(&sum, n)
	}
	var i, n int32
	n = 1000
	for ; i < n; i++ {
		task := Task{
			i,
			testFunc,
		}
		workerPool.SendTask(task)
	}
	workerPool.Release()
	fmt.Println(sum)
}
