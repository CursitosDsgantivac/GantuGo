package RabbitMq

import (
	"context"
	"log"
	"strconv"
	"time"
)

/*
Library github.com/rabbitmq/amqp091-go
*/
const AMQP_URL = "amqp://guest:guest@localhost:5672/"

var queueConfig = QueueConfig{
	ExchangeName: "go-exchange", // exchange name
	ExchangeType: "fanout",      // broadcast to all bound queues
	QueueName:    "go-queue",    // named queue (use "" for auto-generated)
}

func StartRabbitMqConsumer() {
	// create queue and exchange

	cc, err := NewConsumerConnection(AMQP_URL)
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
	pc, err := NewPublisherConnection(AMQP_URL)
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
	// scanner := bufio.NewScanner(os.Stdin)
	for {
		/*
			// Publish message mannually from terminal
			fmt.Println("Write the message and press enter to publish")
			scanner.Scan() // reads the full line including spaces
			message := scanner.Text()
			err = PublishMessage(ctx, pc, cfg.ExchangeName, []byte(message))
			if err != nil {
				log.Fatal(err)
			}
		*/

		rabbitChannel := make(chan []byte, 1000)

		for range 1 {
			go func() {
				for j := 0; j < 1000000; j++ {
					rabbitChannel <- []byte("test " + strconv.Itoa(j))
				}
			}()
		}

		ratePerSecond := 1000

		ticker := time.NewTicker(time.Second / time.Duration(ratePerSecond))
		defer ticker.Stop()

		for msg := range rabbitChannel {
			<-ticker.C // block here until the next tick slot is available
			err = PublishMessage(ctx, pc, cfg.ExchangeName, msg)
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}
