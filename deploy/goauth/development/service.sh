#!/bin/bash

# A unified script to manage the GoAuth service lifecycle.
# It handles different running modes for development, including running the full
# stack with Docker, running only the database, or running the Go service locally.

set -e

# --- Configuration ---
# Path to the full docker-compose file (app + db)
COMPOSE_FILE_FULL="deploy/goauth/development/docker-compose.yml"
# Path to the docker-compose file for running only the database
COMPOSE_FILE_DB_ONLY="deploy/goauth/development/docker-compose.no-service.yml"
# Go service configuration
SERVICE_NAME="goauth"
SERVICE_PATH="./cmd/${SERVICE_NAME}/main.go"

COMMAND=$1

# --- Help Function ---
show_help() {
    echo "Usage: ./service.sh <command>"
    echo ""
    echo "This script is the single entry point for managing the GoAuth service."
    echo ""
    echo "Available commands:"
    echo "  --- Full Dockerized Environment ---"
    echo "  up      Builds and starts the full application (app + db) using ${COMPOSE_FILE_FULL}."
    echo "  stop    Stops the full application stack without deleting data."
    echo "  down    Stops and removes the full application stack and its data volumes."
    echo "  logs    Follows the logs for the full application stack."
    echo ""
    echo "  --- Local Development Helpers ---"
    echo "  up-db   Starts only the database service using ${COMPOSE_FILE_DB_ONLY}."
    echo "  run     Starts the Go service locally (requires the database to be running)."
    echo "  stop-db Stops the standalone database service."
    echo "  down-db Stops and removes the standalone database service and its data."
    echo "  logs-db Follows the logs for the standalone database service."
    echo ""
    echo "  --- General ---"
    echo "  help    Shows this help message."
}

# --- Command Logic ---
case "$COMMAND" in
    # --- Full Dockerized Environment Commands ---
    up)
        echo "--> Starting the full application stack from ${COMPOSE_FILE_FULL}..."
        docker compose -f "${COMPOSE_FILE_FULL}" up --build -d
        echo "--> Full stack is up and running."
        ;;
    stop)
        echo "--> Stopping the full application stack from ${COMPOSE_FILE_FULL}..."
        docker compose -f "${COMPOSE_FILE_FULL}" stop
        echo "--> Full stack stopped."
        ;;
    down)
        echo "--> Tearing down the full application stack from ${COMPOSE_FILE_FULL}..."
        docker compose -f "${COMPOSE_FILE_FULL}" down -v
        echo "--> Full stack and all associated data have been removed."
        ;;
    logs)
        echo "--> Following logs for the full application stack from ${COMPOSE_FILE_FULL}..."
        docker compose -f "${COMPOSE_FILE_FULL}" logs -f
        ;;

    # --- Local Development Helper Commands ---
    up-db)
        echo "--> Starting only the database service from ${COMPOSE_FILE_DB_ONLY}..."
        docker compose -f "${COMPOSE_FILE_DB_ONLY}" up -d
        echo "--> Database service is up and running."
        ;;
    run)
        echo "--> Starting the Go service locally: ${SERVICE_NAME}..."
        echo "--> Ensure the database is running (e.g., via './service.sh up-db')."
        go run "${SERVICE_PATH}" serve --migrate-up
        ;;
    stop-db)
        echo "--> Stopping the standalone database from ${COMPOSE_FILE_DB_ONLY}..."
        docker compose -f "${COMPOSE_FILE_DB_ONLY}" stop
        echo "--> Standalone database stopped."
        ;;
    down-db)
        echo "--> Tearing down the standalone database from ${COMPOSE_FILE_DB_ONLY}..."
        docker compose -f "${COMPOSE_FILE_DB_ONLY}" down -v
        echo "--> Standalone database and its data have been removed."
        ;;
    logs-db)
        echo "--> Following logs for the standalone database from ${COMPOSE_FILE_DB_ONLY}..."
        docker compose -f "${COMPOSE_FILE_DB_ONLY}" logs -f
        ;;

    # --- Help Command ---
    help|--help|-h|*)
        show_help
        exit 1
        ;;
esac

exit 0
