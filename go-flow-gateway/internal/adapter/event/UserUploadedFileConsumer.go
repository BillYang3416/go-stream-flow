package event

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/bgg/go-file-gate/internal/entity"
	"github.com/bgg/go-file-gate/internal/usecase"
	"github.com/bgg/go-file-gate/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserUploadedFileConsumer struct {
	u  usecase.UserUploadedFile
	l  logger.Logger
	ch *amqp.Channel
	q  amqp.Queue
}

func NewUserUploadedFileConsumer(u usecase.UserUploadedFile, ch *amqp.Channel, l logger.Logger) *UserUploadedFileConsumer {
	cs := &UserUploadedFileConsumer{u: u, ch: ch, l: l}

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

		var userUploadedFile entity.UserUploadedFile
		err := json.Unmarshal(d.Body, &userUploadedFile)
		if err != nil {
			cs.l.Error(err, "failed to unmarshal message")
		}
		cs.l.Info("Received a message: %s", userUploadedFile.Name)

		decodedContent, err := base64.StdEncoding.DecodeString(userUploadedFile.Base64Content)
		if err != nil {
			cs.l.Error(err, "failed to decode base64 content")
		}
		userUploadedFile.Content = decodedContent

		err = cs.u.SendEmail(context.Background(), userUploadedFile)
		if err != nil {
			cs.l.Error(err, "failed to send email")
			continue
		}
	}

}
