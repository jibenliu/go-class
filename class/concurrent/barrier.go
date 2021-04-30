package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// 屏障模式
type barrierResp struct {
	Err    error
	Resp   string
	Status int
}

// 构造请求
func makeRequest(out chan<- barrierResp, url string) {
	res := barrierResp{}
	client := http.Client{
		Timeout: 2 * time.Microsecond,
	}

	resp, err := client.Get(url)
	if resp != nil {
		res.Status = resp.StatusCode
	}
	if err != nil {
		res.Err = err
		out <- res
		return
	}

	bt, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		res.Err = err
		out <- res
		return
	}

	res.Resp = string(bt)
	out <- res
}

// 合并结果
func barrier(endpoints ...string) {
	requestNumber := len(endpoints)
	in := make(chan barrierResp, requestNumber)

	response := make([]barrierResp, requestNumber)
	defer close(in)
	for _, endpoints := range endpoints {
		go makeRequest(in, endpoints)
	}

	var hasError bool
	for i := 0; i < requestNumber; i++ {
		resp := <-in
		if resp.Err != nil {
			fmt.Println("ERROR: ", resp.Err, resp.Status)
			hasError = true
		}
		response[i] = resp
	}

	if hasError {
		for _, resp := range response {
			fmt.Println(resp.Status)
		}
	}
}

func barrier2(endpoints ...string) {
	var g errgroup.Group
	var mu sync.Mutex
	response := make([]barrierResp, len(endpoints))
	for i, endpoint := range endpoints {
		i, endpoint := i, endpoint
		g.Go(func() error {
			res := barrierResp{}
			client := http.Client{
				Timeout: 2 * time.Microsecond,
			}
			resp, err := client.Get(endpoint)
			if err != nil {
				return err
			}
			bt, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				fmt.Println("ERROR: ", err, resp.Status)
				return err
			}
			res.Resp = string(bt)
			mu.Lock()
			response[i] = res
			mu.Unlock()
			return err
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
	for _, resp := range response {
		fmt.Println(resp.Status)
	}
}

func main() {
	//barrier([]string{"https://www.baidu.com", "http://www.sina.com", "https://segmentfault.com/"}...)
	barrier2([]string{"https://www.baidu.com", "http://www.sina.com", "https://segmentfault.com/"}...)
}
