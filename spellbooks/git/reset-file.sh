#!/bin/sh
set -e

BOLD='\033[1m'
GREEN='\033[32m'
CYAN='\033[36m'
YELLOW='\033[33m'
RED='\033[31m'
DIM='\033[2m'
RESET='\033[0m'

# Show modified files
changed=$(git diff --name-only)
staged=$(git diff --cached --name-only)
all=$(printf "%s\n%s" "$changed" "$staged" | sort -u | sed '/^$/d')

if [ -z "$all" ]; then
    printf "${DIM}No modified files to reset.${RESET}\n"
    exit 0
fi

printf "${BOLD}${YELLOW} Modified files:${RESET}\n\n"
echo "$all" | while read -r f; do
    printf "  %s\n" "$f"
done
echo ""

printf "${CYAN}File to restore to HEAD: ${RESET}"
read -r filepath

if [ -z "$filepath" ]; then
    printf "${DIM}Aborted.${RESET}\n"
    exit 1
fi

printf "${RED}Discard changes to '%s'? [y/N] ${RESET}" "$filepath"
read -r confirm

if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
    git checkout HEAD -- "$filepath"
    printf "\n${GREEN} Restored '%s' to HEAD.${RESET}\n" "$filepath"
else
    printf "${DIM}Aborted.${RESET}\n"
fi
