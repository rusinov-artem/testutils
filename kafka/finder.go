package kafka

import (
	"fmt"
	"time"
)

type Finder struct {
	Found  bool
	Ch     chan struct{}
	search func(m *Message) bool
}

func NewFinder(search func(m *Message) bool) *Finder {
	f := &Finder{}

	f.Ch = make(chan struct{})
	f.search = search

	return f
}

func (f *Finder) Handle(m *Message) {
	if f.Found {
		return
	}

	r := f.search(m)
	if r {
		f.Found = true
		close(f.Ch)
	}
}

func (f Finder) Wait(d time.Duration) error {
	select {
	case <-f.Ch:
		return nil
	case <-time.After(d):
		return fmt.Errorf("unable to find message")
	}
}
