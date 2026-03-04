#!/bin/sh
set -e

BOLD='\033[1m'
RED='\033[31m'
GREEN='\033[32m'
YELLOW='\033[33m'
DIM='\033[2m'
RESET='\033[0m'

containers=$(docker ps --format '{{.Names}}')

if [ -z "$containers" ]; then
    printf "${DIM}No running containers.${RESET}\n"
    exit 0
fi

count=$(echo "$containers" | wc -l | tr -d ' ')
printf "${BOLD}${YELLOW} %s running container(s):${RESET}\n\n" "$count"
echo "$containers" | while read -r name; do
    printf "  ${RED}%s${RESET}\n" "$name"
done

echo ""
printf "Stop all? [y/N] "
read -r confirm

if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
    docker stop $(docker ps -q)
    printf "\n${GREEN} All containers stopped.${RESET}\n"
else
    printf "${DIM}Aborted.${RESET}\n"
fi
