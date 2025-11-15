-- 000004_commerce_enums.up.sql

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'product_type') THEN
    CREATE TYPE product_type AS ENUM ('physical_book','premium_content','subscription_monthly','course');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
    CREATE TYPE order_status AS ENUM ('pending','paid','failed','canceled');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_gateway') THEN
    CREATE TYPE payment_gateway AS ENUM ('zarinpal');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_status') THEN
    CREATE TYPE payment_status AS ENUM ('init','pending','success','failed');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'plan_type') THEN
    CREATE TYPE plan_type AS ENUM ('monthly');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'subscription_status') THEN
    CREATE TYPE subscription_status AS ENUM ('active','past_due','canceled','expired');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'entitlement_resource_type') THEN
    CREATE TYPE entitlement_resource_type AS ENUM ('chapter','premium_content','course','subscription_feature');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'entitlement_source') THEN
    CREATE TYPE entitlement_source AS ENUM ('order','qr_register','manual');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'entitlement_status') THEN
    CREATE TYPE entitlement_status AS ENUM ('active','revoked');
  END IF;
END
$$;