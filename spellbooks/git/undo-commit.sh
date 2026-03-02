#!/bin/sh
set -e

BOLD='\033[1m'
YELLOW='\033[33m'
GREEN='\033[32m'
DIM='\033[2m'
RESET='\033[0m'

# Show what will be undone
last=$(git log -1 --oneline --color=always)
printf "${BOLD}${YELLOW} Undoing last commit:${RESET}\n"
printf "  %s\n\n" "$last"

printf "Soft reset HEAD~1? [y/N] "
read -r confirm

if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
    git reset --soft HEAD~1
    printf "\n${GREEN} Commit undone. Changes are staged:${RESET}\n\n"
    git diff --cached --stat
else
    printf "${DIM}Aborted.${RESET}\n"
fi
