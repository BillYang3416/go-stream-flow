package event

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/internal/usecase"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserUploadedFileConsumer struct {
	userUploadedFile usecase.UserUploadedFile
	logger           logger.Logger
	ch               *amqp.Channel
	queue            amqp.Queue
}

func NewUserUploadedFileConsumer(u usecase.UserUploadedFile, ch *amqp.Channel, l logger.Logger) *UserUploadedFileConsumer {
	cs := &UserUploadedFileConsumer{userUploadedFile: u, ch: ch, logger: l}

	q, err := ch.QueueDeclare(
		"user-uploaded-file-created-queue",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		cs.logger.Error(err, "failed to declare queue")
	}
	cs.queue = q
	return cs
}

func (cs *UserUploadedFileConsumer) StartConsume() {
	cs.logger.Info("UserUploadedFileConsumer - StartConsume: start consuming messages")

	msgs, err := cs.ch.Consume(
		cs.queue.Name, // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		cs.logger.Error("UserUploadedFileConsumer - StartConsume - ch.Consume: failed to register a consumer", "error", err)
	}
	cs.logger.Info("UserUploadedFileConsumer - StartConsume: successfully registered a consumer")

	for d := range msgs {
		cs.processMessage(d)
	}

}

func (cs *UserUploadedFileConsumer) processMessage(d amqp.Delivery) {
	var userUploadedFile entity.UserUploadedFile
	err := json.Unmarshal(d.Body, &userUploadedFile)
	if err != nil {
		cs.logger.Error("UserUploadedFileConsumer - processMessage - json.Unmarshal: failed to unmarshal message body", "error", err)
	}
	cs.logger.Info("UserUploadedFileConsumer - processMessage: successfully unmarshalled message body", "userUploadedFileID", userUploadedFile.ID)

	decodedContent, err := base64.StdEncoding.DecodeString(userUploadedFile.Base64Content)
	if err != nil {
		cs.logger.Error("UserUploadedFileConsumer - processMessage - base64.StdEncoding.DecodeString: failed to decode base64 content", "error", err)
	}
	userUploadedFile.Content = decodedContent
	cs.logger.Info("UserUploadedFileConsumer - processMessage: successfully decoded base64 content", "userUploadedFileID", userUploadedFile.ID)

	err = cs.userUploadedFile.SendEmail(context.Background(), userUploadedFile)
	if err != nil {
		cs.logger.Error("UserUploadedFileConsumer - processMessage - userUploadedFile.SendEmail: failed to send email", "error", err)
	}
	cs.logger.Info("UserUploadedFileConsumer - processMessage: successfully sent email", "userUploadedFileID", userUploadedFile.ID)
}
