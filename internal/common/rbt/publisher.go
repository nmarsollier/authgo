package rbt

import (
	"encoding/json"

	"github.com/nmarsollier/authgo/internal/common/log"
	"github.com/streadway/amqp"
)

type RabbitPublisher[T any] interface {
	PublishForResult(data T, exchange string, routingKey string) error
	PublishTo(exchange string, routingKey string, data T) error
	Publish(data T) error
	Logger() log.LogRusEntry
}

type rabbitPublisher[T any] struct {
	ch           *amqp.Channel
	log          log.LogRusEntry
	exchangeName string
	channelType  string
	routingKey   string
}

func NewRabbitPublisher[T any](
	log log.LogRusEntry,
	rabbitURL string,
	exchangeName string,
	channelType string,
	routingKey string,
) (RabbitPublisher[T], error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &rabbitPublisher[T]{
		ch:           channel,
		exchangeName: exchangeName,
		channelType:  channelType,
		routingKey:   routingKey,
		log:          log,
	}, nil
}

func (c *rabbitPublisher[T]) exchangeDeclare(exchangeName string, channelType string) error {
	return c.ch.ExchangeDeclare(
		exchangeName, // name
		channelType,  // type
		false,        // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

func (c *rabbitPublisher[T]) Publish(data T) error {
	return c.PublishTo(c.exchangeName, c.routingKey, data)
}

func (c *rabbitPublisher[T]) PublishForResult(data T, exchange string, routingKey string) error {
	return c.publish(c.exchangeName, c.routingKey, data, exchange, routingKey)
}

func (c *rabbitPublisher[T]) PublishTo(exchange string, routingKey string, data T) error {
	return c.publish(exchange, routingKey, data, "", "")
}

func (c *rabbitPublisher[T]) publish(exchange string, routingKey string, data T, fbExchange string, fbRoutingKey string) error {
	correlationId := c.log.CorrelationId()
	c.log.WithField(log.LOG_FIELD_RABBIT_ACTION, "Emit").
		WithField(log.LOG_FIELD_RABBIT_EXCHANGE, exchange).
		WithField(log.LOG_FIELD_RABBIT_QUEUE, routingKey)

	err := c.exchangeDeclare(exchange, c.channelType)
	if err != nil {
		c.log.Error(err)
		return err
	}

	send := publishMessage[T]{
		CorrelationId: correlationId,
		Exchange:      fbExchange,
		RoutingKey:    fbRoutingKey,
		Message:       data,
	}

	body, err := json.Marshal(send)
	if err != nil {
		c.log.Error(err)
		return err
	}

	err = c.ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Body: body,
		})

	if err != nil {
		c.log.Error(err)
		return err
	}

	c.log.Info("Rabbit publish ", exchange+" ", routingKey+" ", string(body))

	return nil
}

func (c *rabbitPublisher[T]) Logger() log.LogRusEntry {
	return c.log
}

type publishMessage[T any] struct {
	CorrelationId string `json:"correlation_id" example:"123123" `
	Exchange      string `json:"exchange" example:"cart"`
	RoutingKey    string `json:"routing_key" example:""`
	Message       T      `json:"message" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbklEIjoiNjZiNjBlYzhlMGYzYzY4OTUzMzJlOWNmIiwidXNlcklEIjoiNjZhZmQ3ZWU4YTBhYjRjZjQ0YTQ3NDcyIn0.who7upBctOpmlVmTvOgH1qFKOHKXmuQCkEjMV3qeySg" `
}

func RbtLogger(fluentUrl string, serverName string, correlationId string) log.LogRusEntry {
	return log.Get(fluentUrl, serverName).
		WithField(log.LOG_FIELD_CORRELATION_ID, correlationId).
		WithField(log.LOG_FIELD_CONTROLLER, "Rabbit")
}
