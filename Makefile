# Makefile for managing the GoAuth service lifecycle.
# This acts as a simple wrapper around the more detailed service.sh script.

# Path to the main management script.
SERVICE_SCRIPT = ./deploy/goauth/development/service.sh

# --- Phony Targets ---
# .PHONY ensures that these commands run even if a file with the same name exists.
.PHONY: help up stop down logs up-db run stop-db down-db logs-db

# --- Default Target ---
# Running 'make' without any arguments will show the help message.
help:
	@$(SERVICE_SCRIPT) help

# --- Full Dockerized Environment Commands ---

# Builds and starts the full application (app + db) in the background.
up:
	@echo "--> Starting full application stack (app + db)..."
	@$(SERVICE_SCRIPT) up

# Stops the full application stack without deleting data.
stop:
	@echo "--> Stopping full application stack..."
	@$(SERVICE_SCRIPT) stop

# Stops and removes the full application stack and its data volumes.
down:
	@echo "--> Tearing down full application stack and data..."
	@$(SERVICE_SCRIPT) down

# Follows the logs for the full application stack.
logs:
	@echo "--> Tailing logs for full application stack..."
	@$(SERVICE_SCRIPT) logs

# --- Local Development Helper Commands ---

# Starts only the database service for local development.
up-db:
	@echo "--> Starting standalone database service..."
	@$(SERVICE_SCRIPT) up-db

# Starts the Go service locally (requires the database to be running).
run:
	@echo "--> Running Go service locally..."
	@$(SERVICE_SCRIPT) run

# Stops the standalone database service.
stop-db:
	@echo "--> Stopping standalone database service..."
	@$(SERVICE_SCRIPT) stop-db

# Stops and removes the standalone database service and its data.
down-db:
	@echo "--> Tearing down standalone database and data..."
	@$(SERVICE_SCRIPT) down-db

# Follows the logs for the standalone database service.
logs-db:
	@echo "--> Tailing logs for standalone database service..."
	@$(SERVICE_SCRIPT) logs-db

