/*
  # Catalogue Service - Initial Schema

  1. New Tables
    - `products` - Product catalog
    - `product_variants` - Product variants (size, color, etc.)
    - `product_categories` - Product categories
    - `product_images` - Product images
    - `inventory` - Stock management
    - `storefront_configs` - Storefront visual configuration

  2. Security
    - Enable RLS for multi-tenant data isolation
    - Merchants can only access their own products
    - Public can read published products

  3. Implementation Notes
    - Supports product variants for flexible inventory
    - Elasticsearch integration for search
    - Storefront configuration per merchant
*/

-- Products table
CREATE TABLE IF NOT EXISTS products (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid NOT NULL,
  name text NOT NULL,
  slug text NOT NULL,
  description text,
  short_description text,
  sku text UNIQUE NOT NULL,
  price numeric(12,2) NOT NULL CHECK (price >= 0),
  compare_at_price numeric(12,2),
  cost numeric(12,2),
  currency text DEFAULT 'USD',
  status text DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
  visibility text DEFAULT 'public' CHECK (visibility IN ('public', 'private', 'draft')),
  category_id uuid,
  collection_ids uuid[] DEFAULT '{}'::uuid[],
  tags text[] DEFAULT '{}'::text[],
  weight numeric(10,3),
  weight_unit text DEFAULT 'kg',
  requires_shipping boolean DEFAULT true,
  track_inventory boolean DEFAULT true,
  allow_backorder boolean DEFAULT false,
  seo_title text,
  seo_description text,
  seo_keywords text[] DEFAULT '{}'::text[],
  rating numeric(3,2) DEFAULT 0,
  rating_count integer DEFAULT 0,
  vendor text,
  metadata jsonb DEFAULT '{}'::jsonb,
  published_at timestamptz,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  CONSTRAINT merchant_sku_unique UNIQUE (merchant_id, sku)
);

-- Product variants
CREATE TABLE IF NOT EXISTS product_variants (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  merchant_id uuid NOT NULL,
  title text,
  sku text UNIQUE NOT NULL,
  price numeric(12,2) NOT NULL CHECK (price >= 0),
  compare_at_price numeric(12,2),
  cost numeric(12,2),
  weight numeric(10,3),
  weight_unit text DEFAULT 'kg',
  barcode text,
  option_1 text,
  option_2 text,
  option_3 text,
  image_id uuid,
  tax_class text,
  tax_code text,
  status text DEFAULT 'active' CHECK (status IN ('active', 'archived')),
  requires_shipping boolean DEFAULT true,
  track_inventory boolean DEFAULT true,
  position integer DEFAULT 1,
  metadata jsonb DEFAULT '{}'::jsonb,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  CONSTRAINT merchant_variant_sku UNIQUE (merchant_id, sku)
);

-- Product categories
CREATE TABLE IF NOT EXISTS product_categories (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid NOT NULL,
  name text NOT NULL,
  slug text NOT NULL,
  description text,
  image_url text,
  parent_category_id uuid REFERENCES product_categories(id) ON DELETE SET NULL,
  position integer DEFAULT 0,
  seo_title text,
  seo_description text,
  status text DEFAULT 'active' CHECK (status IN ('active', 'archived')),
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  CONSTRAINT merchant_category_slug UNIQUE (merchant_id, slug)
);

-- Product images
CREATE TABLE IF NOT EXISTS product_images (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  merchant_id uuid NOT NULL,
  image_url text NOT NULL,
  alt_text text,
  position integer DEFAULT 0,
  is_featured boolean DEFAULT false,
  width integer,
  height integer,
  file_size integer,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- Inventory management
CREATE TABLE IF NOT EXISTS inventory (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  variant_id uuid NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,
  merchant_id uuid NOT NULL,
  quantity_available integer DEFAULT 0 CHECK (quantity_available >= 0),
  quantity_reserved integer DEFAULT 0 CHECK (quantity_reserved >= 0),
  quantity_sold integer DEFAULT 0 CHECK (quantity_sold >= 0),
  reorder_level integer DEFAULT 10,
  warehouse_id text,
  location_id text,
  last_counted_at timestamptz,
  low_stock_notification_sent_at timestamptz,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  CONSTRAINT variant_location_unique UNIQUE (variant_id, warehouse_id, location_id)
);

-- Storefront configurations
CREATE TABLE IF NOT EXISTS storefront_configs (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid NOT NULL UNIQUE,
  theme_colors jsonb DEFAULT '{
    "primary": "#000000",
    "secondary": "#FFFFFF",
    "accent": "#3B82F6"
  }'::jsonb,
  typography jsonb DEFAULT '{
    "heading_font": "sans-serif",
    "body_font": "sans-serif",
    "heading_size": "2rem"
  }'::jsonb,
  layout_settings jsonb DEFAULT '{
    "container_width": "1200px",
    "header_style": "standard"
  }'::jsonb,
  sections jsonb DEFAULT '[]'::jsonb,
  favicon_url text,
  logo_url text,
  banner_image_url text,
  social_links jsonb DEFAULT '{}'::jsonb,
  footer_text text,
  custom_css text,
  google_analytics_id text,
  facebook_pixel_id text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- Indexes for performance
CREATE INDEX idx_products_merchant_id ON products(merchant_id);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_visibility ON products(visibility);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_slug ON products(slug);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_published_at ON products(published_at DESC);
CREATE INDEX idx_products_tags ON products USING GIN(tags);
CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
CREATE INDEX idx_product_variants_merchant_id ON product_variants(merchant_id);
CREATE INDEX idx_product_variants_sku ON product_variants(sku);
CREATE INDEX idx_product_categories_merchant_id ON product_categories(merchant_id);
CREATE INDEX idx_product_categories_slug ON product_categories(slug);
CREATE INDEX idx_product_images_product_id ON product_images(product_id);
CREATE INDEX idx_product_images_featured ON product_images(is_featured);
CREATE INDEX idx_inventory_variant_id ON inventory(variant_id);
CREATE INDEX idx_inventory_merchant_id ON inventory(merchant_id);
CREATE INDEX idx_inventory_warehouse_id ON inventory(warehouse_id);
CREATE INDEX idx_storefront_configs_merchant_id ON storefront_configs(merchant_id);

-- Enable Row Level Security
ALTER TABLE products ENABLE ROW LEVEL SECURITY;
ALTER TABLE product_variants ENABLE ROW LEVEL SECURITY;
ALTER TABLE product_categories ENABLE ROW LEVEL SECURITY;
ALTER TABLE product_images ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory ENABLE ROW LEVEL SECURITY;
ALTER TABLE storefront_configs ENABLE ROW LEVEL SECURITY;

-- RLS Policies for products
CREATE POLICY "Merchants can read own products"
  ON products FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Public can read published products"
  ON products FOR SELECT
  TO anon
  USING (status = 'published' AND visibility = 'public');

CREATE POLICY "Merchants can insert own products"
  ON products FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update own products"
  ON products FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  )
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can delete own products"
  ON products FOR DELETE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

-- Similar policies for variants, categories, images, and inventory
CREATE POLICY "Merchants can read own variants"
  ON product_variants FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can manage own variants"
  ON product_variants FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update own variants"
  ON product_variants FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  )
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can delete own variants"
  ON product_variants FOR DELETE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

-- Policies for categories, images, and inventory follow same pattern
CREATE POLICY "Merchants can manage own categories"
  ON product_categories FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can insert categories"
  ON product_categories FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update categories"
  ON product_categories FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can delete categories"
  ON product_categories FOR DELETE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can read own images"
  ON product_images FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can manage images"
  ON product_images FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update images"
  ON product_images FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can delete images"
  ON product_images FOR DELETE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can read own inventory"
  ON inventory FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can manage inventory"
  ON inventory FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update inventory"
  ON inventory FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can manage storefront config"
  ON storefront_configs FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update own storefront"
  ON storefront_configs FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can insert storefront config"
  ON storefront_configs FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

-- Updated_at triggers
CREATE TRIGGER products_updated_at BEFORE UPDATE ON products
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER product_variants_updated_at BEFORE UPDATE ON product_variants
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER product_categories_updated_at BEFORE UPDATE ON product_categories
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER product_images_updated_at BEFORE UPDATE ON product_images
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER inventory_updated_at BEFORE UPDATE ON inventory
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER storefront_configs_updated_at BEFORE UPDATE ON storefront_configs
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
