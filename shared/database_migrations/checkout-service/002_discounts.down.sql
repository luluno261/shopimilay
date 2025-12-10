-- Rollback migration pour discounts

DROP INDEX IF EXISTS idx_discounts_is_active;
DROP INDEX IF EXISTS idx_discounts_code;
DROP INDEX IF EXISTS idx_discounts_merchant_id;
DROP TABLE IF EXISTS discounts;

