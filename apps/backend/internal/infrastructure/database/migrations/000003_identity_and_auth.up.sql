-- 000003_identity_and_auth.up.sql

-- users
CREATE TABLE IF NOT EXISTS users (
  id            ulid PRIMARY KEY,
  email         CITEXT NOT NULL,
  email_verified BOOLEAN NOT NULL DEFAULT FALSE,
  password_hash TEXT NOT NULL,
  status        user_status NOT NULL DEFAULT 'active',
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at    TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_users_email_active
  ON users (email) WHERE deleted_at IS NULL;

CREATE TRIGGER trg_users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- user_profiles (1:1)
CREATE TABLE IF NOT EXISTS user_profiles (
  user_id      ulid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  full_name    VARCHAR(120) NOT NULL,
  username     VARCHAR(32)  NOT NULL,
  phone        VARCHAR(32)  NULL,
  phone_verified BOOLEAN    NOT NULL DEFAULT FALSE,
  birth_date   DATE         NULL,
  avatar_url   TEXT         NULL,
  language     language     NOT NULL DEFAULT 'fa',
  timezone     VARCHAR(64)  NOT NULL DEFAULT 'Asia/Tehran',
  created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  CONSTRAINT chk_user_profiles_username_format CHECK (username ~ '^[A-Za-z0-9._-]{3,32}$'),
  CONSTRAINT chk_user_profiles_phone_e164      CHECK (phone IS NULL OR phone ~ '^\+[1-9]\d{1,14}$')
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_user_profiles_username_lower
  ON user_profiles ((lower(username)));

CREATE TRIGGER trg_user_profiles_updated_at
  BEFORE UPDATE ON user_profiles
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- user_preferences (1:1)
CREATE TABLE IF NOT EXISTS user_preferences (
  user_id         ulid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  persian_numbers BOOLEAN NOT NULL DEFAULT TRUE,
  shamsi_dates    BOOLEAN NOT NULL DEFAULT TRUE,
  theme           theme   NOT NULL DEFAULT 'system',
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_user_preferences_updated_at
  BEFORE UPDATE ON user_preferences
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- roles
CREATE TABLE IF NOT EXISTS roles (
  id         ulid PRIMARY KEY,
  name       CITEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- user_roles (N:M)
CREATE TABLE IF NOT EXISTS user_roles (
  user_id    ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id    ulid NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, role_id)
);

-- auth_tokens
CREATE TABLE IF NOT EXISTS auth_tokens (
  id          ulid PRIMARY KEY,
  user_id     ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type        token_type NOT NULL,
  token_hash  TEXT NOT NULL,
  user_agent  TEXT NULL,
  ip          INET NULL,
  expires_at  TIMESTAMPTZ NOT NULL,
  used_at     TIMESTAMPTZ NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_auth_tokens_token_hash
  ON auth_tokens (token_hash);

CREATE INDEX IF NOT EXISTS ix_auth_tokens_active
  ON auth_tokens (user_id, type, expires_at)
  WHERE used_at IS NULL AND expires_at > NOW();

-- notifications
CREATE TABLE IF NOT EXISTS notifications (
  id         ulid PRIMARY KEY,
  user_id    ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type       VARCHAR(64) NOT NULL,
  payload    JSONB NOT NULL,
  read_at    TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ix_notifications_user_created
  ON notifications (user_id, created_at DESC);

-- audit_log
CREATE TABLE IF NOT EXISTS audit_log (
  id           ulid PRIMARY KEY,
  actor_id     ulid NULL REFERENCES users(id) ON DELETE SET NULL,
  action       VARCHAR(128) NOT NULL,
  target_ref   JSONB NOT NULL, -- {db:'pg|mongo', table/collection:'...', id:'...'}
  data_before  JSONB NULL,
  data_after   JSONB NULL,
  ip           INET NULL,
  user_agent   TEXT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);