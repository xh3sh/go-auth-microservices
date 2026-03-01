package repository

import (
	"github.com/xh3sh/go-auth-microservices/internal/models"
)

type EventRepository interface {
	PublishAuthEvent(event models.AuthEvent) error
	PublishAPIGatewayEvent(event models.APIGatewayEvent) error
	PublishUserActionEvent(event models.UserActionEvent) error
	PublishNotificationEvent(event models.NotificationEvent) error
	PublishTokenValidationEvent(event models.TokenValidationEvent) error
	Close() error
}
