package client

import (
	"errors"
	"sync"

	"github.com/nats-io/nats.go"
)

var ErrNoSubscription = errors.New("Subscription does not exist.")

type NatsClient struct {
	url  string
	conn *nats.Conn
	opts []nats.Option
	// We store subscriptions in a SyncMap to ensure thread-safety,
	// as it may be accessed/changed concurrently by multiple goroutines.
	subscriptions sync.Map
}

func NewNatsClient(url string, opts []nats.Option) *NatsClient {
	return &NatsClient{url: url, opts: opts}
}

func (c *NatsClient) Connect() (err error) {
	nc, err := nats.Connect(c.url, c.opts...)
	if err != nil {
		return
	}
	c.conn = nc
	return
}

func (c *NatsClient) Disconnect() (err error) {
	if c.conn != nil {
		c.conn.Close()
	}
	return
}

func (c *NatsClient) OnDisconnect(cb func()) {
	c.conn.SetDisconnectHandler(func(_ *nats.Conn) {
		cb()
	})
}

func (c *NatsClient) Publish(subject Subject, data []byte) (err error) {
	err = c.conn.Publish(string(subject), data)
	return
}

func (c *NatsClient) Subscribe(subject Subject, handler func(msg *nats.Msg)) (err error) {
	// Above Subscribe method of `NatsClient` runs the provided handler function, which returns a consumer function.
	// Prior to processing messages, the handler may perform some business logic and initialization steps.
	// The returned consumer function is responsible for consuming messages received from the subscribed subject.
	sub, err := c.conn.Subscribe(string(subject), handler)
	if err != nil {
		return
	}

	// Add the subscription subject/topic names to the c.subscriptions SyncMap to keep track of them.
	// This allows us to unsubscribe from any location without needing to pass the subscription instance
	// and also keeps the interface looking more general.
	// We store it in a SyncMap to ensure thread-safety,
	// as it may be accessed/changed concurrently by multiple goroutines.
	c.subscriptions.Store(subject, sub)
	return
}

func (c *NatsClient) Unsubscribe(subject Subject) {
	item, loaded := c.subscriptions.LoadAndDelete(subject)
	if !loaded {
		return
	}
	sub := item.(*nats.Subscription)
	_ = sub.Unsubscribe()
}
