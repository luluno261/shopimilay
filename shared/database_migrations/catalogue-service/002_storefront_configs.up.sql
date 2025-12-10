CREATE TABLE IF NOT EXISTS storefront_configs (
    merchant_id UUID PRIMARY KEY,
    sections JSONB,
    theme JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_storefront_configs_merchant_id ON storefront_configs(merchant_id);

