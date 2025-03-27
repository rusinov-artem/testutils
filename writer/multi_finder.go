package writer

import (
	"fmt"
	"strings"
	"time"
)

type MultiFinder struct {
	Needle map[string]bool
	Found  bool
	Ch     chan struct{}
}

func NewMultiFinder(substrList ...string) *MultiFinder {
	f := &MultiFinder{
		Needle: make(map[string]bool),
		Ch:     make(chan struct{}),
	}
	for i := range substrList {
		f.Needle[substrList[i]] = false
	}

	return f
}

func (t *MultiFinder) Write(data []byte) (int, error) {
	fmt.Printf("%s: %s", "finder", string(data))

	everyLineFound := true
	for k, v := range t.Needle {
		if strings.Contains(string(data), k) {
			fmt.Printf("Found '%s' in:\n  %s\n", k, string(data))
			v = true
			t.Needle[k] = true
		}

		if !v {
			everyLineFound = false
		}
	}

	if !t.Found && everyLineFound {
		t.Found = true
		close(t.Ch)
	}

	return len(data), nil
}

func (t *MultiFinder) Wait(d time.Duration) error {
	reportBuilder := strings.Builder{}
	for k, v := range t.Needle {
		reportBuilder.WriteString(fmt.Sprintf("\t%s => %v\n", k, v))
	}

	select {
	case <-t.Ch:
		return nil
	case <-time.After(d):
		return fmt.Errorf("timeout. waiting for\n %s", reportBuilder.String())
	}
}
