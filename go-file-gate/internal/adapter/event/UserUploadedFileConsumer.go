package event

import (
	"github.com/bgg/go-file-gate/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserUploadedFileConsumer struct {
	l  logger.Logger
	ch *amqp.Channel
	q  amqp.Queue
}

func NewUserUploadedFileConsumer(ch *amqp.Channel, l logger.Logger) *UserUploadedFileConsumer {
	cs := &UserUploadedFileConsumer{ch: ch, l: l}

	q, err := ch.QueueDeclare(
		"user-uploaded-file-created-queue",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		cs.l.Error(err, "failed to declare queue")
	}
	cs.q = q
	return cs
}

func (cs *UserUploadedFileConsumer) StartConsume() {
	cs.l.Info(" UserUploadedFileConsumer - StartConsume - Waiting for messages.")

	msgs, err := cs.ch.Consume(
		cs.q.Name, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		cs.l.Error(err, "failed to register a consumer")
	}

	for d := range msgs {
		cs.l.Info("Received a message: %s", d.Body)
	}

}
