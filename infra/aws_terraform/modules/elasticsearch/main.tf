resource "aws_elasticsearch_domain" "main" {
  domain_name           = "${var.project_name}-${var.environment}-es"
  elasticsearch_version = "7.10"

  cluster_config {
    instance_type  = var.instance_type
    instance_count = var.instance_count
  }

  vpc_options {
    subnet_ids         = var.private_subnet_ids
    security_group_ids = [aws_security_group.es.id]
  }

  ebs_options {
    ebs_enabled = true
    volume_size = 20
  }

  encrypt_at_rest {
    enabled = true
  }

  node_to_node_encryption {
    enabled = true
  }

  domain_endpoint_options {
    enforce_https       = true
    tls_security_policy = "Policy-Min-TLS-1-2-2019-07"
  }

  tags = {
    Name = "${var.project_name}-${var.environment}-es"
  }
}

resource "aws_security_group" "es" {
  name        = "${var.project_name}-${var.environment}-es-sg"
  description = "Security group pour Elasticsearch"
  vpc_id      = var.vpc_id

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project_name}-${var.environment}-es-sg"
  }
}

