package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

const (
	// DebugMode indicates service mode is debug.
	DebugMode = "debug"
	// TestMode indicates service mode is test.
	TestMode = "test"
	// ReleaseMode indicates service mode is release.
	ReleaseMode = "release"

	TimeExpiredAt = time.Hour * 800
)

type Config struct {
	ServiceName string
	Environment string // debug, test, release
	Version     string

	ServiceHost string
	HTTPPort    string
	HTTPScheme  string
	Domain      string

	DefaultOffset  string
	DefaultLimit   string
	DefaultBarCode string
	DefaultSaleId  string

	UserServiceHost string
	UserGRPCPort    string

	ProductServiceHost string
	ProductGRPCPort    string

	OrderServiceHost string
	OrderGRPCPort    string

	PostgresMaxConnections int32

	SecretKey string
}

const (
	TokenExpireTime = 24 * time.Hour
	JWTSecretKey    = "MySecretKey"
)

// Load ...
func Load() Config {
	if err := godotenv.Load("./.env"); err != nil {
		fmt.Println("No .env file found")
	}

	config := Config{}

	config.ServiceName = cast.ToString(getOrReturnDefaultValue("SERVICE_NAME", "api_gateway_service"))
	config.Environment = cast.ToString(getOrReturnDefaultValue("ENVIRONMENT", DebugMode))
	config.Version = cast.ToString(getOrReturnDefaultValue("VERSION", "1.0"))

	config.ServiceHost = cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "localhost"))
	config.HTTPPort = cast.ToString(getOrReturnDefaultValue("HTTP_PORT", ":8001"))
	config.HTTPScheme = cast.ToString(getOrReturnDefaultValue("HTTP_Scheme", "http"))
	config.Domain = cast.ToString(getOrReturnDefaultValue("DOMAIN", "localhost:8001"))

	config.DefaultOffset = cast.ToString(getOrReturnDefaultValue("OFFSET", "0"))
	config.DefaultLimit = cast.ToString(getOrReturnDefaultValue("LIMIT", "10"))
	config.DefaultBarCode = cast.ToString(getOrReturnDefaultValue("BAR_CODE", ""))
	config.DefaultSaleId = cast.ToString(getOrReturnDefaultValue("SALE_ID", ""))

	config.ProductServiceHost = cast.ToString(getOrReturnDefaultValue("PRODUCT_SERVICE_HOST", "localhost"))
	config.ProductGRPCPort = cast.ToString(getOrReturnDefaultValue("PRODUCT_GRPC_PORT", ":50052"))

	config.OrderServiceHost = cast.ToString(getOrReturnDefaultValue("ORDER_SERVICE_HOST", "localhost"))
	config.OrderGRPCPort = cast.ToString(getOrReturnDefaultValue("ORDER_GRPC_PORT", ":50053"))

	config.UserServiceHost = cast.ToString(getOrReturnDefaultValue("USER_SERVICE_HOST", "localhost"))
	config.UserGRPCPort = cast.ToString(getOrReturnDefaultValue("USER_GRPC_PORT", ":50053"))

	config.PostgresMaxConnections = cast.ToInt32(getOrReturnDefaultValue("POSTGRES_MAX_CONNECTIONS", 30))

	config.SecretKey = cast.ToString(getOrReturnDefaultValue("SECRET_KEY", "final"))

	return config
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	val, exists := os.LookupEnv(key)

	if exists {
		return val
	}

	return defaultValue
}
