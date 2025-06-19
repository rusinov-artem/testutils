package kafka

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

type Checker func(m kafka.Message) bool

// тут лучше перейти на 
// https://github.com/twmb/franz-go
// так как тут есть возможность удалять consumer groups
type Finder struct {
	t *testing.T

	found bool

	topic   string
	broker  string
	checker Checker

	done chan struct{}

	ctx         context.Context
	stopReading context.CancelFunc
}

func NewFinder(t *testing.T, broker, topic string) *Finder {
	return &Finder{
		t:      t,
		topic:  topic,
		broker: broker,
		done:   make(chan struct{}),
	}
}

func (f *Finder) Find(fn Checker) {
	f.checker = fn

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    f.topic,
		GroupID:  f.testName(),
		MaxBytes: 10e6,
		MaxWait:  100 * time.Millisecond,
	})

	f.t.Cleanup(func() {
		r.Close()
	})

	f.ctx, f.stopReading = context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-f.ctx.Done():
				return
			default:
				m, err := r.ReadMessage(f.ctx)
				if err != nil {
					break
				}
				if f.checker(m) {
					f.found = true
					close(f.done)
				}
			}
		}
	}()
}

func (f *Finder) Wait(d time.Duration) {
	f.t.Helper()

	select {
	case <-time.After(d):
		break
	case <-f.done:
		break
	}

	f.stopReading()

	msg := fmt.Sprintf("expected kafka message in topic: %s not found", f.topic)
	assert.True(f.t, f.found, msg)
}

func (f *Finder) testName() string {
	name := f.t.Name()
	return strings.ReplaceAll(name, "_", "/")
}
