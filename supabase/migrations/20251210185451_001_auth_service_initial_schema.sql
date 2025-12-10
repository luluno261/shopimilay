/*
  # Auth Service - Initial Schema

  1. New Tables
    - `users` - Platform users (merchants and admins)
    - `merchant_accounts` - Merchant account information
    - `api_keys` - API authentication tokens
    - `audit_logs` - Audit trail for security events

  2. Security
    - Enable RLS on all tables
    - Create policies for user data access
    - Create policies for merchant accounts
    - Create policies for API key management

  3. Implementation Notes
    - Uses JWT authentication from Supabase auth
    - Users table is separate from auth.users for custom merchant data
    - API keys support for programmatic access
    - Comprehensive audit logging for compliance
*/

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table (merchants and admins)
CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email text UNIQUE NOT NULL,
  full_name text,
  phone text,
  avatar_url text,
  role text DEFAULT 'merchant' CHECK (role IN ('merchant', 'admin', 'support')),
  status text DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
  metadata jsonb DEFAULT '{}'::jsonb,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  last_login_at timestamptz
);

-- Merchant accounts
CREATE TABLE IF NOT EXISTS merchant_accounts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  store_name text NOT NULL,
  store_slug text UNIQUE NOT NULL,
  description text,
  logo_url text,
  banner_url text,
  country text,
  currency text DEFAULT 'USD',
  timezone text DEFAULT 'UTC',
  stripe_account_id text,
  stripe_connected_at timestamptz,
  subscription_tier text DEFAULT 'basic' CHECK (subscription_tier IN ('basic', 'professional', 'enterprise')),
  subscription_status text DEFAULT 'active' CHECK (subscription_status IN ('active', 'trial', 'cancelled', 'suspended')),
  trial_ends_at timestamptz,
  settings jsonb DEFAULT '{}'::jsonb,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- API Keys for programmatic access
CREATE TABLE IF NOT EXISTS api_keys (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid NOT NULL REFERENCES merchant_accounts(id) ON DELETE CASCADE,
  name text NOT NULL,
  key_hash text NOT NULL UNIQUE,
  prefix text NOT NULL,
  last_used_at timestamptz,
  expires_at timestamptz,
  permissions text[] DEFAULT '{"read:products", "read:orders"}'::text[],
  status text DEFAULT 'active' CHECK (status IN ('active', 'revoked')),
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- Audit logs for compliance and security
CREATE TABLE IF NOT EXISTS audit_logs (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid REFERENCES users(id) ON DELETE SET NULL,
  merchant_id uuid REFERENCES merchant_accounts(id) ON DELETE SET NULL,
  action text NOT NULL,
  resource_type text NOT NULL,
  resource_id text,
  changes jsonb,
  ip_address text,
  user_agent text,
  status text DEFAULT 'success' CHECK (status IN ('success', 'failure')),
  error_message text,
  created_at timestamptz DEFAULT now()
);

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_merchant_accounts_user_id ON merchant_accounts(user_id);
CREATE INDEX idx_merchant_accounts_store_slug ON merchant_accounts(store_slug);
CREATE INDEX idx_merchant_accounts_stripe_id ON merchant_accounts(stripe_account_id);
CREATE INDEX idx_api_keys_merchant_id ON api_keys(merchant_id);
CREATE INDEX idx_api_keys_prefix ON api_keys(prefix);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_merchant_id ON audit_logs(merchant_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);

-- Enable Row Level Security
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE merchant_accounts ENABLE ROW LEVEL SECURITY;
ALTER TABLE api_keys ENABLE ROW LEVEL SECURITY;
ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;

-- Users policies
CREATE POLICY "Users can read own profile"
  ON users FOR SELECT
  TO authenticated
  USING (auth.uid()::text = id::text);

CREATE POLICY "Users can update own profile"
  ON users FOR UPDATE
  TO authenticated
  USING (auth.uid()::text = id::text)
  WITH CHECK (auth.uid()::text = id::text);

-- Merchant accounts policies
CREATE POLICY "Merchants can read own account"
  ON merchant_accounts FOR SELECT
  TO authenticated
  USING (
    user_id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
  );

CREATE POLICY "Merchants can update own account"
  ON merchant_accounts FOR UPDATE
  TO authenticated
  USING (
    user_id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
  )
  WITH CHECK (
    user_id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
  );

CREATE POLICY "Merchants can insert own account"
  ON merchant_accounts FOR INSERT
  TO authenticated
  WITH CHECK (
    user_id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
  );

-- API Keys policies
CREATE POLICY "Merchants can manage own API keys"
  ON api_keys FOR SELECT
  TO authenticated
  USING (
    merchant_id IN (
      SELECT id FROM merchant_accounts 
      WHERE user_id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
    )
  );

CREATE POLICY "Merchants can create API keys"
  ON api_keys FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id IN (
      SELECT id FROM merchant_accounts 
      WHERE user_id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
    )
  );

CREATE POLICY "Merchants can update own API keys"
  ON api_keys FOR UPDATE
  TO authenticated
  USING (
    merchant_id IN (
      SELECT id FROM merchant_accounts 
      WHERE user_id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
    )
  )
  WITH CHECK (
    merchant_id IN (
      SELECT id FROM merchant_accounts 
      WHERE user_id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
    )
  );

CREATE POLICY "Merchants can delete own API keys"
  ON api_keys FOR DELETE
  TO authenticated
  USING (
    merchant_id IN (
      SELECT id FROM merchant_accounts 
      WHERE user_id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
    )
  );

-- Audit logs policies
CREATE POLICY "Users can read own audit logs"
  ON audit_logs FOR SELECT
  TO authenticated
  USING (
    user_id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
    OR merchant_id IN (
      SELECT id FROM merchant_accounts 
      WHERE user_id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
    )
  );

CREATE POLICY "System can insert audit logs"
  ON audit_logs FOR INSERT
  WITH CHECK (true);

-- Updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for updated_at
CREATE TRIGGER users_updated_at BEFORE UPDATE ON users
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER merchant_accounts_updated_at BEFORE UPDATE ON merchant_accounts
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER api_keys_updated_at BEFORE UPDATE ON api_keys
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
