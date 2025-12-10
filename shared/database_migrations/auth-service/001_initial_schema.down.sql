-- Rollback de la migration initiale pour auth-service

DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;

