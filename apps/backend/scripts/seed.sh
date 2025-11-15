#!/bin/bash
# apps/backend/scripts/seed.sh

set -e

# خواندن متغیرهای محیطی
if [ -f .env ]; then
  export $(cat .env | sed 's/#.*//g' | xargs)
fi

if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL is not set."
    exit 1
fi

echo "Seeding initial data..."

# اجرای مایگریشن مربوط به seed کردن نقش‌ها
# ما این کار را در یک فایل مایگریشن جداگانه انجام داده‌ایم (000012_seed_roles.up.sql)
# بنابراین، فقط کافی است مایگریشن‌ها را اجرا کنیم.
echo "Running 'roles' seed migration..."
migrate -database "$DATABASE_URL" -path "file://internal/infrastructure/database/migrations" -verbose up

# اگر نیاز به seed کردن داده‌های دیگر باشد، می‌توان از یک اسکریپت Go استفاده کرد:
# echo "Running additional Go seeders..."
# go run ./cmd/seeder/main.go

echo "✅ Seeding completed."