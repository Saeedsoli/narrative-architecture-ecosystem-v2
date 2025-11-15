-- 000003_identity_and_auth.down.sql

-- Drop in dependency-safe order
DROP TABLE IF EXISTS audit_log;
DROP INDEX IF EXISTS ix_notifications_user_created;
DROP TABLE IF EXISTS notifications;
DROP INDEX IF EXISTS ix_auth_tokens_active;
DROP INDEX IF EXISTS ux_auth_tokens_token_hash;
DROP TABLE IF EXISTS auth_tokens;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;
DROP TRIGGER IF EXISTS trg_user_preferences_updated_at ON user_preferences;
DROP TABLE IF EXISTS user_preferences;
DROP TRIGGER IF EXISTS trg_user_profiles_updated_at ON user_profiles;
DROP INDEX IF EXISTS ux_user_profiles_username_lower;
DROP TABLE IF EXISTS user_profiles;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP INDEX IF EXISTS ux_users_email_active;
DROP TABLE IF EXISTS users;