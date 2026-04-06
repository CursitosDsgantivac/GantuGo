package RabbitMq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// QueueConfig holds the configuration needed to declare an exchange, queue and binding.
type QueueConfig struct {
	ExchangeName string
	QueueName    string
	// ExchangeType is typically "fanout", "direct", or "topic"
	ExchangeType string
}

// SetupExchangeAndQueue declares a fanout exchange, a queue, and binds them together.
// It returns the declared queue so callers can reference its server-assigned name
// (important when QueueName is "" and RabbitMQ generates one).
//
// Call this once during startup on both the producer and consumer side so that
// the topology exists before any messages are published or consumed.
func SetupExchangeAndQueue(ch *amqp.Channel, cfg QueueConfig) (amqp.Queue, error) {
	// Declare the exchange.
	// durable=true  → survives a broker restart.
	// autoDelete=false → not deleted when the last queue unbinds.
	// internal=false  → external producers can publish to it.
	err := ch.ExchangeDeclare(
		cfg.ExchangeName, // name
		cfg.ExchangeType, // type  (e.g. "fanout")
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait (wait for server confirmation)
		nil,              // extra arguments
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	// Declare the queue.
	// durable=true  → queue survives a broker restart.
	// exclusive=false → other connections can also access this queue.
	// autoDelete=false → not deleted when the last consumer disconnects.
	//
	// Tip: pass QueueName="" to let RabbitMQ generate a unique name.
	// That is the typical fanout pattern where every consumer gets its own
	// temporary queue.
	q, err := ch.QueueDeclare(
		cfg.QueueName, // name ("" = auto-generated)
		true,          // durable
		false,         // auto-deleted
		false,         // exclusive
		false,         // no-wait
		nil,           // extra arguments
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	// Bind the queue to the exchange.
	// The routing key is ignored by fanout exchanges, so we pass "".
	err = ch.QueueBind(
		q.Name,           // queue name (use the returned name in case it was auto-generated)
		"",               // routing key (irrelevant for fanout)
		cfg.ExchangeName, // exchange
		false,            // no-wait
		nil,              // extra arguments
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	return q, nil
}
