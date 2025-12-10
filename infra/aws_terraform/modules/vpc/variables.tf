variable "vpc_cidr" {
  description = "CIDR block pour le VPC"
  type        = string
}

variable "availability_zones" {
  description = "Zones de disponibilité"
  type        = list(string)
}

variable "environment" {
  description = "Environnement"
  type        = string
}

variable "project_name" {
  description = "Nom du projet"
  type        = string
}

