package kafka

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

type Header struct {
	Key   string
	Value []byte
}

type Message struct {
	Topic     string
	Partition int
	Offset    int64
	Key       []byte
	Value     []byte
	Headers   []Header
	Time      time.Time
}

type HandlerFunc func(m *Message)

func (h *HandlerFunc) Handle(m *Message) {
	(*h)(m)
}

type Handler interface {
	Handle(m *Message)
}

type Consumer struct {
	t       *testing.T
	lock    *sync.Mutex
	handler Handler
}

func NewConsumer(ctx context.Context, t *testing.T, broker []string, topic, group string) *Consumer {
	c := &Consumer{}
	c.t = t
	c.lock = &sync.Mutex{}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     broker,
		Topic:       topic,
		GroupID:     group,
		MaxBytes:    10e6,
		MaxWait:     100 * time.Millisecond,
		StartOffset: kafka.LastOffset,
	})

	r.SetOffset(kafka.LastOffset)

	ctx, stopReading := context.WithCancel(context.Background())

	c.t.Cleanup(func() {
		stopReading()
		r.Close()
	})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				m, err := r.ReadMessage(ctx)
				err = c.consume(m, err)
				if err != nil {
					t.Logf("consumer stopped: %s", err)
					return
				}
			}
		}
	}()

	return c
}

func (c *Consumer) consume(m kafka.Message, err error) error {
	if err != nil {
		return err
	}

	c.lock.Lock()
	if c.handler != nil {
		c.handler.Handle(convert(m))
	}
	c.lock.Unlock()
	return nil
}

func (c *Consumer) SetHandler(handler Handler) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.handler = handler
}

func convert(m kafka.Message) *Message {
	msg := &Message{
		Topic:     m.Topic,
		Partition: m.Partition,
		Offset:    m.Offset,
		Key:       m.Key,
		Value:     m.Value,
		Time:      m.Time,
	}

	msg.Headers = make([]Header, len(m.Headers))
	for i := range m.Headers {
		msg.Headers[i] = Header{
			Key:   m.Headers[i].Key,
			Value: m.Headers[i].Value,
		}
	}

	return msg
}
