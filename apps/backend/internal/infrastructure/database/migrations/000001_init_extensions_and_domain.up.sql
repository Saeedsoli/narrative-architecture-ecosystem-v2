-- 000001_init_extensions_and_domain.up.sql

-- Extensions: citext for case-insensitive text, pg_trgm for similarity indexes (future use)
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- ULID domain: 26 chars, Crockford Base32, uppercase
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_type t
    JOIN pg_namespace n ON n.oid = t.typnamespace
    WHERE t.typname = 'ulid' AND n.nspname = 'public'
  ) THEN
    CREATE DOMAIN ulid AS CHAR(26)
      CHECK (VALUE ~ '^[0-9A-HJKMNP-TV-Z]{26}$' AND VALUE = UPPER(VALUE));
  END IF;
END
$$;