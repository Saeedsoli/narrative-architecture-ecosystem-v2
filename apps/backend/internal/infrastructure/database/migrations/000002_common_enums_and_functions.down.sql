-- 000002_common_enums_and_functions.down.sql

-- Drop trigger function
DROP FUNCTION IF EXISTS set_updated_at() CASCADE;

-- Drop enums (only when no dependent columns exist)
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'language') THEN
    DROP TYPE language;
  END IF;
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'theme') THEN
    DROP TYPE theme;
  END IF;
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'token_type') THEN
    DROP TYPE token_type;
  END IF;
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status') THEN
    DROP TYPE user_status;
  END IF;
END
$$;