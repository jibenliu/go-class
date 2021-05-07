package dbPool

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"testing"
	"time"
)

func connTest()  {
	//容量、扫描时间、键值默认过期时间
	pool, _ := NewGenericPool(10, 10*time.Second, 5*time.Second)
	c := pool.Get()
	// 通过Do函数，发送redis命令
	v, err := c.Conn.Do("SET", "name1", "小王")
	if err != nil {
		fmt.Println(err)
		return
	}

	v, err = redis.String(c.Conn.Do("GET", "name1"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)
	_ = pool.Publish(c)
	time.Sleep(time.Second)
	c = pool.Get()
	//通过Do函数，发送redis命令
	v, err = c.Conn.Do("SET", "name2", "李四")
	if err != nil {
		fmt.Println(err)
		return
	}
	v, err = redis.String(c.Conn.Do("GET", "name2"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)
	_ = pool.Publish(c)

	time.Sleep(time.Second)
	c = pool.Get()
	//通过Do函数，发送redis命令
	v, err = c.Conn.Do("SET", "name3", "sb")
	if err != nil {
		fmt.Println(err)
		return
	}
	v, err = redis.String(c.Conn.Do("GET", "name3"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)
	_ = pool.Publish(c)
}

func TestJanitor_Run(t *testing.T) {
	connTest()
}

func Benchmark_Run(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		connTest()
	}
}


// go test -v -bench=.
// go test -v -bench=. -benchtime=5s
// go test -v -bench=Alloc -benchmem