#!/bin/sh
set -e

BOLD='\033[1m'
GREEN='\033[32m'
CYAN='\033[36m'
RED='\033[31m'
DIM='\033[2m'
RESET='\033[0m'

branch=$(git branch --show-current)
if [ -z "$branch" ]; then
    printf "${RED}Not on a branch (detached HEAD).${RESET}\n"
    exit 1
fi

remote=$(git config "branch.${branch}.remote" 2>/dev/null || echo "origin")

printf "${CYAN}${BOLD} Pulling${RESET} ${DIM}(rebase from ${remote}/${branch})${RESET}\n"
git pull --rebase "$remote" "$branch" 2>&1

printf "\n${CYAN}${BOLD} Pushing${RESET} ${DIM}(to ${remote}/${branch})${RESET}\n"
git push "$remote" "$branch" 2>&1

printf "\n${GREEN} Branch synced.${RESET}\n"
