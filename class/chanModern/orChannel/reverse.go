package orChannel

func or(chs ...<-chan interface{}) <-chan interface{} {
	switch len(chs) {
	case 0:
		return nil
	case 1:
		return chs[0]
	}
	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		switch len(chs) {
		case 2:
			select {
			case <-chs[0]:
			case <-chs[1]:
			}
		default:
			m := len(chs) / 2
			select {
			case <-or(chs[:m]...):
			case <-or(chs[m:]...):
			}
		}
	}()
	return orDone
}
