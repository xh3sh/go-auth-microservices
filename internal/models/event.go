package models

import "time"

// AuthEvent РїСЂРµРґСЃС‚Р°РІР»СЏРµС‚ СЃРѕР±С‹С‚РёРµ Р°СѓС‚РµРЅС‚РёС„РёРєР°С†РёРё
type AuthEvent struct {
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
}

// APIGatewayEvent РїСЂРµРґСЃС‚Р°РІР»СЏРµС‚ Р°РІС‚РѕСЂРёР·РѕРІР°РЅРЅРѕРµ СЃРѕР±С‹С‚РёРµ СЃ Gateway
type APIGatewayEvent struct {
	UserID     string    `json:"user_id"`
	AuthMethod string    `json:"auth_method"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	RequestID  string    `json:"request_id"`
	Timestamp  time.Time `json:"timestamp"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
}

// UserActionEvent РїСЂРµРґСЃС‚Р°РІР»СЏРµС‚ РґРµР№СЃС‚РІРёРµ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ РґР»СЏ Р»РѕРіРѕРІ
type UserActionEvent struct {
	UserID     string                 `json:"user_id"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID string                 `json:"resource_id,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	RequestID  string                 `json:"request_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Status     string                 `json:"status"`
}

// NotificationEvent РїСЂРµРґСЃС‚Р°РІР»СЏРµС‚ СЃРѕР±С‹С‚РёРµ СѓРІРµРґРѕРјР»РµРЅРёСЏ
type NotificationEvent struct {
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
	IsUrgent  bool      `json:"is_urgent"`
}

// TokenValidationEvent РїСЂРµРґСЃС‚Р°РІР»СЏРµС‚ СЃРѕР±С‹С‚РёРµ РїСЂРѕРІРµСЂРєРё С‚РѕРєРµРЅР°
type TokenValidationEvent struct {
	UserID       string    `json:"user_id"`
	AuthMethod   string    `json:"auth_method"`
	TokenID      string    `json:"token_id,omitempty"`
	IsValid      bool      `json:"is_valid"`
	ErrorMessage string    `json:"error_message,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// LogEntry РїСЂРµРґСЃС‚Р°РІР»СЏРµС‚ СѓРЅРёС„РёС†РёСЂРѕРІР°РЅРЅСѓСЋ Р·Р°РїРёСЃСЊ Р»РѕРіР° РІ Redis
type LogEntry struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Service   string    `json:"service"`
	Type      string    `json:"type"`
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
