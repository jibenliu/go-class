package set

func reduce(in <-chan interface{}, fn func(r, v interface{}) interface{}) interface{} {
	if in == nil {
		return nil
	}
	out := <-in
	for v := range in {
		out = fn(out, v)
	}
	return out
}
