output "broker_endpoints" {
  value = aws_msk_cluster.main.bootstrap_brokers_tls
}

