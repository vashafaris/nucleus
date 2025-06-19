#!/bin/bash

# Load environment variables
source .env

# Database connection string
DB_DSN="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

# Migration directory
MIGRATION_DIR="./deployments/migrations"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Create migrations directory if it doesn't exist
mkdir -p ${MIGRATION_DIR}

# Function to run migrations
run_migration() {
  local direction=$1

  echo -e "${GREEN}Running migrations ${direction}...${NC}"

  migrate -path ${MIGRATION_DIR} -database "${DB_DSN}" ${direction}

  if [ $? -eq 0 ]; then
    echo -e "${GREEN}Migrations ${direction} completed successfully!${NC}"
  else
    echo -e "${RED}Migration ${direction} failed!${NC}"
    exit 1
  fi
}

# Check command argument
case "$1" in
up)
  run_migration "up"
  ;;
down)
  run_migration "down 1"
  ;;
drop)
  echo -e "${RED}Dropping all migrations...${NC}"
  migrate -path ${MIGRATION_DIR} -database "${DB_DSN}" drop -f
  ;;
create)
  if [ -z "$2" ]; then
    echo -e "${RED}Please provide a migration name${NC}"
    echo "Usage: ./scripts/migrate.sh create <migration_name>"
    exit 1
  fi
  migrate create -ext sql -dir ${MIGRATION_DIR} -seq $2
  echo -e "${GREEN}Migration files created: ${MIGRATION_DIR}/*_$2.{up,down}.sql${NC}"
  ;;
*)
  echo "Usage: ./scripts/migrate.sh {up|down|drop|create <name>}"
  exit 1
  ;;
esac
