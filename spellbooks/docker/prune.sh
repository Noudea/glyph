#!/bin/sh
set -e

BOLD='\033[1m'
YELLOW='\033[33m'
RED='\033[31m'
GREEN='\033[32m'
DIM='\033[2m'
RESET='\033[0m'

printf "${BOLD}${YELLOW} Docker Disk Usage:${RESET}\n\n"
docker system df
echo ""

printf "${RED}This will remove all unused containers, networks, and images.${RESET}\n"
printf "Also prune volumes? [y/N] "
read -r volumes

if [ "$volumes" = "y" ] || [ "$volumes" = "Y" ]; then
    printf "\n${YELLOW}Pruning everything including volumes...${RESET}\n\n"
    docker system prune -a --volumes -f
else
    printf "\n${YELLOW}Pruning (keeping volumes)...${RESET}\n\n"
    docker system prune -a -f
fi

printf "\n${GREEN} Done.${RESET} New disk usage:\n\n"
docker system df
