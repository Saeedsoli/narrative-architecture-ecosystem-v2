-- 000006_payments_subscriptions_invoices_entitlements.up.sql

-- 1) payments
CREATE TABLE IF NOT EXISTS payments (
  id                 ulid PRIMARY KEY,
  order_id           ulid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  gateway            payment_gateway NOT NULL DEFAULT 'zarinpal',
  amount_cents       INT NOT NULL CHECK (amount_cents >= 0),
  status             payment_status NOT NULL DEFAULT 'init',
  authority          VARCHAR(128) NULL,
  ref_id             VARCHAR(128) NULL,
  callback_payload   JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_payments_authority
  ON payments (authority) WHERE authority IS NOT NULL;

CREATE INDEX IF NOT EXISTS ix_payments_order_status
  ON payments (order_id, status);

CREATE TRIGGER trg_payments_updated_at
  BEFORE UPDATE ON payments
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- 2) subscriptions (one row per user assumed; monthly manual renewal with 7d grace)
CREATE TABLE IF NOT EXISTS subscriptions (
  id                      ulid PRIMARY KEY,
  user_id                 ulid NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  plan                    plan_type NOT NULL DEFAULT 'monthly',
  status                  subscription_status NOT NULL DEFAULT 'active',
  current_period_start    TIMESTAMPTZ NOT NULL,
  current_period_end      TIMESTAMPTZ NOT NULL,
  grace_until             TIMESTAMPTZ NULL,
  origin_order_id         ulid NULL REFERENCES orders(id) ON DELETE SET NULL,
  last_order_id           ulid NULL REFERENCES orders(id) ON DELETE SET NULL,
  created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ix_subscriptions_status_end
  ON subscriptions (status, current_period_end DESC);

CREATE TRIGGER trg_subscriptions_updated_at
  BEFORE UPDATE ON subscriptions
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- 3) invoices (simple)
CREATE TABLE IF NOT EXISTS invoices (
  id               ulid PRIMARY KEY,
  order_id         ulid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  invoice_number   VARCHAR(64) NOT NULL UNIQUE,
  pdf_url          TEXT NOT NULL,
  totals           JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ix_invoices_order
  ON invoices (order_id, created_at DESC);

-- 4) entitlements (digital access control: chapter16, premium content, courses)
CREATE TABLE IF NOT EXISTS entitlements (
  id              ulid PRIMARY KEY,
  user_id         ulid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  resource_type   entitlement_resource_type NOT NULL,
  resource_id     VARCHAR(120) NOT NULL,
  source          entitlement_source NOT NULL,
  status          entitlement_status NOT NULL DEFAULT 'active',
  meta            JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_entitlements_unique_active
  ON entitlements (user_id, resource_type, resource_id) WHERE status = 'active';

CREATE INDEX IF NOT EXISTS ix_entitlements_user
  ON entitlements (user_id, created_at DESC);

CREATE TRIGGER trg_entitlements_updated_at
  BEFORE UPDATE ON entitlements
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();