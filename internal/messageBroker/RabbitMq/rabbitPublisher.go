package RabbitMq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublisherConnection holds the AMQP connection and an open channel ready for publishing.
// Close both when you are done (e.g. via defer).
type PublisherConnection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewPublisherConnection opens a connection and a channel to the RabbitMQ broker.
// amqpURL format: "amqp://user:password@host:port/"
// e.g.           "amqp://guest:guest@localhost:5672/"
//
// Returns a PublisherConnection that must be closed by the caller when no longer needed.
func NewPublisherConnection(amqpURL string) (*PublisherConnection, error) {
	// Establish the TCP connection to the broker.
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	// Open a channel — a lightweight, multiplexed logical connection.
	// Most AMQP operations (declare, publish, consume) happen on channels.
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &PublisherConnection{Conn: conn, Channel: ch}, nil
}

// Close shuts down the channel and the underlying connection gracefully.
func (p *PublisherConnection) Close() {
	if p.Channel != nil {
		p.Channel.Close()
	}
	if p.Conn != nil {
		p.Conn.Close()
	}
}

// PublishMessage sends a message to a fanout exchange.
// The exchange must have been declared beforehand (see SetupExchangeAndQueue).
//
// Parameters:
//   - ctx          → context for cancellation / deadline propagation
//   - pc           → active PublisherConnection returned by NewPublisherConnection
//   - exchangeName → name of the fanout exchange to publish to
//   - body         → raw message payload
func PublishMessage(ctx context.Context, pc *PublisherConnection, exchangeName string, body []byte) error {
	// PublishWithContext is the context-aware publish method.
	// For a fanout exchange the routing key ("") is ignored — the broker
	// delivers the message to every queue bound to the exchange.
	return pc.Channel.PublishWithContext(
		ctx,
		exchangeName, // exchange
		"",           // routing key (ignored by fanout)
		false,        // mandatory: if true, broker returns the msg when no queue is bound
		false,        // immediate: if true, broker returns the msg when no consumer is ready
		amqp.Publishing{
			ContentType: "application/json",
			// Persistent → message survives a broker restart.
			// Transient → message does not survive a broker restart.
			// Requires the target queue to also be durable.
			DeliveryMode: amqp.Transient,
			Body:         body,
			Headers: amqp.Table{
				"test": "test",
			},
		},
	)
}
