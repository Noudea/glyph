#!/bin/sh

BOLD='\033[1m'
CYAN='\033[36m'
DIM='\033[2m'
RESET='\033[0m'

if ! docker compose version > /dev/null 2>&1; then
    printf "${DIM}Docker Compose not available.${RESET}\n"
    exit 1
fi

printf "${BOLD}${CYAN} Compose Services:${RESET}\n\n"
docker compose ps -a 2>/dev/null

if [ $? -ne 0 ] || [ -z "$(docker compose ps -q 2>/dev/null)" ]; then
    printf "${DIM}No compose project found in current directory.${RESET}\n"
fi
