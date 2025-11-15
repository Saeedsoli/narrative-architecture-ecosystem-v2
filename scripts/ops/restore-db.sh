#!/usr/bin/env bash
# scripts/ops/restore-db.sh

set -e

S3_BUCKET="s3://narrative-arch-backups"
PG_FILE=$1
MONGO_FILE=$2

if [ -z "$PG_FILE" ] || [ -z "$MONGO_FILE" ]; then
  echo "Usage: $0 <postgres_backup_file.sql.gz> <mongo_backup_file.archive.gz>"
  exit 1
fi

# --- PostgreSQL Restore ---
echo "Downloading PostgreSQL backup: $PG_FILE"
aws s3 cp $S3_BUCKET/postgres/$PG_FILE /tmp/$PG_FILE

echo "Restoring PostgreSQL..."
gunzip < /tmp/$PG_FILE | psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB
echo "PostgreSQL restore completed."

# --- MongoDB Restore ---
echo "Downloading MongoDB backup: $MONGO_FILE"
aws s3 cp $S3_BUCKET/mongo/$MONGO_FILE /tmp/$MONGO_FILE

echo "Restoring MongoDB..."
mongorestore --uri="$MONGODB_URI" --archive --gzip < /tmp/$MONGO_FILE
echo "MongoDB restore completed."

# --- Cleanup ---
rm /tmp/$PG_FILE
rm /tmp/$MONGO_FILE

echo "✅ Database restore process completed successfully!"