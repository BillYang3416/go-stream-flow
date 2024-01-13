package rabbitmq

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/bgg/go-file-gate/internal/entity"
	"github.com/bgg/go-file-gate/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserUploadedFilePublisher struct {
	l          logger.Logger
	ch         *amqp.Channel
	exchange   string
	routingKey string
}

func NewUserUploadedFilePublisher(l logger.Logger, ch *amqp.Channel) *UserUploadedFilePublisher {
	pub := &UserUploadedFilePublisher{l: l, ch: ch, exchange: "user-uploaded-file", routingKey: "user-uploaded-file.event.created"}

	// declare exchange
	err := ch.ExchangeDeclare(
		pub.exchange,
		"direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		pub.l.Error(err, "failed to declare exchange")
	}

	queueName := "user-uploaded-file-created-queue"
	// declare queue
	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-deleted
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		pub.l.Error(err, "failed to declare queue")
	}

	// declare binding
	err = ch.QueueBind(
		queueName,
		pub.routingKey,
		pub.exchange,
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		pub.l.Error(err, "failed to declare binding")
	}

	return pub
}

func (pub *UserUploadedFilePublisher) Publish(ctx context.Context, file entity.UserUploadedFile) error {

	file.Base64Content = base64.StdEncoding.EncodeToString(file.Content)
	body, err := json.Marshal(file)
	if err != nil {
		pub.l.Error(err, "failed to marshal file")
		return err
	}

	err = pub.ch.PublishWithContext(
		ctx,
		pub.exchange,
		pub.routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		pub.l.Error(err, "failed to publish message")
		return err
	}

	return nil
}
