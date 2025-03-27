package writer

import "fmt"

type PrefixWriter struct {
	Prefix string
}

func (t *PrefixWriter) Write(data []byte) (int, error) {
	fmt.Printf("%s: %s", t.Prefix, string(data))
	return len(data), nil
}
