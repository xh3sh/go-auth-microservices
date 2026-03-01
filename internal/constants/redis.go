package constants

import "time"

const (
	RedisPrefix = "api-microservice:"
	DefaultTTL  = 24 * time.Hour
	APIKeyTTL   = 24 * time.Hour
)

const (
	PrefixAPIKey       = "api_key:"
	PrefixSession      = "session:"
	PrefixUser         = "user:"
	PrefixUsername     = "username:"
	PrefixBlacklist    = "blacklist:"
	PrefixRefreshToken = "refresh_token:"
	PrefixLogs         = "logs:all"
	PrefixLogUser      = "logs:user:"
	PrefixLogService   = "logs:service:"
	PrefixLogType      = "logs:type:"
	LogTTL             = 24 * time.Hour
)

const (
	ValueRevoked = "revoked"
)

func BuildRedisKey(prefix, identifier string) string {
	return prefix + identifier
}
