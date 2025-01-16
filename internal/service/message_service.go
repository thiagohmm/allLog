package service

import (
    "log"
    "project/internal/entity"
    "project/internal/repository"
)

type MessageService struct {
    Repo *repository.MessageRepository
}




   func (s *MessageService) ListenToQueue(rabbitmqURL string) {
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
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true, // auto-ack changed to true for at-least-once delivery
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %s", err)
	}

	// Use a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case d, ok := <-msgs:
				if !ok {
					return // Channel closed, exit
				}
				log.Printf("Received a message: %s", d.Body)
				message := entity.Message{
					Table:  "messages",
					Fields: map[string]interface{}{
						"message": string(d.Body),
					},
				}
				err := s.Repo.SaveMessage(message) // Ensure this function handles db errors gracefully
				if err != nil {
					// Implement retry logic, dead-letter queue, or other error handling mechanisms.
					log.Printf("failed to save message: %s", err)
                    // Example retry logic with exponential backoff
					time.Sleep(time.Second * 2) // Initial retry after 2 seconds
					for i := 0; i < 3 && err != nil ; i++ { // Retry 3 times
                        err = s.Repo.SaveMessage(message)
						if err != nil {
							time.Sleep(time.Second * time.Duration(2 << i)) // Exponential backoff
						}
					}
                    if err != nil {
                        // Log the error and potentially publish the message to a dead-letter queue
                        log.Printf("Failed to save message after multiple retries: %s", err)
                        // ... publish to dead-letter queue ...
                    }
				}


			case <-ctx.Done():
				return // Context cancelled, exit
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-ctx.Done() // Block until context is cancelled (e.g., by a signal handler)

}

