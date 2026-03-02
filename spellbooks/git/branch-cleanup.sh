#!/bin/sh
set -e

BOLD='\033[1m'
GREEN='\033[32m'
YELLOW='\033[33m'
RED='\033[31m'
DIM='\033[2m'
RESET='\033[0m'

current=$(git branch --show-current)

# Prune stale remote-tracking branches
printf "${DIM}Pruning stale remote references...${RESET}\n"
git remote prune origin 2>/dev/null || true
echo ""

merged=$(git branch --merged | grep -v "^\*" | grep -v " main$" | grep -v " master$" | sed 's/^[ ]*//')

if [ -z "$merged" ]; then
    printf "${GREEN}No merged branches to clean up.${RESET}\n"
    exit 0
fi

count=$(echo "$merged" | wc -l | tr -d ' ')
printf "${BOLD}${YELLOW} %s merged branch(es)${RESET} ${DIM}(current: %s)${RESET}\n\n" "$count" "$current"

echo "$merged" | while read -r branch; do
    printf "  ${RED}%s${RESET}\n" "$branch"
done

echo ""
printf "Delete these branches? [y/N] "
read -r confirm

if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
    echo "$merged" | while read -r branch; do
        git branch -d "$branch" 2>&1
    done
    printf "\n${GREEN} Cleanup complete.${RESET}\n"
else
    printf "${DIM}Aborted.${RESET}\n"
fi
