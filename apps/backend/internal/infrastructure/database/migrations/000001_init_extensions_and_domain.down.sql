-- 000001_init_extensions_and_domain.down.sql

-- Drop ULID domain (ensure no table depends on it; run after dropping dependent tables)
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM pg_type t
    JOIN pg_namespace n ON n.oid = t.typnamespace
    WHERE t.typname = 'ulid' AND n.nspname = 'public'
  ) THEN
    DROP DOMAIN ulid;
  END IF;
END
$$;

-- Extensions can be left installed; drop if you need a clean rollback
-- DROP EXTENSION IF EXISTS pg_trgm;
-- DROP EXTENSION IF EXISTS citext;