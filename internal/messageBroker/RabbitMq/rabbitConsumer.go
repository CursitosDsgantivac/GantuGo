package RabbitMq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConsumerConnection holds the AMQP connection and an open channel ready for consuming.
// Close both when you are done (e.g. via defer).
type ConsumerConnection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewConsumerConnection opens a connection and a channel to the RabbitMQ broker.
// amqpURL format: "amqp://user:password@host:port/"
// e.g.           "amqp://guest:guest@localhost:5672/"
//
// Returns a ConsumerConnection that must be closed by the caller when no longer needed.
func NewConsumerConnection(amqpURL string) (*ConsumerConnection, error) {
	// Establish the TCP connection to the broker.
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	// Open a channel for consuming.
	// Best practice: use a dedicated channel per consumer.
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &ConsumerConnection{Conn: conn, Channel: ch}, nil
}

// Close shuts down the channel and the underlying connection gracefully.
func (c *ConsumerConnection) Close() {
	if c.Channel != nil {
		c.Channel.Close()
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
}

// MessageHandler is the callback signature your application must implement.
// It receives one delivery at a time.
// Return an error to Nack (reject) the message; return nil to Ack it.
type MessageHandler func(body []byte) error

// ConsumeMessages starts consuming messages from queueName using the provided
// ConsumerConnection. It blocks until the deliveries channel is closed (e.g.
// when the connection drops or the context is cancelled by the caller).
//
// Parameters:
//   - cc          → active ConsumerConnection returned by NewConsumerConnection
//   - queueName   → name of the queue to consume from (must already be declared & bound)
//   - consumerTag → unique identifier for this consumer; "" lets the broker generate one
//   - handler     → callback invoked for every incoming message
//
// Usage: run this in a goroutine so it doesn't block your main flow.
//
//	go ConsumeMessages(cc, q.Name, "", func(body []byte) error {
//	    fmt.Println(string(body))
//	    return nil
//	})
func ConsumeMessages(cc *ConsumerConnection, queueName, consumerTag string, handler MessageHandler) error {
	// Register the consumer.
	// autoAck=false → we manually Ack/Nack after processing.
	//   This is the safer pattern: if the consumer crashes before Ack-ing,
	//   the broker re-queues the message for another consumer.
	msgs, err := cc.Channel.Consume(
		queueName,   // queue
		consumerTag, // consumer tag ("" = auto-generated)
		false,       // auto-ack (false = manual ack)
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // extra arguments
	)
	if err != nil {
		return err
	}

	// Range over the deliveries channel; this blocks until the channel is closed.
	for d := range msgs {
		if d.Headers["test"] != nil {
			fmt.Println("Header received with message: ", d.Headers["test"])
		}
		if err := handler(d.Body); err != nil {
			// Nack without requeue=true re-delivers to another consumer.
			// Change the second bool to false to dead-letter instead of requeue.
			d.Nack(
				false, // multiple: false = only this message
				true,  // requeue: true = put it back in the queue
			)
			continue
		}

		// Ack tells the broker the message was processed successfully and can be removed.
		// multiple=false → only acknowledge this specific delivery tag.
		d.Ack(false)
	}

	return nil
}
