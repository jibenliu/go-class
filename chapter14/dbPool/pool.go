package dbPool

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"reflect"
	"sync"
	"time"
)

var (
	ErrCapacityNotAllowed = errors.New("pool capacity < 0,not create")
)

type newConn func() redis.Conn

type janitor struct {
	c        *ConnPool
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run() {
	//创建定时器
	ticker := time.NewTimer(j.Interval)
	print("开启定时器\n")
	for {
		select {
		case <-ticker.C: // 当定时器每次到达设置的时间时就会向管道发送消息，检测链接队列中链接是否过期
			print("开始扫描\n")
			j.c.DeleteExpired()
		case <-j.stop: //监视器退出信道
			ticker.Stop()
			close(j.stop)
			return
		}
	}
}

func (j *janitor) stopJanitor() {
	j.stop <- true
}
func runJanitor(c *ConnPool, ci time.Duration) *janitor {
	j := &janitor{
		c:        c,
		Interval: ci,
		stop:     make(chan bool),
	}
	go j.Run() //运行监控器
	return j
}

type ConnPoolInterface interface {
	Get() (*Conn, error) //获取资源
	Publish(*Conn) error //释放资源，返回池中
	ShutDown() error     //关闭池
}

type ConnPool struct {
	lock              sync.Mutex
	ConnList          []*Conn       //链接
	capacity          int32         // 链接池最大链接限制
	numOpen           int32         // 当前池中空闲链接数
	running           int32         // 正在使用的链接数
	expiryDuration    time.Duration //扫描时间
	defaultExpiration time.Duration //链接的过期时间
	factory           newConn       // 创建连接的方法
	j                 *janitor      //监视器
	isClose           bool          //链接池是否关闭
}

func (c *ConnPool) monitorAndClear() *janitor {
	return runJanitor(c, c.expiryDuration)
}

///////////////////////////////

func (c *ConnPool) DeleteExpired() {
	// 现在时间戳
	now := time.Now().UnixNano()
	// 加互斥锁
	c.lock.Lock()
	for i, conn := range c.ConnList {
		// "Inlining" of expired
		if conn.Expiration > 0 && now > conn.Expiration {
			// 超时则删除
			o, ok := c.delete(c.ConnList, i)
			// 类型断言
			if ok {
				c.ConnList = o.([]*Conn)
			}
			conn.close()
		}
	}
	c.lock.Unlock()
}

func (c *ConnPool) delete(slice interface{}, index int) (interface{}, bool) {
	// 判断是否是切片类型
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, false
	}
	// 参数检查
	if v.Len() == 0 || index < 0 || index > v.Len()-1 {
		return nil, false
	}
	return reflect.AppendSlice(v.Slice(0, index), v.Slice(index+1, v.Len())).Interface(), true
}

// Get 获取资源
func (c *ConnPool) Get() *Conn {
	if c.isClose {
		return nil
	}
	var conn *Conn
	// 标志，表示当前运行的链接数量是否已达容量上限
	waiting := false
	// 涉及从链接队列取可用链接，需要加锁
	c.lock.Lock()
	ConnList := c.ConnList
	n := len(ConnList) - 1
	fmt.Println("空闲链接数量:", n+1)
	fmt.Println("链接池现在运行的链接数量：", c.running)
	// 当前链接队列为空(无空闲链接)
	if n < 0 {
		// 没有空闲的链接有两种可能：
		// 1.运行的链接超出了pool容量
		// 2.当前是空pool，从未往pool添加链接或者一段时间内没有链接添加，被定期清除
		// 运行链接数目已达到该Pool的容量上限，置等待标志
		if c.running >= c.capacity {
			print("超过上限")
			waiting = true
		} else {
			// 当前无空闲链接但是Pool还没有满，
			// 则可以直接新开一个链接执行任务
			c.running++
			conn = &Conn{
				Expiration: time.Now().Add(c.expiryDuration).UnixNano(),
				Conn:       c.factory(),
			}
		}
	} else { // 有空闲链接，从队列尾部取出一个使用
		conn = ConnList[n]
		ConnList[n] = nil
		c.ConnList = ConnList[:n]
		c.running++
	}
	// 解锁
	c.lock.Unlock()
	if waiting {
		//当一个链接执行完以后会添加到池中，有了空闲的链接就可以继续执行：
		// 阻塞等待直到有空闲链接
		for {
			for len(c.ConnList) == 0 {
				continue
			}
			c.lock.Lock()
			ConnList = c.ConnList
			l := len(ConnList) - 1
			conn = ConnList[l]
			ConnList[l] = nil
			c.ConnList = ConnList[:l]
			c.running++
			c.lock.Unlock()
			break
		}
	}
	return conn
}

// Publish 释放资源
func (c *ConnPool) Publish(conn *Conn) error {
	if c.isClose {
		return nil
	}
	conn.Expiration = time.Now().UnixNano()
	c.lock.Lock()
	c.running--
	c.ConnList = append(c.ConnList, conn)
	c.lock.Unlock()
	return nil
}
func (c *ConnPool) ShutDown() error {
	c.isClose = true
	for _, conn := range c.ConnList {
		conn.close()
	}
	c.j.stopJanitor()
	return nil
}

type Conn struct {
	Expiration int64
	Conn       redis.Conn //链接
}

func (c *Conn) close() {
	err := c.Conn.Close()
	if err != nil {
		log.Println("关闭链接失败!", err)
	}
}

func NewGenericPool(capacity int32, expiryDuration time.Duration, defaultExpiration time.Duration) (*ConnPool, error) {
	if capacity <= 0 {
		return nil, ErrCapacityNotAllowed
	}
	p := &ConnPool{
		capacity:          capacity,
		expiryDuration:    expiryDuration,
		running:           0,
		numOpen:           0,
		defaultExpiration: defaultExpiration,
		factory: func() redis.Conn {
			rs, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil
			}
			return rs
		},
		isClose: false,
	}
	// 启动定期清理过期链接，独立goroutine运行，
	// 进一步节省系统资源
	j := p.monitorAndClear()
	p.j = j
	return p, nil
}
