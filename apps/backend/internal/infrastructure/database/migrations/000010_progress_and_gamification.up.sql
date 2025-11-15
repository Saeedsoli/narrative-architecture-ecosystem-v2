-- 000010_progress_and_gamification.up.sql

-- Progress per chapter
CREATE TABLE IF NOT EXISTS progress (
  id           ulid PRIMARY KEY,
  user_id      ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  chapter_id   VARCHAR(64) NOT NULL,
  percent      NUMERIC(5,2) NOT NULL CHECK (percent >= 0 AND percent <= 100),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (user_id, chapter_id)
);

-- User streaks (optional)
CREATE TABLE IF NOT EXISTS user_streaks (
  user_id          ulid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  current_streak   INT NOT NULL DEFAULT 0 CHECK (current_streak >= 0),
  longest_streak   INT NOT NULL DEFAULT 0 CHECK (longest_streak >= 0),
  last_activity_on DATE NULL
);

-- Badges (optional)
CREATE TABLE IF NOT EXISTS badges (
  id         ulid PRIMARY KEY,
  code       VARCHAR(64) NOT NULL UNIQUE,
  name       VARCHAR(120) NOT NULL,
  criteria   JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_badges (
  id         ulid PRIMARY KEY,
  user_id    ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  badge_id   ulid NOT NULL REFERENCES badges(id) ON DELETE RESTRICT,
  earned_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (user_id, badge_id)
);