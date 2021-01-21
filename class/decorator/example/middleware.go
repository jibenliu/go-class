package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

//1.实现一个httpServer
//2.实现一个handler
//3.实现中间件的功能， 1）记录请求的URL和方法 2）记录请求的网络的地址 3）记录方法的执行时间
func hello(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	log.Printf("记录请求的网络地址：%s", r.RemoteAddr)
	log.Printf("记录方法的执行时间：%s", time.Since(startTime))
	fmt.Fprintf(w, "hello")
}

func main() {
	http.Handle("/", tracing(http.HandlerFunc(hello)))
	http.ListenAndServe(":8080", nil)
}

func tracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("记录请求的URL和方法：%s %s", r.URL, r.Method)
		next.ServeHTTP(w, r)
	})
}
