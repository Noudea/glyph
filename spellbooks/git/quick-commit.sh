#!/bin/sh
set -e

BOLD='\033[1m'
GREEN='\033[32m'
YELLOW='\033[33m'
DIM='\033[2m'
RESET='\033[0m'

# Check for changes
if git diff --quiet && git diff --cached --quiet && [ -z "$(git ls-files --others --exclude-standard)" ]; then
    printf "${DIM}Nothing to commit — working tree clean.${RESET}\n"
    exit 0
fi

# Show what will be committed
printf "${BOLD}Changes to be committed:${RESET}\n"
git status --short
echo ""

printf "${YELLOW}Commit message: ${RESET}"
read -r msg

if [ -z "$msg" ]; then
    printf "${DIM}Aborted — empty message.${RESET}\n"
    exit 1
fi

git add -A
git commit -m "$msg"

printf "\n${GREEN} Committed.${RESET}\n"
