package repository

import (
	"context"
	"github.com/xh3sh/go-auth-microservices/internal/models"
)

type LogRepository interface {
	SaveLog(ctx context.Context, log models.LogEntry) error
	GetLogs(ctx context.Context, page, pageSize int) ([]models.LogEntry, int64, error)
	FilterLogs(ctx context.Context, userID string, service, logType string, page, pageSize int) ([]models.LogEntry, int64, error)
}
