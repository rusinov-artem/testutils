package writer

import (
	"io"
	"sync"
)

type ProxyWriter struct {
	w io.Writer
	sync.Mutex
}

func NewProxy() *ProxyWriter {
	return &ProxyWriter{
		w: &PrefixWriter{Prefix: "Empty Proxy"},
	}
}

func (t *ProxyWriter) Write(data []byte) (int, error) {
	t.Lock()
	defer t.Unlock()
	return t.w.Write(data)
}

func (t *ProxyWriter) SetWriter(w io.Writer) {
	t.Lock()
	defer t.Unlock()
	t.w = w
}
