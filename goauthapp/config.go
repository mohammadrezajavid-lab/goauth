package goauthapp

import (
	"github.com/mohammadrezajavid-lab/goauth/goauthapp/delivery/http/middleware"
	"github.com/mohammadrezajavid-lab/goauth/pkg/database"
	"github.com/mohammadrezajavid-lab/goauth/pkg/httpserver"
	"github.com/mohammadrezajavid-lab/goauth/pkg/logger"
	"github.com/mohammadrezajavid-lab/goauth/pkg/token"
	"time"
)

type Config struct {
	HTTPServer           httpserver.Config `koanf:"http_server"`
	RateLimiter          middleware.Config `koanf:"rate_limiter"`
	JWT                  token.Config      `koanf:"jwt"`
	PostgresDB           database.Config   `koanf:"postgres_db"`
	Logger               logger.Config     `koanf:"logger"`
	TotalShutdownTimeout time.Duration     `koanf:"total_shutdown_timeout"`
	PathOfMigration      string            `koanf:"path_of_migration"`
}
