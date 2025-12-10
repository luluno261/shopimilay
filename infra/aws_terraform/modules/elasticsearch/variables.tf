variable "vpc_id" {
  type = string
}

variable "vpc_cidr" {
  type = string
}

variable "private_subnet_ids" {
  type = list(string)
}

variable "instance_type" {
  type = string
}

variable "instance_count" {
  type = number
}

variable "environment" {
  type = string
}

variable "project_name" {
  type = string
}

