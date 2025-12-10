variable "vpc_id" {
  type = string
}

variable "vpc_cidr" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
}

variable "kafka_version" {
  type = string
}

variable "instance_type" {
  type = string
}

variable "broker_count" {
  type = number
}

variable "environment" {
  type = string
}

variable "project_name" {
  type = string
}

