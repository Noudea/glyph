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
printf "${CYAN}Container number: ${RESET}"
read -r num

target=$(echo "$containers" | sed -n "${num}p")

if [ -z "$target" ]; then
    printf "${DIM}Invalid selection.${RESET}\n"
    exit 1
fi

printf "${CYAN}Lines to show ${DIM}(default: 50):${RESET} "
read -r lines
lines="${lines:-50}"

printf "\n${BOLD}${CYAN} Logs for %s${RESET} ${DIM}(last %s lines, following):${RESET}\n\n" "$target" "$lines"
exec docker logs -f --tail "$lines" "$target"
