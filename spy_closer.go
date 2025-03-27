package testutils

import "bytes"

// ReadCloserSpy Позволят проверить что был вызват Close()
type ReadCloserSpy struct {
	IsClosed bool
	Data     *bytes.Buffer
}

// NewReadCloserSpy конструктор
func NewReadCloserSpy() *ReadCloserSpy {
	return &ReadCloserSpy{
		IsClosed: false,
		Data:     bytes.NewBufferString(""),
	}
}

// Read читать
func (r *ReadCloserSpy) Read(p []byte) (n int, err error) {
	return r.Data.Read(p)
}

// Close закрыть
func (r *ReadCloserSpy) Close() error {
	r.IsClosed = true
	return nil
}
