#!/bin/sh
set -e

BOLD='\033[1m'
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
DIM='\033[2m'
RESET='\033[0m'

printf "${BOLD}${YELLOW} Rebuilding${RESET} ${DIM}(no cache)...${RESET}\n\n"
docker compose build --no-cache

printf "\n${BOLD}${CYAN} Starting services...${RESET}\n\n"
docker compose up -d

printf "\n${GREEN} Done.${RESET} Running services:\n\n"
docker compose ps
