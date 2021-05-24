package main

import (
	"fmt"
	"github.com/eapache/channels"
)

func main() {
	fmt.Println("pipe:")
	a := channels.NewNativeChannel(channels.None)
	b := channels.NewNativeChannel(channels.None)
	channels.Pipe(a, b)
	// channels.WeakPipe(a, b)
	go func() {
		for i := 0; i < 5; i++ {
			a.In() <- i
		}
		a.Close()
	}()
	for v := range b.Out() {
		fmt.Printf("%d ", v)
	}
}
