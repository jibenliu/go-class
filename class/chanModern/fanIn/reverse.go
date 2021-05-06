package fanIn

func fanInRec(chs ...<-chan interface{}) <-chan interface{} {
	switch len(chs) {
	case 0:
		c := make(chan interface{})
		close(c)
		return c
	case 1:
		return chs[0]
	case 2:
		return mergeTwo(chs[0], chs[1])
	default:
		m := len(chs) / 2
		return mergeTwo(
			fanInRec(chs[:m]...),
			fanInRec(chs[m:]...))
	}
}
func mergeTwo(a, b <-chan interface{}) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				c <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				c <- v
			}
		}
	}()
	return c
}
