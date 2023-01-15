package client

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
)

func runAmqpConsume(cmd *cobra.Command, args []string) {
	url := "amqp://guest:guest@localhost:5672/"
	consumer := "amqp-consume"

	conn, err := amqp091.Dial(url)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}

	dCh, err := ch.Consume(q.Name, consumer, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to register a consumer: %v", err)
	}

	d, ok := <-dCh
	if !ok {
		log.Fatalf("failed to receive a message")
	}

	log.Printf("received a message with %d bytes [%v] %q", len(d.Body), d.DeliveryTag, d.Body)

	if err := d.Ack(false); err != nil {
		log.Fatalf("failed to ack a message: %v", err)
	}

	if err := ch.Cancel(consumer, false); err != nil {
		log.Fatalf("failed to cancel a consumer: %v", err)
	}
}
