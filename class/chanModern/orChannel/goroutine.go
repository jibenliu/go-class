package orChannel

import "sync"

func or(chs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		var once sync.Once
		for _, c := range chs {
			go func(c <-chan interface{}) {
				select {
				case <-c:
					once.Do(func() { close(out) })
				case <-out:
				}
			}(c)
		}
	}()
	return out
}
