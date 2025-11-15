// apps/backend/test/integration/main_test.go

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	pgDB    *sql.DB
	mongoDB *mongo.Database
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// --- 1. راه‌اندازی کانتینر PostgreSQL ---
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute),
		),
	)
	if err != nil {
		log.Fatalf("Could not start postgres container: %s", err)
	}

	// --- 2. راه‌اندازی کانتینر MongoDB ---
	mongoContainer, err := mongodb.RunContainer(ctx,
		testcontainers.WithImage("mongo:7"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Waiting for connections").
				WithOccurrence(1).
				WithStartupTimeout(5*time.Minute),
		),
	)
	if err != nil {
		log.Fatalf("Could not start mongo container: %s", err)
	}

	// --- 3. اتصال به دیتابیس‌ها ---
	pgConnStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
	pgDB, err = sql.Open("postgres", pgConnStr)
	if err != nil {
		log.Fatalf("Could not connect to postgres: %s", err)
	}

	mongoConnStr, _ := mongoContainer.ConnectionString(ctx)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConnStr))
	if err != nil {
		log.Fatalf("Could not connect to mongo: %s", err)
	}
	mongoDB = mongoClient.Database("test_content_db")

	// --- 4. اجرای مایگریشن‌ها ---
	mig, err := migrate.New("file://../../internal/infrastructure/database/migrations", pgConnStr)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %s", err)
	}
	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not run up migrations: %s", err)
	}
	
	log.Println("Test databases and migrations are ready.")

	// --- 5. اجرای تمام تست‌ها ---
	exitCode := m.Run()

	// --- 6. پاک‌سازی و خاموش کردن کانتینرها ---
	log.Println("Terminating test containers...")
	if err := pgContainer.Terminate(ctx); err != nil {
		log.Fatalf("Could not terminate postgres container: %s", err)
	}
	if err := mongoContainer.Terminate(ctx); err != nil {
		log.Fatalf("Could not terminate mongo container: %s", err)
	}

	os.Exit(exitCode)
}