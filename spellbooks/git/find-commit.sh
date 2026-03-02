#!/bin/sh
set -e

BOLD='\033[1m'
CYAN='\033[36m'
YELLOW='\033[33m'
DIM='\033[2m'
RESET='\033[0m'

printf "${CYAN}Search commit messages for: ${RESET}"
read -r query

if [ -z "$query" ]; then
    printf "${DIM}Aborted — empty query.${RESET}\n"
    exit 1
fi

results=$(git log --all --oneline --color=always \
    --format="%C(bold blue)%h%C(reset) %C(white)%s%C(reset) %C(dim)(%cr)%C(reset)" \
    --grep="$query" -i -20)

if [ -z "$results" ]; then
    printf "\n${DIM}No commits matching '%s'.${RESET}\n" "$query"
    exit 0
fi

count=$(echo "$results" | wc -l | tr -d ' ')
printf "\n${BOLD}${YELLOW} %s result(s) for '%s':${RESET}\n\n" "$count" "$query"
echo "$results"
printf "\n${DIM}Press Enter to return...${RESET}"
read -r _
