# Infrastructure Terraform - OmniSphere

Infrastructure as Code pour déployer OmniSphere sur AWS.

## Architecture

- **VPC** : Réseau virtuel avec sous-réseaux publics et privés
- **RDS** : PostgreSQL pour les bases de données
- **MSK** : Amazon Managed Streaming for Apache Kafka
- **ECS** : Cluster ECS pour les microservices
- **ALB** : Application Load Balancer
- **Elasticsearch** : Pour la recherche de produits

## Utilisation

```bash
# Initialiser Terraform
terraform init

# Planifier les changements
terraform plan

# Appliquer les changements
terraform apply

# Détruire l'infrastructure
terraform destroy
```

## Configuration

1. Copier `terraform.tfvars.example` vers `terraform.tfvars`
2. Remplir les variables nécessaires
3. Configurer le backend S3 pour le state Terraform

## Modules

- `modules/vpc` - VPC et sous-réseaux
- `modules/rds` - Base de données PostgreSQL
- `modules/msk` - Cluster Kafka
- `modules/ecs` - Cluster ECS
- `modules/alb` - Application Load Balancer
- `modules/elasticsearch` - Cluster Elasticsearch

