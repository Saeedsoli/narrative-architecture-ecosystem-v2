// apps/backend/internal/infrastructure/database/postgres/connection.go

package postgres

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq" // درایور PostgreSQL
)

// DBConfig شامل تنظیمات اتصال به دیتابیس است.
type DBConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// ConnectDB یک Connection Pool به PostgreSQL ایجاد و برمی‌گرداند.
func ConnectDB(cfg DBConfig) *sql.DB {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		log.Fatalf("FATAL: Failed to open database connection: %v", err)
	}

	// تنظیمات Connection Pool برای عملکرد بهینه
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// بررسی صحت اتصال
	if err := db.Ping(); err != nil {
		log.Fatalf("FATAL: Failed to ping database: %v", err)
	}

	log.Println("PostgreSQL connection established successfully.")
	return db
}