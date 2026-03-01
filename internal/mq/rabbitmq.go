package mq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	URL       string
	QueueName string
}

// Client РїСЂРµРґРѕСЃС‚Р°РІР»СЏРµС‚ РѕР±РµСЂС‚РєСѓ РґР»СЏ СЂР°Р±РѕС‚С‹ СЃ RabbitMQ СЃРѕРµРґРёРЅРµРЅРёРµРј
type Client struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func NewClient(cfg Config) (*Client, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		cfg.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	log.Printf("Connected to RabbitMQ, Queue: %s", q.Name)

	return &Client{
		Conn:    conn,
		Channel: ch,
		Queue:   q,
	}, nil
}

func (c *Client) Close() {
	if c.Channel != nil {
		c.Channel.Close()
	}
	if c.Conn != nil {
		c.Conn.Close()
	}
}
