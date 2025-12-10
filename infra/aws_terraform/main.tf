terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  
  backend "s3" {
    # Configuration du backend S3 pour le state Terraform
    # À configurer selon votre environnement
    bucket = "omnisphere-terraform-state"
    key    = "omnisphere/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = var.aws_region
}

# VPC
module "vpc" {
  source = "./modules/vpc"
  
  vpc_cidr             = var.vpc_cidr
  availability_zones   = var.availability_zones
  environment          = var.environment
  project_name         = var.project_name
}

# RDS PostgreSQL
module "rds" {
  source = "./modules/rds"
  
  vpc_id              = module.vpc.vpc_id
  vpc_cidr            = module.vpc.vpc_cidr
  private_subnet_ids  = module.vpc.private_subnet_ids
  db_instance_class   = var.db_instance_class
  db_allocated_storage = var.db_allocated_storage
  db_name             = var.db_name
  db_username         = var.db_username
  db_password         = var.db_password
  environment         = var.environment
  project_name         = var.project_name
}

# MSK (Kafka)
module "msk" {
  source = "./modules/msk"
  
  vpc_id            = module.vpc.vpc_id
  vpc_cidr          = module.vpc.vpc_cidr
  subnet_ids        = module.vpc.private_subnet_ids
  kafka_version     = var.kafka_version
  instance_type     = var.msk_instance_type
  broker_count      = var.msk_broker_count
  environment       = var.environment
  project_name       = var.project_name
}

# ECS Cluster
module "ecs" {
  source = "./modules/ecs"
  
  vpc_id            = module.vpc.vpc_id
  public_subnet_ids = module.vpc.public_subnet_ids
  private_subnet_ids = module.vpc.private_subnet_ids
  environment       = var.environment
  project_name       = var.project_name
}

# Application Load Balancer
module "alb" {
  source = "./modules/alb"
  
  vpc_id            = module.vpc.vpc_id
  public_subnet_ids = module.vpc.public_subnet_ids
  environment       = var.environment
  project_name       = var.project_name
}

# Elasticsearch (optionnel, peut utiliser AWS OpenSearch)
module "elasticsearch" {
  source = "./modules/elasticsearch"
  
  vpc_id            = module.vpc.vpc_id
  vpc_cidr          = module.vpc.vpc_cidr
  private_subnet_ids = module.vpc.private_subnet_ids
  instance_type     = var.elasticsearch_instance_type
  instance_count    = var.elasticsearch_instance_count
  environment       = var.environment
  project_name       = var.project_name
}

