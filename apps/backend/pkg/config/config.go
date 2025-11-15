// apps/backend/pkg/config/config.go

package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config ساختار کلی پیکربندی پروژه است.
type Config struct {
	Server   ServerConfig
	DB       DBConfig
	Mongo    MongoConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Services ServiceConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DBConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type MongoConfig struct {
	URI    string
	DBName string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type ServiceConfig struct {
	AI               string
	S3Bucket         string
	S3Region         string
	DecryptionKey    string
}

// LoadConfig پیکربندی را از متغیرهای محیطی می‌خواند.
func LoadConfig() *Config {
	vp := viper.New()
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vp.AutomaticEnv()

	// تنظیم مقادیر پیش‌فرض
	vp.SetDefault("server.port", "8080")
	vp.SetDefault("server.ginmode", "debug")
	vp.SetDefault("db.maxopenconns", 25)
	vp.SetDefault("db.maxidleconns", 25)
	vp.SetDefault("db.connmaxlifetime", 5*time.Minute)
	vp.SetDefault("jwt.accessttl", 15*time.Minute)
	vp.SetDefault("jwt.refreshttl", 30*24*time.Hour)

	var cfg Config
	if err := vp.Unmarshal(&cfg); err != nil {
		log.Fatalf("FATAL: unable to decode config into struct, %v", err)
	}

	return &cfg
}