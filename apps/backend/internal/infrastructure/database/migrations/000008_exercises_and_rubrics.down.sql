-- 000008_exercises_and_rubrics.down.sql

DROP INDEX IF EXISTS ix_exercise_rubrics_exercise;
DROP TABLE IF EXISTS exercise_rubrics;

DROP TRIGGER IF EXISTS trg_exercises_updated_at ON exercises;
DROP INDEX IF EXISTS ix_exercises_type_difficulty;
DROP INDEX IF EXISTS ix_exercises_chapter;
DROP TABLE IF EXISTS exercises;