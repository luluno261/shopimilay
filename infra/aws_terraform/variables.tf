variable "aws_region" {
  description = "Région AWS"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environnement (dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "project_name" {
  description = "Nom du projet"
  type        = string
  default     = "omnisphere"
}

variable "vpc_cidr" {
  description = "CIDR block pour le VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "Zones de disponibilité"
  type        = list(string)
  default     = ["us-east-1a", "us-east-1b"]
}

# RDS Variables
variable "db_instance_class" {
  description = "Classe d'instance RDS"
  type        = string
  default     = "db.t3.medium"
}

variable "db_allocated_storage" {
  description = "Stockage alloué pour RDS (GB)"
  type        = number
  default     = 100
}

variable "db_name" {
  description = "Nom de la base de données"
  type        = string
  default     = "omnisphere"
}

variable "db_username" {
  description = "Nom d'utilisateur de la base de données"
  type        = string
  sensitive   = true
}

variable "db_password" {
  description = "Mot de passe de la base de données"
  type        = string
  sensitive   = true
}

# MSK Variables
variable "kafka_version" {
  description = "Version de Kafka"
  type        = string
  default     = "3.5.1"
}

variable "msk_instance_type" {
  description = "Type d'instance MSK"
  type        = string
  default     = "kafka.m5.large"
}

variable "msk_broker_count" {
  description = "Nombre de brokers MSK"
  type        = number
  default     = 3
}

# Elasticsearch Variables
variable "elasticsearch_instance_type" {
  description = "Type d'instance Elasticsearch"
  type        = string
  default     = "t3.small.elasticsearch"
}

variable "elasticsearch_instance_count" {
  description = "Nombre d'instances Elasticsearch"
  type        = number
  default     = 2
}

