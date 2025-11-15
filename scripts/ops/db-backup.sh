#!/usr/bin/env bash
# scripts/ops/db-backup.sh

set -e

# --- Configuration ---
S3_BUCKET="s3://narrative-arch-backups"
DATE=$(date +%Y-%m-%d-%H%M)
PG_BACKUP_FILE="postgres-backup-${DATE}.sql.gz"
MONGO_BACKUP_FILE="mongo-backup-${DATE}.archive.gz"

# --- PostgreSQL Backup ---
echo "Starting PostgreSQL backup..."
pg_dump -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -C -x --if-exists --clean \
  | gzip > /tmp/$PG_BACKUP_FILE
echo "PostgreSQL backup created: /tmp/$PG_BACKUP_FILE"

# --- MongoDB Backup ---
echo "Starting MongoDB backup..."
mongodump --uri="$MONGODB_URI" --archive --gzip > /tmp/$MONGO_BACKUP_FILE
echo "MongoDB backup created: /tmp/$MONGO_BACKUP_FILE"

# --- Upload to S3 ---
echo "Uploading backups to S3 bucket: $S3_BUCKET"
aws s3 cp /tmp/$PG_BACKUP_FILE $S3_BUCKET/postgres/
aws s3 cp /tmp/$MONGO_BACKUP_FILE $S3_BUCKET/mongo/

# --- Cleanup ---
rm /tmp/$PG_BACKUP_FILE
rm /tmp/$MONGO_BACKUP_FILE

echo "✅ Database backup process completed successfully!"