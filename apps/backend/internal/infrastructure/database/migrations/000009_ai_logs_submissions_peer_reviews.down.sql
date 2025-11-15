-- 000009_ai_logs_submissions_peer_reviews.down.sql

DROP TRIGGER IF EXISTS trg_peer_reviews_updated_at ON peer_reviews;
DROP TABLE IF EXISTS peer_reviews;

DROP INDEX IF EXISTS ix_submissions_status_submitted;
DROP INDEX IF EXISTS ix_submissions_user_submitted;
DROP INDEX IF EXISTS ix_submissions_exercise_user;
DROP TABLE IF EXISTS submissions;

DROP INDEX IF EXISTS ix_ai_logs_model_created;
DROP INDEX IF EXISTS ix_ai_logs_user_created;
DROP TABLE IF EXISTS ai_logs;