package client

import "github.com/nats-io/nats.go"

// The IMessageClient interface defines a set of methods that any messaging
// system client implementation should implement.
// It provides an abstract instance that can be used for different message clients
// such as RabbitMQ, Kafka, etc. By using this interface, we can switch between
// different messaging systems without changing the rest of the code that uses it.
type IMessageClient interface {
	Connect() error
	Disconnect() error
	OnDisconnect(func())
	Publish(Subject, []byte) error
	Subscribe(Subject, func(msg *nats.Msg)) error
	Unsubscribe(Subject)
}
