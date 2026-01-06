-- Create coupons table
CREATE TABLE IF NOT EXISTS coupons (
    id INTEGER NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR NOT NULL,
    amount INTEGER NOT NULL CHECK (amount >= 0),
    remaining_amount INTEGER NOT NULL CHECK (remaining_amount >= 0)
);

-- Create unique index on name
CREATE UNIQUE INDEX idx_coupons_name ON coupons(name);