#!/bin/sh

BOLD='\033[1m'
GREEN='\033[32m'
RED='\033[31m'
CYAN='\033[36m'
DIM='\033[2m'
RESET='\033[0m'

running=$(docker ps -q | wc -l | tr -d ' ')
stopped=$(docker ps -aq --filter "status=exited" | wc -l | tr -d ' ')

printf "${BOLD}${CYAN} Containers${RESET} ${DIM}(%s running, %s stopped)${RESET}\n\n" "$running" "$stopped"

if [ "$running" -gt 0 ]; then
    printf "${GREEN}${BOLD}Running:${RESET}\n"
    docker ps --format "  {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}" | column -t -s "$(printf '\t')"
    echo ""
fi

if [ "$stopped" -gt 0 ]; then
    printf "${RED}${BOLD}Stopped:${RESET}\n"
    docker ps -a --filter "status=exited" --format "  {{.Names}}\t{{.Image}}\t{{.Status}}" | column -t -s "$(printf '\t')"
    echo ""
fi

if [ "$running" -eq 0 ] && [ "$stopped" -eq 0 ]; then
    printf "${DIM}No containers.${RESET}\n"
fi
