output "vpc_id" {
  description = "ID du VPC"
  value       = module.vpc.vpc_id
}

output "rds_endpoint" {
  description = "Endpoint RDS"
  value       = module.rds.db_endpoint
  sensitive   = true
}

output "msk_broker_endpoints" {
  description = "Endpoints des brokers MSK"
  value       = module.msk.broker_endpoints
}

output "alb_dns_name" {
  description = "DNS name de l'ALB"
  value       = module.alb.alb_dns_name
}

output "ecs_cluster_name" {
  description = "Nom du cluster ECS"
  value       = module.ecs.cluster_name
}

output "elasticsearch_endpoint" {
  description = "Endpoint Elasticsearch"
  value       = module.elasticsearch.endpoint
}

