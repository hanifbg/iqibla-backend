-- Migration: add_shipping_district_column
-- Description: Adds shipping_district column to orders table
-- Created: 2025-07-31

-- Up Migration
ALTER TABLE orders 
ADD COLUMN shipping_district VARCHAR(100);

-- Add index for location-based queries
CREATE INDEX idx_orders_shipping_location 
ON orders (shipping_province, shipping_city, shipping_district);

-- Set default value for existing records
UPDATE orders 
SET shipping_district = 'TBD' 
WHERE shipping_district IS NULL;

-- Down Migration
DROP INDEX IF EXISTS idx_orders_shipping_location;
ALTER TABLE orders DROP COLUMN IF EXISTS shipping_district;
