package RabbitMq

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
)

/*
Library github.com/rabbitmq/amqp091-go
*/

var queueConfig = QueueConfig{
	ExchangeName: "go-exchange", // exchange name
	ExchangeType: "fanout",      // broadcast to all bound queues
	QueueName:    "go-queue",    // named queue (use "" for auto-generated)
}

func StartRabbitMqConsumer() {
	// create queue and exchange

	cc, err := NewConsumerConnection("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	// 2. Define the exchange + queue topology
	cfg := queueConfig

	// 3. Declare the exchange, queue, and binding in one call
	q, err := SetupExchangeAndQueue(cc.Channel, cfg)
	if err != nil {
		log.Fatal(err)
	}

	go ConsumeMessages(cc, q.Name, "", func(body []byte) error {
		log.Printf("Received: %s", body)
		return nil
	})

	// block forever
	select {}

}

func StartRabbitMqPublisher() {
	// create queue and exchange

	ctx := context.Background()

	// 1. Open a publisher connection (or consumer — same setup call)
	pc, err := NewPublisherConnection("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	// 2. Define the exchange + queue topology
	cfg := queueConfig

	// 3. Declare the exchange, queue, and binding in one call
	_, err = SetupExchangeAndQueue(pc.Channel, cfg)
	if err != nil {
		log.Fatal(err)
	}

	// 4a. Publish a message using the exchange
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Write the message and press enter to publish")
		scanner.Scan() // reads the full line including spaces
		message := scanner.Text()
		err = PublishMessage(ctx, pc, cfg.ExchangeName, []byte(message))
		if err != nil {
			log.Fatal(err)
		}
	}
}
