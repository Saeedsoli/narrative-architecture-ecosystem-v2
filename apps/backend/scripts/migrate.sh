#!/bin/bash
# apps/backend/scripts/migrate.sh

set -e

# مطمئن شوید که golang-migrate نصب است
if ! command -v migrate &> /dev/null
then
    echo "migrate could not be found. Please install it:"
    echo "go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

# خواندن متغیرهای محیطی از فایل .env در صورت وجود
if [ -f .env ]; then
  export $(cat .env | sed 's/#.*//g' | xargs)
fi

# بررسی وجود DATABASE_URL
if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL is not set."
    exit 1
fi

MIGRATIONS_PATH="file://internal/infrastructure/database/migrations"

# اجرای دستور بر اساس آرگومان اول
case "$1" in
  up)
    echo "Running migrations up..."
    migrate -database "$DATABASE_URL" -path "$MIGRATIONS_PATH" up
    ;;
  down)
    echo "Running migrations down..."
    migrate -database "$DATABASE_URL" -path "$MIGRATIONS_PATH" down
    ;;
  create)
    if [ -z "$2" ]; then
        echo "Error: Please provide a name for the migration."
        echo "Usage: ./scripts/migrate.sh create <migration_name>"
        exit 1
    fi
    echo "Creating migration: $2"
    migrate create -ext sql -dir internal/infrastructure/database/migrations -seq "$2"
    ;;
  *)
    echo "Usage: ./scripts/migrate.sh {up|down|create <name>}"
    exit 1
    ;;
esac

echo "Migration script finished."