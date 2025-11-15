-- 000002_common_enums_and_functions.up.sql

-- Common enums
DO $$
BEGIN
  -- User & Auth
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status') THEN
    CREATE TYPE user_status AS ENUM ('active','suspended','deleted');
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'token_type') THEN
    CREATE TYPE token_type AS ENUM ('refresh','email_verify','password_reset');
  END IF;

  -- UI & Locale
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'theme') THEN
    CREATE TYPE theme AS ENUM ('light','dark','system');
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'language') THEN
    CREATE TYPE language AS ENUM ('fa','en');
  END IF;
END
$$;

-- updated_at trigger function
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;