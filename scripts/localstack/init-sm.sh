#!/bin/sh

echo "Initializing Secrets Manager..."

awslocal secretsmanager create-secret \
    --name db-secret-url \
    --description "DB URL" \
    --secret-string "postgres://production:production@localhost:5432/production_db?sslmode=disable"