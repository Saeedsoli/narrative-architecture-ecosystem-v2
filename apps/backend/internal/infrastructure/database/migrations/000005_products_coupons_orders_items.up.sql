-- 000005_products_coupons_orders_items.up.sql

-- 1) products
CREATE TABLE IF NOT EXISTS products (
  id             ulid PRIMARY KEY,
  type           product_type NOT NULL,
  sku            VARCHAR(64) NOT NULL UNIQUE,
  title          VARCHAR(200) NOT NULL,
  description    TEXT NOT NULL,
  price_cents    INT NOT NULL CHECK (price_cents >= 0),
  currency       VARCHAR(3) NOT NULL DEFAULT 'IRR' CHECK (currency = 'IRR'),
  active         BOOLEAN NOT NULL DEFAULT TRUE,
  meta           JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at     TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS ix_products_type_active
  ON products (type, active);

CREATE TRIGGER trg_products_updated_at
  BEFORE UPDATE ON products
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- 2) coupons
CREATE TABLE IF NOT EXISTS coupons (
  id            ulid PRIMARY KEY,
  code          VARCHAR(50) NOT NULL UNIQUE,
  type          VARCHAR(10) NOT NULL CHECK (type IN ('percent','fixed')),
  value_int     INT NOT NULL CHECK (
                  (type = 'percent' AND value_int BETWEEN 1 AND 100) OR
                  (type = 'fixed'   AND value_int > 0)
                ),
  starts_at     TIMESTAMPTZ NOT NULL,
  ends_at       TIMESTAMPTZ NOT NULL,
  usage_limit   INT NOT NULL CHECK (usage_limit >= 0),
  active        BOOLEAN NOT NULL DEFAULT TRUE,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at    TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS ix_coupons_active_window
  ON coupons (active, starts_at, ends_at);

CREATE TRIGGER trg_coupons_updated_at
  BEFORE UPDATE ON coupons
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- 3) coupon_redemptions
CREATE TABLE IF NOT EXISTS coupon_redemptions (
  id           ulid PRIMARY KEY,
  coupon_id    ulid NOT NULL REFERENCES coupons(id) ON DELETE RESTRICT,
  user_id      ulid NOT NULL REFERENCES users(id)   ON DELETE CASCADE,
  order_id     ulid NOT NULL,
  redeemed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (coupon_id, user_id, order_id)
);

CREATE INDEX IF NOT EXISTS ix_coupon_redemptions_coupon
  ON coupon_redemptions (coupon_id, redeemed_at);

-- 4) orders (snapshots for shipping info may be NULL for digital-only orders)
CREATE TABLE IF NOT EXISTS orders (
  id                 ulid PRIMARY KEY,
  user_id            ulid NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  status             order_status NOT NULL DEFAULT 'pending',
  total_cents        INT NOT NULL CHECK (total_cents >= 0),
  shipping_cents     INT NOT NULL DEFAULT 0 CHECK (shipping_cents >= 0),
  discount_cents     INT NOT NULL DEFAULT 0 CHECK (discount_cents >= 0),
  currency           VARCHAR(3) NOT NULL DEFAULT 'IRR' CHECK (currency = 'IRR'),
  shipping_name      VARCHAR(120) NULL,
  shipping_phone     VARCHAR(32)  NULL,
  shipping_address_line TEXT      NULL,
  city               VARCHAR(100) NULL,
  province           VARCHAR(100) NULL,
  postal_code        VARCHAR(20)  NULL,
  notes              TEXT NULL,
  created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  paid_at            TIMESTAMPTZ NULL,
  canceled_at        TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS ix_orders_user_created
  ON orders (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS ix_orders_status_created
  ON orders (status, created_at DESC);

-- 5) order_items
CREATE TABLE IF NOT EXISTS order_items (
  id                 ulid PRIMARY KEY,
  order_id           ulid NOT NULL REFERENCES orders(id)   ON DELETE CASCADE,
  product_id         ulid NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
  quantity           INT NOT NULL CHECK (quantity >= 1),
  unit_price_cents   INT NOT NULL CHECK (unit_price_cents >= 0),
  total_price_cents  INT GENERATED ALWAYS AS (quantity * unit_price_cents) STORED,
  meta               JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS ix_order_items_order
  ON order_items (order_id);