package goroutinePool

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type sig struct{}

var (
	ErrSizeNotAllowed = errors.New("pool size <0,not create")
	ErrChanClosed     = errors.New("pool is close")
)

type Pool struct {
	// capacity of the pool.
	// capacity是该Pool的容量，也就是开启worker数量的上限，每一个worker需要一个goroutine去执行；
	// worker类型为任务类。
	capacity int32
	// running is the number of the currently running goroutines.
	// running是当前正在执行任务的worker数量
	running int32
	// expiryDuration set the expired time (second) of every worker.
	// expiryDuration是worker的过期时长，在空闲队列中的worker的最新一次运行时间与当前时间之差如果大于这个值则表示已过期，定时清理任务会清理掉这个worker；
	expiryDuration time.Duration
	// workers is a slice that store the available workers
	// 任务队列
	workers []*Worker
	// release is used to notice the pool to closed itself.
	// 当关闭该Pool支持通知所有worker退出运行以防goroutine泄露
	release chan sig
	// lock for synchronous operation
	// 用以支持Pool的同步操作
	lock sync.Mutex
	// once用在确保Pool关闭操作只会执行一次
	once sync.Once
}

func NewPool(size, expiry int) (*Pool, error) {
	if size <= 0 {
		return nil, ErrSizeNotAllowed
	}

	p := &Pool{
		capacity:       int32(size),
		release:        make(chan sig, 1),
		expiryDuration: time.Duration(expiry) * time.Second,
		running:        0,
	}
	// 启动定期清理过期worker任务，独立goroutine运行，
	// 进一步节省系统资源
	p.monitorAndClear()
	return p, nil
}

func (p *Pool) GetWorker() *Worker {
	var w *Worker
	// 标志，表示当前运行的worker数量是否已达容量上限
	waiting := false
	// 涉及从workers队列取可用worker，需要加锁
	p.lock.Lock()
	workers := p.workers
	n := len(workers) - 1
	fmt.Println("空闲worker数量:", n+1)
	fmt.Println("协程池现在运行的worker数量：", p.running)
	// 当前worker队列为空(无空闲worker)
	if n < 0 {
		//没有空闲的worker有两种可能：
		//1.运行的worker超出了pool容量
		//2.当前是空pool，从未往pool添加任务或者一段时间内没有任务添加，被定期清除
		// 运行worker数目已达到该Pool的容量上限，置等待标志
		if p.running >= p.capacity {
			print("超过上限")
			waiting = true
		} else {
			// 当前无空闲worker但是Pool还没有满，
			// 则可以直接新开一个worker执行任务
			p.running++
			w = &Worker{
				pool: p,
				task: make(chan funcType),
				str:  make(chan string),
			}
		}
	} else {
		// 有空闲worker，从队列尾部取出一个使用
		w = workers[n]
		workers[n] = nil
		p.workers = workers[:n]
		p.running++
	}
	// 判断是否有worker可用结束，解锁
	p.lock.Unlock()
	if waiting {
		for {
			// 自旋等待获取worker
			// 当一个任务执行完以后会添加到池中，有了空闲的任务就可以继续执行：
			// 阻塞等待直到有空闲worker
			for len(p.workers) == 0 {
				continue
			}
			p.lock.Lock()
			workers = p.workers
			l := len(workers) - 1
			w = workers[l]
			workers[l] = nil
			p.workers = workers[:1]
			p.running++
			p.lock.Unlock()
			break
		}
	}
	return w
}

// 定期清理过期worker
func (p *Pool) monitorAndClear() {
	go func() {
		for {
			// 周期性循环检查过期worker并清理
			time.Sleep(p.expiryDuration)
			currentTime := time.Now()
			p.lock.Lock()
			idleWorkers := p.workers
			n := 0
			for i, w := range idleWorkers {
				// 计算当前时间减去该worker的最后运行时间之差是否符合过期时长
				if currentTime.Sub(w.recycleTime) <= p.expiryDuration {
					break
				}
				n = i
				w.stop()
				idleWorkers[i] = nil
				p.running--
			}
			if n > 0 {
				n++
				p.workers = idleWorkers[n:]
			}
			p.lock.Unlock()
		}
	}()
}

// Worker回收（goroutine复用）
// putWorker puts a worker back into free pool, recycling the goroutines.
func (p *Pool) putWorker(worker *Worker) {
	// 写入回收时间，亦即该worker的最后运行时间
	worker.recycleTime = time.Now()
	p.lock.Lock()
	p.running--
	p.workers = append(p.workers, worker)
	p.lock.Unlock()
}

// Resize 动态扩容或者缩小池容量
// ReSize change the capacity of this pool
func (p *Pool) Resize(size int) {
	c := int(p.capacity)
	if size < c {
		diff := c - size
		for i := 0; i < diff; i++ {
			p.GetWorker().stop()
		}
	} else if size == c {
		return
	}
	atomic.StoreInt32(&p.capacity, int32(size))
}

func (p *Pool) Submit(task funcType, str string) error {
	if len(p.release) > 0 {
		return ErrChanClosed
	}
	//创建或得到一个空闲的worker
	w := p.GetWorker()
	w.run()
	//将任务参数通过信道传递给它
	w.sendArg(str)
	//将任务通过信道传递给它
	w.sendTask(task)
	return nil
}
