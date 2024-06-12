package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/storage"
)

type Producer struct {
	uri        string
	exchange   string
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewProducer(uri, exchange string) *Producer {
	return &Producer{
		uri:      uri,
		exchange: exchange,
	}
}

func (p *Producer) Connect() error {
	connection, err := amqp.Dial(p.uri)
	if err != nil {
		return fmt.Errorf("dial error: %w", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("channel error: %w", err)
	}

	if err := channel.ExchangeDeclare(
		p.exchange, // name
		"direct",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return fmt.Errorf("exchange declare error: %w", err)
	}

	p.connection = connection
	p.channel = channel

	return nil
}

func (p *Producer) Disconnect() {
	if p.connection != nil {
		p.connection.Close()
		p.connection = nil
		p.channel = nil
	}
}

func (p *Producer) ProduceEvent(event storage.Event) error {
	eventJSON, err := event.MarshalJSON()
	if err != nil {
		return fmt.Errorf("error marshalling event: %w", err)
	}

	err = p.channel.Publish(
		p.exchange,
		"",
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            eventJSON,
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
	)
	if err != nil {
		return fmt.Errorf("error while publishing to rabbit mq: %w", err)
	}

	return nil
}
