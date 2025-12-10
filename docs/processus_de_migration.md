# Processus de migration

Guide pour migrer les données et déployer les mises à jour de la plateforme OmniSphere.

## Types de migrations

### 1. Migrations de schéma de base de données

Les migrations de schéma sont gérées via des fichiers SQL dans `shared/database_migrations/`.

#### Exécution des migrations

```bash
# Migrations vers le haut (up)
./infra/devops_scripts/run-migrations.sh auth-service up

# Migrations vers le bas (rollback)
./infra/devops_scripts/run-migrations.sh auth-service down

# Toutes les migrations
./infra/devops_scripts/run-migrations.sh all up
```

#### Bonnes pratiques

1. **Nommage** : Utiliser un format séquentiel (001_, 002_, etc.)
2. **Réversibilité** : Toujours créer un fichier `.down.sql` pour le rollback
3. **Tests** : Tester les migrations en environnement de développement d'abord
4. **Backup** : Faire un backup de la base de données avant les migrations en production

### 2. Migrations de données

Pour les migrations de données volumineuses :

1. Créer un script de migration dans `infra/devops_scripts/`
2. Tester sur une copie de la base de données de production
3. Exécuter pendant une fenêtre de maintenance
4. Vérifier l'intégrité des données après migration

### 3. Migrations de code (déploiement)

#### Déploiement d'un service

```bash
# Déployer un service spécifique
./infra/devops_scripts/deploy-prod.sh auth-service

# Déployer tous les services
./infra/devops_scripts/deploy-prod.sh all
```

#### Stratégie de déploiement

1. **Blue-Green Deployment** : Déployer la nouvelle version parallèlement
2. **Canary Deployment** : Déployer progressivement à un pourcentage d'utilisateurs
3. **Rolling Update** : Mise à jour progressive (ECS par défaut)

## Checklist de migration

### Avant la migration

- [ ] Backup de la base de données
- [ ] Tests des migrations en environnement de staging
- [ ] Documentation des changements
- [ ] Plan de rollback préparé
- [ ] Fenêtre de maintenance planifiée

### Pendant la migration

- [ ] Exécution des migrations de schéma
- [ ] Vérification de l'intégrité des données
- [ ] Déploiement des nouveaux services
- [ ] Tests de smoke sur les endpoints critiques

### Après la migration

- [ ] Monitoring des métriques (erreurs, latence)
- [ ] Vérification des logs
- [ ] Tests de régression
- [ ] Communication aux utilisateurs si nécessaire

## Rollback

En cas de problème :

1. **Rollback des migrations de schéma** :
   ```bash
   ./infra/devops_scripts/run-migrations.sh [service] down
   ```

2. **Rollback du déploiement** :
   - Utiliser la version précédente de l'image Docker
   - Forcer un nouveau déploiement ECS avec l'ancienne version

3. **Restaurer la base de données** :
   - Restaurer depuis le backup créé avant la migration

## Migration de données volumineuses

Pour les migrations de données importantes :

1. Utiliser des scripts Python/Go dédiés
2. Traiter par lots (batch processing)
3. Utiliser des transactions pour garantir la cohérence
4. Monitorer la progression
5. Avoir un mécanisme de reprise en cas d'échec

## Exemple de script de migration de données

```python
# infra/devops_scripts/migrate_data.py
import psycopg2
from psycopg2.extras import execute_batch

def migrate_users():
    # Connexion à la base de données
    conn = psycopg2.connect(...)
    cur = conn.cursor()
    
    # Migration par lots
    batch_size = 1000
    offset = 0
    
    while True:
        cur.execute("SELECT * FROM old_users LIMIT %s OFFSET %s", (batch_size, offset))
        users = cur.fetchall()
        
        if not users:
            break
        
        # Traitement et insertion
        execute_batch(cur, insert_query, users)
        conn.commit()
        
        offset += batch_size
        print(f"Migré {offset} utilisateurs...")
    
    cur.close()
    conn.close()
```

