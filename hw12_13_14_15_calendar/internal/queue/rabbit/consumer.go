package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Consumer struct {
	uri        string
	exchange   string
	queue      string
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewConsumer(uri, exchange, queue string) *Consumer {
	return &Consumer{
		uri:      uri,
		exchange: exchange,
		queue:    queue,
	}
}

func (c *Consumer) Connect() error {
	connection, err := amqp.Dial(c.uri)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("channel error: %w", err)
	}

	if err = channel.ExchangeDeclare(
		c.exchange, // name
		"direct",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return fmt.Errorf("exchange declare error: %w", err)
	}

	if _, err = channel.QueueDeclare(
		c.queue, // name of the queue
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // noWait
		nil,     // arguments
	); err != nil {
		return fmt.Errorf("queue declare error: %w", err)
	}

	if err = channel.QueueBind(
		c.queue,    // name of the queue
		"",         // bindingKey
		c.exchange, // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return fmt.Errorf("queue bind error: %w", err)
	}

	c.connection = connection
	c.channel = channel

	return nil
}

func (c *Consumer) Disconnect() {
	if c.connection != nil {
		c.connection.Close()
	}
}

func (c *Consumer) ConsumeEvents() (<-chan amqp.Delivery, error) {
	deliveries, err := c.channel.Consume(
		c.queue,
		"simple-consumer",
		false, // noAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue consume error: %w", err)
	}

	return deliveries, nil
}
