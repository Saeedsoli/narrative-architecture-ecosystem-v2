-- 000009_ai_logs_submissions_peer_reviews.up.sql

-- AI logs (detailed, 365-day retention handled by jobs)
CREATE TABLE IF NOT EXISTS ai_logs (
  id             ulid PRIMARY KEY,
  user_id        ulid NULL REFERENCES users(id) ON DELETE SET NULL,
  model          VARCHAR(64) NOT NULL,
  prompt         JSONB NOT NULL,
  retrieved_refs JSONB NULL,
  response       JSONB NOT NULL,
  tokens_input   INT NULL CHECK (tokens_input IS NULL OR tokens_input >= 0),
  tokens_output  INT NULL CHECK (tokens_output IS NULL OR tokens_output >= 0),
  cost_cents     INT NULL CHECK (cost_cents IS NULL OR cost_cents >= 0),
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ix_ai_logs_user_created
  ON ai_logs (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS ix_ai_logs_model_created
  ON ai_logs (model, created_at DESC);

-- Submissions
CREATE TABLE IF NOT EXISTS submissions (
  id            ulid PRIMARY KEY,
  exercise_id   ulid NOT NULL REFERENCES exercises(id) ON DELETE RESTRICT,
  user_id       ulid NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
  status        submission_status NOT NULL DEFAULT 'pending',
  answer        JSONB NOT NULL,
  score         INT NULL,
  feedback      TEXT NULL,
  ai_summary    JSONB NULL,
  ai_log_id     ulid NULL REFERENCES ai_logs(id) ON DELETE SET NULL,
  submitted_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  graded_at     TIMESTAMPTZ NULL,
  deleted_at    TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS ix_submissions_exercise_user
  ON submissions (exercise_id, user_id);

CREATE INDEX IF NOT EXISTS ix_submissions_user_submitted
  ON submissions (user_id, submitted_at DESC);

CREATE INDEX IF NOT EXISTS ix_submissions_status_submitted
  ON submissions (status, submitted_at DESC);

-- Peer Reviews (feature-flaggable)
CREATE TABLE IF NOT EXISTS peer_reviews (
  id             ulid PRIMARY KEY,
  submission_id  ulid NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
  reviewer_id    ulid NOT NULL REFERENCES users(id)       ON DELETE CASCADE,
  rubric_scores  JSONB NOT NULL,
  comments       TEXT NULL,
  status         peer_review_status NOT NULL DEFAULT 'pending',
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (submission_id, reviewer_id)
);

CREATE TRIGGER trg_peer_reviews_updated_at
  BEFORE UPDATE ON peer_reviews
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();