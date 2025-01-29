package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/thiagohmm/allLog/internal/rabbitmq"
	"github.com/thiagohmm/allLog/internal/usecase"
)

type MessageService struct {
	UseCase *usecase.LogUseCase // Changed to UseCase
}

func (s *MessageService) ListenToQueue(ctx context.Context, rabbitmqURL string, queueName string) { // Added context
	conn, err := rabbitmq.GetRabbitMQConnection(rabbitmqURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack  <- Important: changed to false for proper error handling
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %s", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var dto usecase.DTOIN               // Use the correct DTO type
			err := json.Unmarshal(d.Body, &dto) // Unmarshal the message body
			if err != nil {
				log.Printf("failed to unmarshal message: %s", err)
				d.Nack(false, true) // Requeue the message if unmarshaling fails.
				continue
			}

			err = s.UseCase.SaveLog(ctx, dto) // Call SaveLog on the UseCase
			if err != nil {
				log.Printf("failed to process message: %s", err)
				d.Nack(false, true) // Requeue the message on error.
			} else {
				d.Ack(false) // Acknowledge successful processing.
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
