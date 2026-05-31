package models

import "time"

// AuthEvent представляет событие аутентификации
type AuthEvent struct {
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	TraceID   string    `json:"trace_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// APIGatewayEvent представляет авторизованное событие с Gateway
type APIGatewayEvent struct {
	UserID     string    `json:"user_id"`
	AuthMethod string    `json:"auth_method"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	RequestID  string    `json:"request_id"`
	TraceID    string    `json:"trace_id"`
	Timestamp  time.Time `json:"timestamp"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
}

// UserActionEvent представляет действие пользователя для логов
type UserActionEvent struct {
	UserID     string                 `json:"user_id"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID string                 `json:"resource_id,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	RequestID  string                 `json:"request_id"`
	TraceID    string                 `json:"trace_id,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Status     string                 `json:"status"`
}

// NotificationEvent представляет событие уведомления
type NotificationEvent struct {
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
	IsUrgent  bool      `json:"is_urgent"`
}

// TokenValidationEvent представляет событие проверки токена
type TokenValidationEvent struct {
	UserID       string    `json:"user_id"`
	AuthMethod   string    `json:"auth_method"`
	TokenID      string    `json:"token_id,omitempty"`
	IsValid      bool      `json:"is_valid"`
	TraceID      string    `json:"trace_id,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// LogEntry представляет унифицированную запись лога в Redis
type LogEntry struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Service   string    `json:"service"`
	Type      string    `json:"type"`
	TraceID   string    `json:"trace_id,omitempty"`
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
