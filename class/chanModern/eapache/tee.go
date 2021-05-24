package main

import (
	"fmt"
	"github.com/eapache/channels"
)

func main() {
	fmt.Println("tee:")
	a := channels.NewNativeChannel(channels.None)
	outputs := []channels.Channel{
		channels.NewNativeChannel(channels.None),
		channels.NewNativeChannel(channels.None),
		channels.NewNativeChannel(channels.None),
		channels.NewNativeChannel(channels.None),
	}
	channels.Tee(a, outputs[0], outputs[1], outputs[2], outputs[3])
	//channels.WeakTee(a, outputs[0], outputs[1], outputs[2], outputs[3])
	go func() {
		for i := 0; i < 5; i++ {
			a.In() <- i
		}
		a.Close()
	}()
	for i := 0; i < 20; i++ {
		var v interface{}
		var j int
		select {
		case v = <-outputs[0].Out():
			j = 0
		case v = <-outputs[1].Out():
			j = 1
		case v = <-outputs[2].Out():
			j = 2
		case v = <-outputs[3].Out():
			j = 3
		}
		fmt.Printf("channel#%d: %d\n", j, v)
	}
}
