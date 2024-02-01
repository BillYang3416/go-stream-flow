package rabbitmq

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/bgg/go-flow-gateway/internal/entity"
	"github.com/bgg/go-flow-gateway/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserUploadedFilePublisher struct {
	logger     logger.Logger
	ch         *amqp.Channel
	exchange   string
	routingKey string
}

func NewUserUploadedFilePublisher(l logger.Logger, ch *amqp.Channel) *UserUploadedFilePublisher {
	pub := &UserUploadedFilePublisher{logger: l, ch: ch, exchange: "user-uploaded-file", routingKey: "user-uploaded-file.event.created"}

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
		pub.logger.Error("UserUploadedFilePublisher - NewUserUploadedFilePublisher - ch.ExchangeDeclare: failed to declare exchange", "error", err)
	}
	pub.logger.Info("UserUploadedFilePublisher - NewUserUploadedFilePublisher: successfully declared exchange", "exchange", pub.exchange)

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
		pub.logger.Error("UserUploadedFilePublisher - NewUserUploadedFilePublisher - ch.QueueDeclare: failed to declare queue", "error", err)
	}
	pub.logger.Info("UserUploadedFilePublisher - NewUserUploadedFilePublisher: successfully declared queue", "queue", queueName)

	// declare binding
	err = ch.QueueBind(
		queueName,
		pub.routingKey,
		pub.exchange,
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		pub.logger.Error("UserUploadedFilePublisher - NewUserUploadedFilePublisher - ch.QueueBind: failed to bind queue", "error", err)
	}
	pub.logger.Info("UserUploadedFilePublisher - NewUserUploadedFilePublisher: successfully bound queue", "queue", queueName, "exchange", pub.exchange, "routingKey", pub.routingKey)

	return pub
}

func (pub *UserUploadedFilePublisher) Publish(ctx context.Context, file entity.UserUploadedFile) error {

	file.Base64Content = base64.StdEncoding.EncodeToString(file.Content)
	body, err := json.Marshal(file)
	if err != nil {
		pub.logger.Error("UserUploadedFilePublisher - Publish - json.Marshal: failed to marshal file", "error", err)
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
		pub.logger.Error("UserUploadedFilePublisher - Publish - ch.Publish: failed to publish message", "error", err)
		return err
	}

	pub.logger.Info("UserUploadedFilePublisher - Publish: successfully published message", "exchange", pub.exchange, "routingKey", pub.routingKey)
	return nil
}
