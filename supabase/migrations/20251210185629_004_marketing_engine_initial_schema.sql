/*
  # Marketing Engine - Initial Schema

  1. New Tables
    - `customers` - Customer profiles with attributes
    - `customer_events` - Event tracking for CDP
    - `customer_segments` - Dynamic audience segments
    - `segment_members` - Customers in segments
    - `email_campaigns` - Email campaign definitions
    - `email_templates` - Reusable email templates
    - `automation_workflows` - Automation workflow definitions
    - `workflow_steps` - Steps in workflows
    - `email_sends` - Email send history
    - `email_opens` - Email engagement tracking
    - `email_clicks` - Link click tracking
    - `capture_popups` - Email capture popup definitions

  2. Security
    - Enable RLS for customer data isolation
    - Merchants can only access their customers and campaigns

  3. Implementation Notes
    - Full CDP with event tracking
    - Email marketing capabilities
    - Automation workflows for lifecycle marketing
    - Conversion optimization with capture popups
*/

-- Customers
CREATE TABLE IF NOT EXISTS customers (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid NOT NULL,
  email text NOT NULL,
  phone text,
  first_name text,
  last_name text,
  avatar_url text,
  status text DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'unsubscribed', 'bounced')),
  email_subscription_status text DEFAULT 'unsubscribed' CHECK (email_subscription_status IN ('subscribed', 'unsubscribed', 'pending', 'bounced')),
  sms_subscription_status text DEFAULT 'unsubscribed' CHECK (sms_subscription_status IN ('subscribed', 'unsubscribed', 'pending')),
  lifetime_value numeric(12,2) DEFAULT 0 CHECK (lifetime_value >= 0),
  total_orders integer DEFAULT 0,
  total_spent numeric(12,2) DEFAULT 0 CHECK (total_spent >= 0),
  average_order_value numeric(12,2) DEFAULT 0,
  last_order_date timestamptz,
  first_order_date timestamptz,
  tags text[] DEFAULT '{}'::text[],
  metadata jsonb DEFAULT '{}'::jsonb,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  CONSTRAINT merchant_customer_email UNIQUE (merchant_id, email)
);

-- Customer events (CDP data)
CREATE TABLE IF NOT EXISTS customer_events (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id uuid REFERENCES customers(id) ON DELETE CASCADE,
  merchant_id uuid NOT NULL,
  event_type text NOT NULL,
  event_data jsonb DEFAULT '{}'::jsonb,
  source text DEFAULT 'system' CHECK (source IN ('system', 'storefront', 'popup', 'form', 'api')),
  user_agent text,
  ip_address text,
  session_id text,
  created_at timestamptz DEFAULT now()
);

-- Customer segments for targeting
CREATE TABLE IF NOT EXISTS customer_segments (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid NOT NULL,
  name text NOT NULL,
  description text,
  criteria jsonb NOT NULL,
  member_count integer DEFAULT 0,
  status text DEFAULT 'active' CHECK (status IN ('active', 'paused', 'archived')),
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  CONSTRAINT merchant_segment_name UNIQUE (merchant_id, name)
);

-- Members of segments (denormalized for performance)
CREATE TABLE IF NOT EXISTS segment_members (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  segment_id uuid NOT NULL REFERENCES customer_segments(id) ON DELETE CASCADE,
  customer_id uuid NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
  added_at timestamptz DEFAULT now(),
  CONSTRAINT segment_member_unique UNIQUE (segment_id, customer_id)
);

-- Email campaigns
CREATE TABLE IF NOT EXISTS email_campaigns (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid NOT NULL,
  name text NOT NULL,
  subject_line text NOT NULL,
  preview_text text,
  template_id uuid,
  segment_id uuid,
  from_name text,
  from_email text,
  reply_to_email text,
  status text DEFAULT 'draft' CHECK (status IN ('draft', 'scheduled', 'sending', 'sent', 'paused', 'cancelled')),
  scheduled_at timestamptz,
  sent_at timestamptz,
  total_recipients integer DEFAULT 0,
  total_sent integer DEFAULT 0,
  total_opened integer DEFAULT 0,
  total_clicked integer DEFAULT 0,
  total_bounced integer DEFAULT 0,
  total_unsubscribed integer DEFAULT 0,
  open_rate numeric(5,2) DEFAULT 0,
  click_rate numeric(5,2) DEFAULT 0,
  conversion_rate numeric(5,2) DEFAULT 0,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- Email templates
CREATE TABLE IF NOT EXISTS email_templates (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid NOT NULL,
  name text NOT NULL,
  subject_line text NOT NULL,
  preheader text,
  html_content text NOT NULL,
  text_content text,
  variables text[] DEFAULT '{"first_name", "email", "order_id"}'::text[],
  category text,
  status text DEFAULT 'active' CHECK (status IN ('active', 'archived')),
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  CONSTRAINT merchant_template_name UNIQUE (merchant_id, name)
);

-- Automation workflows
CREATE TABLE IF NOT EXISTS automation_workflows (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid NOT NULL,
  name text NOT NULL,
  description text,
  trigger_type text NOT NULL CHECK (trigger_type IN ('customer_signup', 'order_placed', 'abandoned_cart', 'customer_email_subscribed', 'tag_added', 'custom')),
  trigger_data jsonb,
  status text DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'paused', 'archived')),
  active_member_count integer DEFAULT 0,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  CONSTRAINT merchant_workflow_name UNIQUE (merchant_id, name)
);

-- Workflow steps
CREATE TABLE IF NOT EXISTS workflow_steps (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  workflow_id uuid NOT NULL REFERENCES automation_workflows(id) ON DELETE CASCADE,
  position integer NOT NULL,
  step_type text NOT NULL CHECK (step_type IN ('delay', 'email', 'condition', 'tag', 'filter')),
  step_data jsonb NOT NULL,
  next_step_id uuid REFERENCES workflow_steps(id) ON DELETE SET NULL,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

-- Email send history
CREATE TABLE IF NOT EXISTS email_sends (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  campaign_id uuid REFERENCES email_campaigns(id) ON DELETE SET NULL,
  workflow_step_id uuid REFERENCES workflow_steps(id) ON DELETE SET NULL,
  customer_id uuid NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
  merchant_id uuid NOT NULL,
  recipient_email text NOT NULL,
  subject_line text,
  status text DEFAULT 'sent' CHECK (status IN ('sent', 'bounced', 'complained', 'rejected', 'dropped')),
  bounce_type text CHECK (bounce_type IN ('permanent', 'temporary', 'unknown')),
  error_message text,
  sent_at timestamptz DEFAULT now(),
  opened_at timestamptz,
  opened_count integer DEFAULT 0,
  first_click_at timestamptz,
  last_click_at timestamptz,
  click_count integer DEFAULT 0,
  unsubscribed_at timestamptz,
  complained_at timestamptz
);

-- Email opens (detailed tracking)
CREATE TABLE IF NOT EXISTS email_opens (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email_send_id uuid NOT NULL REFERENCES email_sends(id) ON DELETE CASCADE,
  opened_at timestamptz DEFAULT now(),
  user_agent text,
  ip_address text,
  location jsonb
);

-- Email clicks (detailed tracking)
CREATE TABLE IF NOT EXISTS email_clicks (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email_send_id uuid NOT NULL REFERENCES email_sends(id) ON DELETE CASCADE,
  link_url text NOT NULL,
  link_text text,
  clicked_at timestamptz DEFAULT now(),
  user_agent text,
  ip_address text,
  location jsonb
);

-- Capture popups for email collection
CREATE TABLE IF NOT EXISTS capture_popups (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  merchant_id uuid NOT NULL,
  name text NOT NULL,
  type text DEFAULT 'popup' CHECK (type IN ('popup', 'banner', 'slide_in', 'form')),
  title text,
  description text,
  cta_text text,
  background_color text DEFAULT '#FFFFFF',
  text_color text DEFAULT '#000000',
  button_color text DEFAULT '#3B82F6',
  button_text_color text DEFAULT '#FFFFFF',
  image_url text,
  targeting_rules jsonb,
  redirect_url text,
  offer_discount_code text,
  thank_you_message text,
  status text DEFAULT 'active' CHECK (status IN ('active', 'paused', 'archived')),
  display_frequency text DEFAULT 'once' CHECK (display_frequency IN ('once', 'always', 'daily', 'weekly')),
  impressions integer DEFAULT 0,
  conversions integer DEFAULT 0,
  conversion_rate numeric(5,2) DEFAULT 0,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now(),
  CONSTRAINT merchant_popup_name UNIQUE (merchant_id, name)
);

-- Indexes for performance
CREATE INDEX idx_customers_merchant_id ON customers(merchant_id);
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_status ON customers(status);
CREATE INDEX idx_customers_email_subscription ON customers(email_subscription_status);
CREATE INDEX idx_customers_created_at ON customers(created_at DESC);
CREATE INDEX idx_customer_events_customer_id ON customer_events(customer_id);
CREATE INDEX idx_customer_events_merchant_id ON customer_events(merchant_id);
CREATE INDEX idx_customer_events_type ON customer_events(event_type);
CREATE INDEX idx_customer_events_created_at ON customer_events(created_at DESC);
CREATE INDEX idx_customer_events_source ON customer_events(source);
CREATE INDEX idx_customer_events_session_id ON customer_events(session_id);
CREATE INDEX idx_customer_segments_merchant_id ON customer_segments(merchant_id);
CREATE INDEX idx_customer_segments_status ON customer_segments(status);
CREATE INDEX idx_segment_members_segment_id ON segment_members(segment_id);
CREATE INDEX idx_segment_members_customer_id ON segment_members(customer_id);
CREATE INDEX idx_email_campaigns_merchant_id ON email_campaigns(merchant_id);
CREATE INDEX idx_email_campaigns_status ON email_campaigns(status);
CREATE INDEX idx_email_campaigns_created_at ON email_campaigns(created_at DESC);
CREATE INDEX idx_email_templates_merchant_id ON email_templates(merchant_id);
CREATE INDEX idx_email_templates_status ON email_templates(status);
CREATE INDEX idx_automation_workflows_merchant_id ON automation_workflows(merchant_id);
CREATE INDEX idx_automation_workflows_status ON automation_workflows(status);
CREATE INDEX idx_automation_workflows_trigger_type ON automation_workflows(trigger_type);
CREATE INDEX idx_workflow_steps_workflow_id ON workflow_steps(workflow_id);
CREATE INDEX idx_email_sends_campaign_id ON email_sends(campaign_id);
CREATE INDEX idx_email_sends_customer_id ON email_sends(customer_id);
CREATE INDEX idx_email_sends_merchant_id ON email_sends(merchant_id);
CREATE INDEX idx_email_sends_status ON email_sends(status);
CREATE INDEX idx_email_sends_recipient_email ON email_sends(recipient_email);
CREATE INDEX idx_email_sends_sent_at ON email_sends(sent_at DESC);
CREATE INDEX idx_email_opens_email_send_id ON email_opens(email_send_id);
CREATE INDEX idx_email_clicks_email_send_id ON email_clicks(email_send_id);
CREATE INDEX idx_capture_popups_merchant_id ON capture_popups(merchant_id);
CREATE INDEX idx_capture_popups_status ON capture_popups(status);

-- Enable Row Level Security
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;
ALTER TABLE customer_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE customer_segments ENABLE ROW LEVEL SECURITY;
ALTER TABLE segment_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE email_campaigns ENABLE ROW LEVEL SECURITY;
ALTER TABLE email_templates ENABLE ROW LEVEL SECURITY;
ALTER TABLE automation_workflows ENABLE ROW LEVEL SECURITY;
ALTER TABLE workflow_steps ENABLE ROW LEVEL SECURITY;
ALTER TABLE email_sends ENABLE ROW LEVEL SECURITY;
ALTER TABLE email_opens ENABLE ROW LEVEL SECURITY;
ALTER TABLE email_clicks ENABLE ROW LEVEL SECURITY;
ALTER TABLE capture_popups ENABLE ROW LEVEL SECURITY;

-- RLS Policies for customers
CREATE POLICY "Merchants can read own customers"
  ON customers FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can manage customers"
  ON customers FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update customers"
  ON customers FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

-- RLS Policies for customer events
CREATE POLICY "Merchants can read customer events"
  ON customer_events FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "System can insert events"
  ON customer_events FOR INSERT
  WITH CHECK (true);

-- RLS Policies for segments
CREATE POLICY "Merchants can manage segments"
  ON customer_segments FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can create segments"
  ON customer_segments FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update segments"
  ON customer_segments FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

-- RLS Policies for campaigns, templates, and workflows (similar pattern)
CREATE POLICY "Merchants can manage email campaigns"
  ON email_campaigns FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can create campaigns"
  ON email_campaigns FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update campaigns"
  ON email_campaigns FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can manage email templates"
  ON email_templates FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can create templates"
  ON email_templates FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update templates"
  ON email_templates FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can manage workflows"
  ON automation_workflows FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can create workflows"
  ON automation_workflows FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update workflows"
  ON automation_workflows FOR UPDATE
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can manage popups"
  ON capture_popups FOR SELECT
  TO authenticated
  USING (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can create popups"
  ON capture_popups FOR INSERT
  TO authenticated
  WITH CHECK (
    merchant_id = (
      SELECT merchant_accounts.id FROM merchant_accounts
      JOIN users ON merchant_accounts.user_id = users.id
      WHERE users.id = (SELECT id FROM users WHERE auth.uid()::text = id::text LIMIT 1)
      LIMIT 1
    )
  );

CREATE POLICY "Merchants can update popups"
  ON capture_popups FOR UPDATE
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
CREATE TRIGGER customers_updated_at BEFORE UPDATE ON customers
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER customer_segments_updated_at BEFORE UPDATE ON customer_segments
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER email_campaigns_updated_at BEFORE UPDATE ON email_campaigns
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER email_templates_updated_at BEFORE UPDATE ON email_templates
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER automation_workflows_updated_at BEFORE UPDATE ON automation_workflows
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER workflow_steps_updated_at BEFORE UPDATE ON workflow_steps
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER capture_popups_updated_at BEFORE UPDATE ON capture_popups
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
