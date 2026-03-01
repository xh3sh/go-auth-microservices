package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config С…СЂР°РЅРёС‚ РєРѕРЅС„РёРіСѓСЂР°С†РёСЋ РїСЂРёР»РѕР¶РµРЅРёСЏ РёР· РїРµСЂРµРјРµРЅРЅС‹С… РѕРєСЂСѓР¶РµРЅРёСЏ
type Config struct {
	RabbitMQUser      string
	RabbitMQPassword  string
	RabbitMQHost      string
	RabbitMQPort      string
	RabbitMQQueueName string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	JWTSecret            string
	JWTExpiration        int64
	JWTRefreshExpiration int64

	AuthServicePort string
	AuthServiceHost string
	APIGatewayPort  string
	APIGatewayHost  string
	UserServicePort string
	UserServiceHost string
	FrontendPort    string
	FrontendHost    string
	LogConsumerPort string
	LogConsumerHost string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

// Load Р·Р°РіСЂСѓР¶Р°РµС‚ РєРѕРЅС„РёРіСѓСЂР°С†РёСЋ РёР· .env С„Р°Р№Р»Р° РёР»Рё СЃРёСЃС‚РµРјРЅС‹С… РїРµСЂРµРјРµРЅРЅС‹С…
func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		RabbitMQUser:      getEnv("RABBITMQ_USER", "guest"),
		RabbitMQPassword:  getEnv("RABBITMQ_PASSWORD", "guest"),
		RabbitMQHost:      getEnv("RABBITMQ_HOST", "localhost"),
		RabbitMQPort:      getEnv("RABBITMQ_PORT", "5672"),
		RabbitMQQueueName: getEnv("RABBITMQ_QUEUE_NAME", "auth_events_queue"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,

		JWTSecret:            getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiration:        int64(getEnvAsInt("JWT_EXPIRATION", 3600)),
		JWTRefreshExpiration: int64(getEnvAsInt("JWT_REFRESH_EXPIRATION", 86400)),

		AuthServicePort: getEnv("AUTH_SERVICE_PORT", "8081"),
		AuthServiceHost: getEnv("AUTH_SERVICE_HOST", "0.0.0.0"),
		APIGatewayPort:  getEnv("API_GATEWAY_PORT", "8080"),
		APIGatewayHost:  getEnv("API_GATEWAY_HOST", "0.0.0.0"),
		UserServicePort: getEnv("USER_SERVICE_PORT", "8082"),
		UserServiceHost: getEnv("USER_SERVICE_HOST", "0.0.0.0"),
		FrontendPort:    getEnv("FRONTEND_PORT", "8083"),
		FrontendHost:    getEnv("FRONTEND_HOST", "0.0.0.0"),
		LogConsumerPort: getEnv("LOG_CONSUMER_PORT", "8084"),
		LogConsumerHost: getEnv("LOG_CONSUMER_HOST", "0.0.0.0"),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8081/auth/oauth/google/callback"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
