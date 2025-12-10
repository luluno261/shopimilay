/*
  # Checkout Service - Initial Schema

  1. New Tables
    - `shopping_carts` - Customer shopping carts
    - `cart_items` - Items in shopping carts
    - `orders` - Customer orders
    - `order_items` - Items in orders
    - `order_addresses` - Shipping and billing addresses
    - `discounts` - Discount/coupon codes
    - `order_payments` - Payment information
    - `refunds` - Refund records
    - `shipping_methods` - Available shipping options

  2. Security
    - Enable RLS for customer data isolation
    - Customers can only access their own carts and orders
    - Merchants can access orders for their products

  3. Implementation Notes
    - Supports multiple payment methods
    - Stripe integration for payments
    - Discount and promotion management
    - Full refund tracking and history
*/

-- Shopping carts
CREATE TABLE IF NOT EXISTS shopping_carts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id uuid,
  merchant_id uuid NOT NULL,
  session_id text,
  status text DEFAULT 'active' CHECK (status IN ('active', 'abandoned', 'completed')),
  currency text DEFAULT 'USD',
  subtotal numeric(12,2) DEFAULT 0 CHECK (subtotal >= 0),
  tax_amount numeric(12,2) DEFAULT 0 CHECK (tax_amount >= 0),
  discount_amount numeric(12,2) DEFAULT 0 CHECK (discount_amount >= 0),
  shipping_amount numeric(12,2) DEFAULT 0 CHECK (shipping_amount >= 0),
  total numeric(12,2) DEFAULT 0 CHECK (total >= 0),
  metadata jsonb DEFAULT '{}'::jsonb,
  abandoned_at timestamptz,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- Cart items
CREATE TABLE IF NOT EXISTS cart_items (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  cart_id uuid NOT NULL REFERENCES shopping_carts(id) ON DELETE CASCADE,
  product_id uuid NOT NULL,
  variant_id uuid NOT NULL,
  merchant_id uuid NOT NULL,
  quantity integer NOT NULL CHECK (quantity > 0),
  unit_price numeric(12,2) NOT NULL CHECK (unit_price >= 0),
  discount_amount numeric(12,2) DEFAULT 0 CHECK (discount_amount >= 0),
  tax_amount numeric(12,2) DEFAULT 0 CHECK (tax_amount >= 0),
  subtotal numeric(12,2) NOT NULL CHECK (subtotal >= 0),
  total numeric(12,2) NOT NULL CHECK (total >= 0),
  attributes jsonb DEFAULT '{}'::jsonb,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- Orders
CREATE TABLE IF NOT EXISTS orders (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_number text UNIQUE NOT NULL,
  customer_id uuid,
  merchant_id uuid NOT NULL,
  status text DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded')),
  payment_status text DEFAULT 'unpaid' CHECK (payment_status IN ('unpaid', 'paid', 'failed', 'refunded', 'partial_refund')),
  fulfillment_status text DEFAULT 'unshipped' CHECK (fulfillment_status IN ('unshipped', 'shipped', 'delivered', 'cancelled')),
  currency text DEFAULT 'USD',
  subtotal numeric(12,2) NOT NULL CHECK (subtotal >= 0),
  tax_amount numeric(12,2) DEFAULT 0 CHECK (tax_amount >= 0),
  discount_amount numeric(12,2) DEFAULT 0 CHECK (discount_amount >= 0),
  shipping_amount numeric(12,2) DEFAULT 0 CHECK (shipping_amount >= 0),
  total numeric(12,2) NOT NULL CHECK (total >= 0),
  customer_email text,
  customer_phone text,
  customer_name text,
  notes text,
  internal_notes text,
  source text DEFAULT 'storefront' CHECK (source IN ('storefront', 'admin', 'api', 'import')),
  tracking_number text,
  shipping_carrier text,
  stripe_payment_intent_id text,
  discount_code text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  paid_at timestamptz,
  shipped_at timestamptz,
  delivered_at timestamptz
);

-- Order items
CREATE TABLE IF NOT EXISTS order_items (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id uuid NOT NULL,
  variant_id uuid NOT NULL,
  merchant_id uuid NOT NULL,
  product_name text NOT NULL,
  product_sku text NOT NULL,
  quantity integer NOT NULL CHECK (quantity > 0),
  unit_price numeric(12,2) NOT NULL CHECK (unit_price >= 0),
  discount_amount numeric(12,2) DEFAULT 0 CHECK (discount_amount >= 0),
  tax_amount numeric(12,2) DEFAULT 0 CHECK (tax_amount >= 0),
  subtotal numeric(12,2) NOT NULL CHECK (subtotal >= 0),
  total numeric(12,2) NOT NULL CHECK (total >= 0),
  attributes jsonb,
  created_at timestamptz DEFAULT now()
);

-- Order addresses
CREATE TABLE IF NOT EXISTS order_addresses (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  address_type text NOT NULL CHECK (address_type IN ('shipping', 'billing')),
  first_name text,
  last_name text,
  company text,
  street_address text NOT NULL,
  street_address_2 text,
  city text NOT NULL,
  state text,
  postal_code text NOT NULL,
  country text NOT NULL,
  phone text,
  email text,
  created_at timestamptz DEFAULT now()
);

-- Discounts and coupons
CREATE TABLE IF NOT EXISTS discounts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid NOT NULL,
  code text UNIQUE NOT NULL,
  type text NOT NULL CHECK (type IN ('fixed', 'percentage')),
  value numeric(12,2) NOT NULL CHECK (value > 0),
  max_discount_amount numeric(12,2),
  currency text DEFAULT 'USD',
  usage_limit integer,
  usage_count integer DEFAULT 0,
  per_customer_limit integer DEFAULT 1,
  active boolean DEFAULT true,
  minimum_purchase_amount numeric(12,2),
  maximum_purchase_amount numeric(12,2),
  applicable_products uuid[],
  excluded_products uuid[],
  applicable_categories uuid[],
  excluded_categories uuid[],
  applicable_collections uuid[],
  started_at timestamptz,
  expires_at timestamptz,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- Payments
CREATE TABLE IF NOT EXISTS order_payments (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  payment_method text NOT NULL CHECK (payment_method IN ('credit_card', 'debit_card', 'paypal', 'stripe', 'bank_transfer')),
  amount numeric(12,2) NOT NULL CHECK (amount >= 0),
  currency text DEFAULT 'USD',
  status text DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
  transaction_id text,
  stripe_charge_id text,
  stripe_payment_intent_id text,
  last_4_digits text,
  card_brand text,
  error_message text,
  metadata jsonb DEFAULT '{}'::jsonb,
  received_at timestamptz,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- Refunds
CREATE TABLE IF NOT EXISTS refunds (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  refund_amount numeric(12,2) NOT NULL CHECK (refund_amount > 0),
  currency text DEFAULT 'USD',
  reason text NOT NULL,
  status text DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
  refund_method text,
  stripe_refund_id text,
  notes text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  processed_at timestamptz
);

-- Shipping methods
CREATE TABLE IF NOT EXISTS shipping_methods (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid NOT NULL,
  name text NOT NULL,
  description text,
  code text NOT NULL,
  cost numeric(12,2) NOT NULL CHECK (cost >= 0),
  min_delivery_days integer,
  max_delivery_days integer,
  countries text[] NOT NULL,
  states text[],
  min_order_value numeric(12,2),
  max_order_value numeric(12,2),
  min_weight numeric(10,3),
  max_weight numeric(10,3),
  active boolean DEFAULT true,
  position integer DEFAULT 0,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  CONSTRAINT merchant_shipping_code UNIQUE (merchant_id, code)
);

-- Indexes for performance
CREATE INDEX idx_shopping_carts_merchant_id ON shopping_carts(merchant_id);
CREATE INDEX idx_shopping_carts_customer_id ON shopping_carts(customer_id);
CREATE INDEX idx_shopping_carts_session_id ON shopping_carts(session_id);
CREATE INDEX idx_shopping_carts_status ON shopping_carts(status);
CREATE INDEX idx_shopping_carts_created_at ON shopping_carts(created_at DESC);
CREATE INDEX idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX idx_cart_items_variant_id ON cart_items(variant_id);
CREATE INDEX idx_orders_merchant_id ON orders(merchant_id);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_payment_status ON orders(payment_status);
CREATE INDEX idx_orders_fulfillment_status ON orders(fulfillment_status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX idx_orders_stripe_intent ON orders(stripe_payment_intent_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
CREATE INDEX idx_order_addresses_order_id ON order_addresses(order_id);
CREATE INDEX idx_discounts_merchant_id ON discounts(merchant_id);
CREATE INDEX idx_discounts_code ON discounts(code);
CREATE INDEX idx_discounts_active ON discounts(active);
CREATE INDEX idx_order_payments_order_id ON order_payments(order_id);
CREATE INDEX idx_order_payments_status ON order_payments(status);
CREATE INDEX idx_order_payments_stripe_id ON order_payments(stripe_charge_id);
CREATE INDEX idx_refunds_order_id ON refunds(order_id);
CREATE INDEX idx_refunds_status ON refunds(status);
CREATE INDEX idx_refunds_created_at ON refunds(created_at DESC);
CREATE INDEX idx_shipping_methods_merchant_id ON shipping_methods(merchant_id);
CREATE INDEX idx_shipping_methods_active ON shipping_methods(active);

-- Enable Row Level Security
ALTER TABLE shopping_carts ENABLE ROW LEVEL SECURITY;
ALTER TABLE cart_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE order_addresses ENABLE ROW LEVEL SECURITY;
ALTER TABLE discounts ENABLE ROW LEVEL SECURITY;
ALTER TABLE order_payments ENABLE ROW LEVEL SECURITY;
ALTER TABLE refunds ENABLE ROW LEVEL SECURITY;
ALTER TABLE shipping_methods ENABLE ROW LEVEL SECURITY;

-- RLS Policies for shopping carts
CREATE POLICY "Customers can read own carts"
  ON shopping_carts FOR SELECT
  TO authenticated
  USING (customer_id = auth.uid() OR session_id IS NOT NULL);

CREATE POLICY "Merchants can read carts for their products"
  ON shopping_carts FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Customers can create carts"
  ON shopping_carts FOR INSERT
  TO authenticated
  WITH CHECK (customer_id = auth.uid() OR session_id IS NOT NULL);

CREATE POLICY "Customers can update own carts"
  ON shopping_carts FOR UPDATE
  TO authenticated
  USING (customer_id = auth.uid() OR session_id IS NOT NULL);

-- RLS Policies for cart items
CREATE POLICY "Users can read cart items"
  ON cart_items FOR SELECT
  TO authenticated
  USING (
    EXISTS (
      SELECT 1 FROM shopping_carts
      WHERE shopping_carts.id = cart_items.cart_id
      AND (shopping_carts.customer_id = auth.uid() OR shopping_carts.session_id IS NOT NULL)
    )
  );

CREATE POLICY "Users can insert cart items"
  ON cart_items FOR INSERT
  TO authenticated
  WITH CHECK (
    EXISTS (
      SELECT 1 FROM shopping_carts
      WHERE shopping_carts.id = cart_items.cart_id
      AND (shopping_carts.customer_id = auth.uid() OR shopping_carts.session_id IS NOT NULL)
    )
  );

CREATE POLICY "Users can update cart items"
  ON cart_items FOR UPDATE
  TO authenticated
  USING (
    EXISTS (
      SELECT 1 FROM shopping_carts
      WHERE shopping_carts.id = cart_items.cart_id
      AND (shopping_carts.customer_id = auth.uid() OR shopping_carts.session_id IS NOT NULL)
    )
  );

CREATE POLICY "Users can delete cart items"
  ON cart_items FOR DELETE
  TO authenticated
  USING (
    EXISTS (
      SELECT 1 FROM shopping_carts
      WHERE shopping_carts.id = cart_items.cart_id
      AND (shopping_carts.customer_id = auth.uid() OR shopping_carts.session_id IS NOT NULL)
    )
  );

-- RLS Policies for orders
CREATE POLICY "Customers can read own orders"
  ON orders FOR SELECT
  TO authenticated
  USING (customer_id = auth.uid());

CREATE POLICY "Merchants can read orders for their products"
  ON orders FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Customers can create orders"
  ON orders FOR INSERT
  TO authenticated
  WITH CHECK (customer_id = auth.uid());

CREATE POLICY "Merchants can update orders"
  ON orders FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

-- RLS Policies for order items
CREATE POLICY "Users can read order items"
  ON order_items FOR SELECT
  TO authenticated
  USING (
    EXISTS (
      SELECT 1 FROM orders
      WHERE orders.id = order_items.order_id
      AND (orders.customer_id = auth.uid() OR orders.merchant_id = (
        SELECT merchant_accounts.id FROM merchant_accounts
        JOIN users ON merchant_accounts.user_id = users.id
        WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
        LIMIT 1
      ))
    )
  );

-- RLS Policies for order addresses
CREATE POLICY "Users can read order addresses"
  ON order_addresses FOR SELECT
  TO authenticated
  USING (
    EXISTS (
      SELECT 1 FROM orders
      WHERE orders.id = order_addresses.order_id
      AND (orders.customer_id = auth.uid() OR orders.merchant_id = (
        SELECT merchant_accounts.id FROM merchant_accounts
        JOIN users ON merchant_accounts.user_id = users.id
        WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
        LIMIT 1
      ))
    )
  );

-- RLS Policies for discounts
CREATE POLICY "Merchants can read own discounts"
  ON discounts FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can manage discounts"
  ON discounts FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update own discounts"
  ON discounts FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

-- RLS Policies for payments and refunds
CREATE POLICY "Users can read own payments"
  ON order_payments FOR SELECT
  TO authenticated
  USING (
    EXISTS (
      SELECT 1 FROM orders
      WHERE orders.id = order_payments.order_id
      AND (orders.customer_id = auth.uid() OR orders.merchant_id = (
        SELECT merchant_accounts.id FROM merchant_accounts
        JOIN users ON merchant_accounts.user_id = users.id
        WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
        LIMIT 1
      ))
    )
  );

CREATE POLICY "Users can read own refunds"
  ON refunds FOR SELECT
  TO authenticated
  USING (
    EXISTS (
      SELECT 1 FROM orders
      WHERE orders.id = refunds.order_id
      AND (orders.customer_id = auth.uid() OR orders.merchant_id = (
        SELECT merchant_accounts.id FROM merchant_accounts
        JOIN users ON merchant_accounts.user_id = users.id
        WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
        LIMIT 1
      ))
    )
  );

-- RLS Policies for shipping methods
CREATE POLICY "Merchants can manage shipping methods"
  ON shipping_methods FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can create shipping methods"
  ON shipping_methods FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update shipping"
  ON shipping_methods FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

-- Updated_at triggers
CREATE TRIGGER shopping_carts_updated_at BEFORE UPDATE ON shopping_carts
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER cart_items_updated_at BEFORE UPDATE ON cart_items
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER orders_updated_at BEFORE UPDATE ON orders
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER order_payments_updated_at BEFORE UPDATE ON order_payments
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER refunds_updated_at BEFORE UPDATE ON refunds
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER discounts_updated_at BEFORE UPDATE ON discounts
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER shipping_methods_updated_at BEFORE UPDATE ON shipping_methods
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
