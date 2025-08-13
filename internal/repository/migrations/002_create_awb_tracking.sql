-- Migration: Create AWB tracking table
-- Purpose: Store AWB numbers linked to orders for tracking shipments

CREATE TABLE IF NOT EXISTS awb_tracking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    awb_number VARCHAR(255) NOT NULL,
    courier VARCHAR(50) NOT NULL,
    last_phone_number VARCHAR(5),
    is_validated BOOLEAN DEFAULT FALSE,
    tracking_data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    -- Constraints
    CONSTRAINT fk_awb_tracking_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT uk_awb_tracking_awb_courier UNIQUE (awb_number, courier),
    CONSTRAINT ck_awb_tracking_courier CHECK (courier IN ('jne', 'jnt', 'ninja', 'tiki', 'pos', 'anteraja', 'sicepat', 'sap', 'lion', 'wahana', 'first', 'ide'))
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_awb_tracking_order_id ON awb_tracking(order_id);
CREATE INDEX IF NOT EXISTS idx_awb_tracking_awb_number ON awb_tracking(awb_number);
CREATE INDEX IF NOT EXISTS idx_awb_tracking_courier ON awb_tracking(courier);
CREATE INDEX IF NOT EXISTS idx_awb_tracking_created_at ON awb_tracking(created_at);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_awb_tracking_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_awb_tracking_updated_at
    BEFORE UPDATE ON awb_tracking
    FOR EACH ROW
    EXECUTE FUNCTION update_awb_tracking_updated_at();
