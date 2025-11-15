-- 000008_exercises_and_rubrics.up.sql

CREATE TABLE IF NOT EXISTS exercises (
  id            ulid PRIMARY KEY,
  chapter_id    VARCHAR(64) NOT NULL, -- e.g. 'chapter-01'
  type          exercise_type NOT NULL,
  difficulty    exercise_difficulty NOT NULL,
  points        INT NOT NULL CHECK (points >= 0),
  content       JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at    TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS ix_exercises_chapter
  ON exercises (chapter_id);

CREATE INDEX IF NOT EXISTS ix_exercises_type_difficulty
  ON exercises (type, difficulty);

CREATE TRIGGER trg_exercises_updated_at
  BEFORE UPDATE ON exercises
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS exercise_rubrics (
  id          ulid PRIMARY KEY,
  exercise_id ulid NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
  criteria    JSONB NOT NULL, -- [{name, points, description}]
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ix_exercise_rubrics_exercise
  ON exercise_rubrics (exercise_id);