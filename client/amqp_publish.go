package client

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
)

func runAmqpPublish(cmd *cobra.Command, args []string) {
	url := "amqp://guest:guest@localhost:5672/"

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

	if err := ch.Confirm(false); err != nil {
		log.Fatalf("failed to enable publisher confirms: %v", err)
	}

	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}

	ctx := cmd.Context()

	df, err := ch.PublishWithDeferredConfirmWithContext(ctx, "", q.Name, false, false, amqp091.Publishing{Body: []byte("Hello")})
	if err != nil {
		log.Fatalf("failed to publish a message: %v", err)
	}

	if ack := df.Wait(); ack {
		log.Printf("message is confirmed to be delivered with delivery tag: %d", df.DeliveryTag)
	} else {
		log.Printf("message is not confirmed to be delivered")
	}

	log.Println("published a message")
}
