-- Create coupon_claims table
CREATE TABLE IF NOT EXISTS coupon_claims (
    id INTEGER NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    coupon_name VARCHAR NOT NULL,
    user_id VARCHAR NOT NULL
);

-- Create unique composite index on coupon_name and user_id to prevent duplicate claims
CREATE UNIQUE INDEX idx_coupon_claims_coupon_name_user_id ON coupon_claims(coupon_name, user_id);