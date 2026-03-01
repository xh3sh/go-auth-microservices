package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xh3sh/go-auth-microservices/internal/constants"
	"github.com/xh3sh/go-auth-microservices/internal/models"
)

// RabbitMQRepository РїСЂРµРґРѕСЃС‚Р°РІР»СЏРµС‚ С„СѓРЅРєС†РёРѕРЅР°Р» РґР»СЏ СЂР°Р±РѕС‚С‹ СЃ РѕС‡РµСЂРµРґСЏРјРё СЃРѕРѕР±С‰РµРЅРёР№
type RabbitMQRepository struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	mu      sync.Mutex
}

func NewRabbitMQRepository(user, password, host, port string) (*RabbitMQRepository, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	
	var conn *amqp.Connection
	var err error
	for i := 0; i < constants.RabbitMQRetryAttempts; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ, retrying in %d seconds... (%d/%d)", constants.RabbitMQRetryInterval, i+1, constants.RabbitMQRetryAttempts)
		time.Sleep(constants.RabbitMQRetryInterval * time.Second)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	err = ch.ExchangeDeclare(constants.ExchangeName, constants.ExchangeType, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	repo := &RabbitMQRepository{
		conn:    conn,
		channel: ch,
	}

	if err := repo.SetupTopology(); err != nil {
		log.Printf("SetupTopology error: %v", err)
	}

	return repo, nil
}

// SetupTopology РЅР°СЃС‚СЂР°РёРІР°РµС‚ РѕС‡РµСЂРµРґРё Рё РїСЂРёРІСЏР·РєРё Рє РѕР±РјРµРЅРЅРёРєСѓ
func (r *RabbitMQRepository) SetupTopology() error {
	queues := map[string][]string{
		constants.QueueAuthEvents: {constants.PatternAuthEvents, constants.PatternTokenEvents},
		constants.QueueAPIEvents:  {constants.PatternAPIEvents},
		constants.QueueUserEvents: {constants.PatternUserEvents},
	}

	for qName, keys := range queues {
		_, err := r.channel.QueueDeclare(
			qName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", qName, err)
		}

		for _, key := range keys {
			err = r.channel.QueueBind(qName, key, constants.ExchangeName, false, nil)
			if err != nil {
				return fmt.Errorf("failed to bind queue %s with key %s: %w", qName, key, err)
			}
		}
	}
	return nil
}

func (r *RabbitMQRepository) publish(routingKey string, body interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel == nil || r.conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection is closed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.RabbitMQPublishTimeout*time.Second)
	defer cancel()

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	err = r.channel.PublishWithContext(ctx,
		constants.ExchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: constants.ContentTypeJSON,
			Body:        jsonBody,
		})
	
	if err != nil {
		log.Printf("RabbitMQ Publish Error (key: %s): %v", routingKey, err)
	}
	return err
}

func (r *RabbitMQRepository) PublishAuthEvent(event models.AuthEvent) error {
	return r.publish(fmt.Sprintf("%s%s", constants.RoutingKeyAuthPrefix, event.EventType), event)
}

func (r *RabbitMQRepository) PublishAPIGatewayEvent(event models.APIGatewayEvent) error {
	return r.publish(constants.RoutingKeyAPIRequest, event)
}

func (r *RabbitMQRepository) PublishUserActionEvent(event models.UserActionEvent) error {
	return r.publish(fmt.Sprintf("user.%s", event.Action), event)
}

func (r *RabbitMQRepository) PublishNotificationEvent(event models.NotificationEvent) error {
	return r.publish(constants.RoutingKeyUserNotification, event)
}

func (r *RabbitMQRepository) PublishTokenValidationEvent(event models.TokenValidationEvent) error {
	return r.publish(constants.RoutingKeyTokenValidation, event)
}

// Consume РїРѕРґРїРёСЃС‹РІР°РµС‚СЃСЏ РЅР° РѕС‡РµСЂРµРґСЊ Рё РІРѕР·РІСЂР°С‰Р°РµС‚ РєР°РЅР°Р» СЃРѕРѕР±С‰РµРЅРёР№
func (r *RabbitMQRepository) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQRepository) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
