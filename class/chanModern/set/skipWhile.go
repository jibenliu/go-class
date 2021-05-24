package set

func skipWhile(done <-chan struct{}, valueStream <-chan interface{}, fn func(interface{}) bool) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		take := false
		for {
			select {
			case <-done:
				return
			case v := <-valueStream:
				if !take {
					take = !fn(v)
					if !take {
						continue
					}
				}
				takeStream <- v
			}
		}
	}()
	return takeStream
}
