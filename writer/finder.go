package writer

import (
	"fmt"
	"strings"
	"time"
)

type Finder struct {
	Needle string
	Found  bool
	Ch     chan struct{}
}

func NewFinder(substr string) *Finder {
	return &Finder{
		Needle: substr,
		Ch:     make(chan struct{}),
	}
}

func (t *Finder) Write(data []byte) (int, error) {
	fmt.Printf("%s: %s", t.Needle, string(data))
	if strings.Contains(string(data), t.Needle) {
		fmt.Printf("Found '%s' in:\n  %s\n", t.Needle, string(data))
		if !t.Found {
			t.Found = true
			close(t.Ch)
		}
	}
	return len(data), nil
}

func (t *Finder) Wait(d time.Duration) error {
	select {
	case <-t.Ch:
		return nil
	case <-time.After(d):
		return fmt.Errorf("timeout. waiting for %s", t.Needle)
	}
}
