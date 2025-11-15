-- 000011_moderation_and_metrics.up.sql

CREATE TABLE IF NOT EXISTS moderation_queue (
  id            ulid PRIMARY KEY,
  target_type   moderation_target_type NOT NULL,
  target_ref    JSONB NOT NULL,  -- {db:'pg|mongo', table/collection:'...', id:'...'}
  reason        moderation_reason NOT NULL,
  flags         JSONB NOT NULL DEFAULT '{}'::jsonb, -- {toxicity:..., spam:..., engine:...}
  status        moderation_status NOT NULL DEFAULT 'pending',
  assigned_to   ulid NULL REFERENCES users(id) ON DELETE SET NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  resolved_at   TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS ix_moderation_queue_status_created
  ON moderation_queue (status, created_at);

CREATE INDEX IF NOT EXISTS ix_moderation_queue_target
  ON moderation_queue (target_type);

CREATE TABLE IF NOT EXISTS metrics_daily (
  date     DATE NOT NULL,
  metric   metric_type NOT NULL,
  key      VARCHAR(120) NOT NULL,
  locale   language NOT NULL,
  count    BIGINT NOT NULL DEFAULT 0 CHECK (count >= 0),
  PRIMARY KEY (date, metric, key, locale)
);