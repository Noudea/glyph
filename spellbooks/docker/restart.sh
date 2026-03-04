#!/bin/sh
set -e

BOLD='\033[1m'
CYAN='\033[36m'
GREEN='\033[32m'
DIM='\033[2m'
RESET='\033[0m'

containers=$(docker ps --format '{{.Names}}')

if [ -z "$containers" ]; then
    printf "${DIM}No running containers.${RESET}\n"
    exit 0
fi

printf "${BOLD}${CYAN} Running Containers:${RESET}\n\n"
i=1
for name in $containers; do
    printf "  ${GREEN}%d)${RESET} %s\n" "$i" "$name"
    i=$((i + 1))
done

echo ""
printf "${CYAN}Container number to restart: ${RESET}"
read -r num

target=$(echo "$containers" | sed -n "${num}p")

if [ -z "$target" ]; then
    printf "${DIM}Invalid selection.${RESET}\n"
    exit 1
fi

printf "\n${CYAN}Restarting %s...${RESET}\n" "$target"
docker restart "$target"
printf "\n${GREEN} %s restarted.${RESET}\n" "$target"
