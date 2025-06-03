package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnv(filenames ...string) error {
	err := godotenv.Load(filenames...)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error loading .env file(s) %v: %w", filenames, err)
		}
	}

	requiredVars := []string{
		"MONGO_URI",
		"MONGO_DATABASE",
		"KAFKA_BROKERS_NOTIFICATION",
		"KAFKA_TOPIC_NOTIFICATION",
		"AUTH_SERVICE_ADDRESS",
		"KAFKA_BROKERS_BOARD_EVENTS",
		"KAFKA_TOPIC_BOARD_EVENTS",
		"KAFKA_GROUP_ID_BOARD_EVENTS",
	}

	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return fmt.Errorf("required environment variable %s is not set", v)
		}
	}
	return nil
}

func GetMongoURI() string {
	return os.Getenv("MONGO_URI")
}

func GetMongoDatabase() string {
	return os.Getenv("MONGO_DATABASE")
}

func GetMongoTimeout() time.Duration {
	timeoutStr := GetEnvDefault("MONGO_TIMEOUT", "10s")
	duration, err := time.ParseDuration(timeoutStr)
	if err != nil {
		defaultDuration, _ := time.ParseDuration("10s")
		return defaultDuration
	}
	return duration
}

func GetKafkaBrokersNotification() []string {
	brokers := os.Getenv("KAFKA_BROKERS_NOTIFICATION")
	if brokers == "" {
		return []string{}
	}
	return strings.Split(brokers, ",")
}

func GetKafkaTopicNotification() string {
	return os.Getenv("KAFKA_TOPIC_NOTIFICATION")
}

func GetAuthServiceAddress() string {
	return os.Getenv("AUTH_SERVICE_ADDRESS")
}

func GetKafkaBrokersBoardEvents() []string {
	brokers := os.Getenv("KAFKA_BROKERS_BOARD_EVENTS")
	if brokers == "" {
		return []string{}
	}
	return strings.Split(brokers, ",")
}

func GetKafkaTopicBoardEvents() string {
	return os.Getenv("KAFKA_TOPIC_BOARD_EVENTS")
}

func GetKafkaGroupIDBoardEvents() string {
	return os.Getenv("KAFKA_GROUP_ID_BOARD_EVENTS")
}

func GetAppPort() string {
	return GetEnvDefault("APP_PORT", "9090")
}

func GetAppName() string {
	return GetEnvDefault("APP_NAME", "calendar_service")
}

func GetAppVersion() string {
	return GetEnvDefault("APP_VERSION", "1.0.0")
}

func GetAppReadTimeout() time.Duration {
	return parseDurationWithDefault(GetEnvDefault("APP_READ_TIMEOUT", "5s"), 5*time.Second)
}

func GetAppWriteTimeout() time.Duration {
	return parseDurationWithDefault(GetEnvDefault("APP_WRITE_TIMEOUT", "10s"), 10*time.Second)
}

func GetAppIdleTimeout() time.Duration {
	return parseDurationWithDefault(GetEnvDefault("APP_IDLE_TIMEOUT", "120s"), 120*time.Second)
}

func GetEnvDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseDurationWithDefault(valueStr string, defaultValue time.Duration) time.Duration {
	duration, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return duration
}
