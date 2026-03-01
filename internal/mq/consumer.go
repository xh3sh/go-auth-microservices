package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/xh3sh/go-auth-microservices/internal/constants"
	"github.com/xh3sh/go-auth-microservices/internal/models"
	"github.com/xh3sh/go-auth-microservices/internal/repository"
)

// LogConsumer РїРѕС‚СЂРµР±Р»СЏРµС‚ СЃРѕР±С‹С‚РёСЏ РёР· РѕС‡РµСЂРµРґРµР№ Рё СЃРѕС…СЂР°РЅСЏРµС‚ РёС… РєР°Рє Р»РѕРіРё РІ СЂРµРїРѕР·РёС‚РѕСЂРёР№
type LogConsumer struct {
	repo      repository.Repository
	eventRepo *repository.RabbitMQRepository
}

func NewLogConsumer(repo repository.Repository, eventRepo *repository.RabbitMQRepository) *LogConsumer {
	return &LogConsumer{
		repo:      repo,
		eventRepo: eventRepo,
	}
}

// Start Р·Р°РїСѓСЃРєР°РµС‚ РїСЂРѕС†РµСЃСЃ РїРѕС‚СЂРµР±Р»РµРЅРёСЏ РёР· РІСЃРµС… РЅР°СЃС‚СЂРѕРµРЅРЅС‹С… РѕС‡РµСЂРµРґРµР№
func (c *LogConsumer) Start(ctx context.Context) {
	queues := []string{constants.QueueAuthEvents, constants.QueueAPIEvents, constants.QueueUserEvents}

	for _, q := range queues {
		go func(queueName string) {
			msgs, err := c.eventRepo.Consume(queueName)
			if err != nil {
				log.Printf("Failed to start consuming from %s: %v", queueName, err)
				return
			}

			log.Printf("Started consuming from %s", queueName)

			for {
				select {
				case <-ctx.Done():
					return
				case d, ok := <-msgs:
					if !ok {
						log.Printf("Channel closed for queue %s", queueName)
						return
					}
					c.processMessage(queueName, d.RoutingKey, d.Body)
				}
			}
		}(q)
	}
}

func (c *LogConsumer) processMessage(queueName, routingKey string, body []byte) {
	var entry models.LogEntry
	entry.Data = string(body)

	if len(routingKey) > 0 {
		switch {
		case routingKey == constants.RoutingKeyTokenValidation:
			var tokenEvent models.TokenValidationEvent
			if err := json.Unmarshal(body, &tokenEvent); err == nil {
				entry.UserID = tokenEvent.UserID
				entry.Service = "auth"
				entry.Type = "token_validation"
				entry.Timestamp = tokenEvent.Timestamp
			}
		case len(routingKey) >= len(constants.RoutingKeyAuthPrefix) && routingKey[:len(constants.RoutingKeyAuthPrefix)] == constants.RoutingKeyAuthPrefix:
			var authEvent models.AuthEvent
			if err := json.Unmarshal(body, &authEvent); err == nil {
				entry.UserID = authEvent.UserID
				entry.Service = "auth"
				entry.Type = authEvent.EventType
				entry.Timestamp = authEvent.Timestamp
			}
		case routingKey == constants.RoutingKeyAPIRequest:
			var apiEvent models.APIGatewayEvent
			if err := json.Unmarshal(body, &apiEvent); err == nil {
				entry.UserID = apiEvent.UserID
				entry.Service = "gateway"
				entry.Type = apiEvent.Action
				entry.Timestamp = apiEvent.Timestamp
			}
		case len(routingKey) > 5 && routingKey[:5] == "user.":
			var userEvent models.UserActionEvent
			if err := json.Unmarshal(body, &userEvent); err == nil {
				entry.UserID = userEvent.UserID
				entry.Service = "user"
				entry.Type = userEvent.Action
				entry.Timestamp = userEvent.Timestamp
			}
		}
	}

	if entry.Type == "" {
		switch queueName {
		case constants.QueueAuthEvents:
			var authEvent models.AuthEvent
			if err := json.Unmarshal(body, &authEvent); err == nil && authEvent.EventType != "" {
				entry.UserID = authEvent.UserID
				entry.Service = "auth"
				entry.Type = authEvent.EventType
				entry.Timestamp = authEvent.Timestamp
			}
		case constants.QueueAPIEvents:
			var apiEvent models.APIGatewayEvent
			if err := json.Unmarshal(body, &apiEvent); err == nil {
				entry.UserID = apiEvent.UserID
				entry.Service = "gateway"
				entry.Type = apiEvent.Action
				entry.Timestamp = apiEvent.Timestamp
			}
		case constants.QueueUserEvents:
			var userEvent models.UserActionEvent
			if err := json.Unmarshal(body, &userEvent); err == nil {
				entry.UserID = userEvent.UserID
				entry.Service = "user"
				entry.Type = userEvent.Action
				entry.Timestamp = userEvent.Timestamp
			}
		}
	}

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	
	entry.ID = fmt.Sprintf("%d-%s-%s", entry.Timestamp.UnixNano(), entry.Service, entry.UserID)

	if err := c.repo.SaveLog(context.Background(), entry); err != nil {
		log.Printf("Failed to save log: %v", err)
	}
}
