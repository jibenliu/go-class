package main

import (
	"fmt"
	"github.com/eapache/channels"
)

func main() {
	fmt.Println("multi:")
	a := channels.NewNativeChannel(channels.None)
	inputs := []channels.Channel{
		channels.NewNativeChannel(channels.None),
		channels.NewNativeChannel(channels.None),
		channels.NewNativeChannel(channels.None),
		channels.NewNativeChannel(channels.None),
	}
	channels.Multiplex(a, inputs[0], inputs[1], inputs[2], inputs[3])
	//channels.WeakMultiplex(a, inputs[0], inputs[1], inputs[2], inputs[3])
	go func() {
		for i := 0; i < 5; i++ {
			for j := range inputs {
				inputs[j].In() <- i
			}
		}
		for i := range inputs {
			inputs[i].Close()
		}
	}()
	for v := range a.Out() {
		fmt.Printf("%d ", v)
	}
}
