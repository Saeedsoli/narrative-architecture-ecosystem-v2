-- 000005_products_coupons_orders_items.down.sql

-- Drop in reverse dependency order
DROP INDEX IF EXISTS ix_order_items_order;
DROP TABLE IF EXISTS order_items;

DROP INDEX IF EXISTS ix_orders_status_created;
DROP INDEX IF EXISTS ix_orders_user_created;
DROP TABLE IF EXISTS orders;

DROP INDEX IF EXISTS ix_coupon_redemptions_coupon;
DROP TABLE IF EXISTS coupon_redemptions;

DROP TRIGGER IF EXISTS trg_coupons_updated_at ON coupons;
DROP INDEX IF EXISTS ix_coupons_active_window;
DROP TABLE IF EXISTS coupons;

DROP TRIGGER IF EXISTS trg_products_updated_at ON products;
DROP INDEX IF EXISTS ix_products_type_active;
DROP TABLE IF EXISTS products;