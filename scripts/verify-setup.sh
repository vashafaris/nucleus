#!/bin/bash

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}=== Nucleus Step 1 Verification ===${NC}\n"

# Function to check if command exists
check_command() {
  if command -v $1 &>/dev/null; then
    echo -e "✅ $1 is installed"
  else
    echo -e "❌ $1 is not installed"
  fi
}

# Function to check if service is running
check_service() {
  local service=$1
  local port=$2
  if nc -z localhost $port 2>/dev/null; then
    echo -e "✅ $service is running on port $port"
  else
    echo -e "❌ $service is not accessible on port $port"
  fi
}

# Function to check if directory exists
check_directory() {
  if [ -d "$1" ]; then
    echo -e "✅ Directory exists: $1"
  else
    echo -e "❌ Directory missing: $1"
  fi
}

echo -e "${YELLOW}1. Checking Prerequisites:${NC}"
check_command go
check_command docker
check_command docker-compose
check_command make

echo -e "\n${YELLOW}2. Checking Project Structure:${NC}"
check_directory "cmd/api"
check_directory "internal/domain"
check_directory "internal/application"
check_directory "internal/infrastructure"
check_directory "pkg/config"

echo -e "\n${YELLOW}3. Checking Configuration Files:${NC}"
[ -f ".env" ] && echo -e "✅ .env file exists" || echo -e "❌ .env file missing"
[ -f "docker-compose.yml" ] && echo -e "✅ docker-compose.yml exists" || echo -e "❌ docker-compose.yml missing"
[ -f "Makefile" ] && echo -e "✅ Makefile exists" || echo -e "❌ Makefile missing"
[ -f "go.mod" ] && echo -e "✅ go.mod exists" || echo -e "❌ go.mod missing"

echo -e "\n${YELLOW}4. Checking Running Services:${NC}"
check_service "PostgreSQL" 5433
check_service "Redis" 6379
check_service "RabbitMQ" 5672
check_service "RabbitMQ Management" 15672
check_service "Kafka" 9092
check_service "Keycloak" 8180
check_service "Prometheus" 9090
check_service "Grafana" 3000
check_service "Loki" 3101

echo -e "\n${YELLOW}5. Checking Docker Containers:${NC}"
running_containers=$(docker-compose ps --services --filter "status=running" | wc -l)
total_containers=$(docker-compose ps --services | wc -l)
echo -e "Running containers: $running_containers/$total_containers"

if [ $running_containers -eq $total_containers ]; then
  echo -e "${GREEN}✅ All containers are running!${NC}"
else
  echo -e "${RED}❌ Some containers are not running${NC}"
  echo -e "${YELLOW}Not running:${NC}"
  docker-compose ps --services --filter "status=exited"
fi

echo -e "\n${YELLOW}6. Testing Basic Endpoints:${NC}"
# Test Redis
redis_test=$(docker exec nucleus-redis redis-cli -a redis123 ping 2>/dev/null)
[ "$redis_test" == "PONG" ] && echo -e "✅ Redis is responding" || echo -e "❌ Redis test failed"

# Test PostgreSQL
pg_test=$(docker exec nucleus-postgres psql -U nucleus -d nucleus_db -c "SELECT 1;" 2>/dev/null | grep -c "1 row")
[ $pg_test -eq 1 ] && echo -e "✅ PostgreSQL is responding" || echo -e "❌ PostgreSQL test failed"

echo -e "\n${GREEN}=== Step 1 Verification Complete ===${NC}"
echo -e "\nIf all checks pass, you're ready to proceed to Step 2!"
echo -e "Run: ${YELLOW}make docker-logs${NC} to see service logs if any service is failing."
