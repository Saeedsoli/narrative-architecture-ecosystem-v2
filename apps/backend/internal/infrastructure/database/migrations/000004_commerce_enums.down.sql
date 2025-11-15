-- 000004_commerce_enums.down.sql

DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'entitlement_status') THEN
    DROP TYPE entitlement_status;
  END IF;
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'entitlement_source') THEN
    DROP TYPE entitlement_source;
  END IF;
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'entitlement_resource_type') THEN
    DROP TYPE entitlement_resource_type;
  END IF;

  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'subscription_status') THEN
    DROP TYPE subscription_status;
  END IF;
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'plan_type') THEN
    DROP TYPE plan_type;
  END IF;

  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_status') THEN
    DROP TYPE payment_status;
  END IF;
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_gateway') THEN
    DROP TYPE payment_gateway;
  END IF;

  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
    DROP TYPE order_status;
  END IF;
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'product_type') THEN
    DROP TYPE product_type;
  END IF;
END
$$;