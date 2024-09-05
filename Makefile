# Variables
DOCKER_COMPOSE=docker-compose -f docker-compose.yaml

# Build and run the Docker containers
.PHONY: up
up:
	$(DOCKER_COMPOSE) up --build -d

# Stop the running containers
.PHONY: down
down:
	$(DOCKER_COMPOSE) down

# Restart the Docker containers
.PHONY: restart
restart:
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) up --build

# Clean up any existing containers and volumes
.PHONY: clean
clean:
	$(DOCKER_COMPOSE) down -v

# Tail logs from the app service
.PHONY: logs
logs:
	$(DOCKER_COMPOSE) logs -f app
