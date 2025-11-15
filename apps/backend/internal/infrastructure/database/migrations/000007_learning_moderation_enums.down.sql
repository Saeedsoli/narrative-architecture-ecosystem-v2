-- 000007_learning_moderation_enums.down.sql

DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'metric_type') THEN
    DROP TYPE metric_type;
  END IF;

  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'moderation_status') THEN
    DROP TYPE moderation_status;
  END IF;
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'moderation_reason') THEN
    DROP TYPE moderation_reason;
  END IF;
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'moderation_target_type') THEN
    DROP TYPE moderation_target_type;
  END IF;

  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'peer_review_status') THEN
    DROP TYPE peer_review_status;
  END IF;
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'submission_status') THEN
    DROP TYPE submission_status;
  END IF;
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'exercise_difficulty') THEN
    DROP TYPE exercise_difficulty;
  END IF;
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'exercise_type') THEN
    DROP TYPE exercise_type;
  END IF;
END
$$;