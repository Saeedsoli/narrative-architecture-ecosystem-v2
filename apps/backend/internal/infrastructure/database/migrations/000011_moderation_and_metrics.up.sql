-- 000011_moderation_and_metrics.down.sql

DROP TABLE IF EXISTS metrics_daily;

DROP INDEX IF EXISTS ix_moderation_queue_target;
DROP INDEX IF EXISTS ix_moderation_queue_status_created;
DROP TABLE IF EXISTS moderation_queue;