#!/bin/sh
set -e

BOLD='\033[1m'
GREEN='\033[32m'
CYAN='\033[36m'
DIM='\033[2m'
RESET='\033[0m'

current=$(git branch --show-current 2>/dev/null || echo "detached")
printf "${DIM}Current branch: %s${RESET}\n\n" "$current"

printf "${CYAN}New branch name: ${RESET}"
read -r name

if [ -z "$name" ]; then
    printf "${DIM}Aborted — empty name.${RESET}\n"
    exit 1
fi

git checkout -b "$name"
printf "\n${GREEN} Switched to new branch '${name}'.${RESET}\n"
