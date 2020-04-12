package util

// ErrChain return an error when any of the funcs in chain fail
func ErrChain(funcs ...func() error) (err error) {
	for _, fn := range funcs {
		err = fn()
		if err != nil {
			break
		}
	}
	return
}
