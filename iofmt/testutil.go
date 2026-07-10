package iofmt

type ErrWriter struct {
	err error
}

func NewErrWriter(err error) *ErrWriter {
	w := &ErrWriter{
		err: err,
	}
	return w
}

func (e *ErrWriter) Write(p []byte) (int, error) {
	return -1, e.err
}
