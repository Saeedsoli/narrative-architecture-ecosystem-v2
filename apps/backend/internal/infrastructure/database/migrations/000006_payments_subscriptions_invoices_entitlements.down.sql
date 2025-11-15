-- 000006_payments_subscriptions_invoices_entitlements.down.sql

-- Drop in reverse dependency order
DROP TRIGGER IF EXISTS trg_entitlements_updated_at ON entitlements;
DROP INDEX IF EXISTS ix_entitlements_user;
DROP INDEX IF EXISTS ux_entitlements_unique_active;
DROP TABLE IF EXISTS entitlements;

DROP INDEX IF EXISTS ix_invoices_order;
DROP TABLE IF EXISTS invoices;

DROP TRIGGER IF EXISTS trg_subscriptions_updated_at ON subscriptions;
DROP INDEX IF EXISTS ix_subscriptions_status_end;
DROP TABLE IF EXISTS subscriptions;

DROP TRIGGER IF EXISTS trg_payments_updated_at ON payments;
DROP INDEX IF EXISTS ix_payments_order_status;
DROP INDEX IF EXISTS ux_payments_authority;
DROP TABLE IF EXISTS payments;