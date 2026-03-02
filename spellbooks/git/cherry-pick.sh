#!/bin/sh
set -e

BOLD='\033[1m'
GREEN='\033[32m'
CYAN='\033[36m'
YELLOW='\033[33m'
DIM='\033[2m'
RESET='\033[0m'

printf "${DIM}Recent commits across all branches:${RESET}\n\n"
git log --all --oneline --color=always --format="%C(bold blue)%h%C(reset) %C(white)%s%C(reset) %C(dim)(%cr, %an)%C(reset)" -15
echo ""

printf "${CYAN}Commit hash to cherry-pick: ${RESET}"
read -r hash

if [ -z "$hash" ]; then
    printf "${DIM}Aborted — empty hash.${RESET}\n"
    exit 1
fi

# Show what we're picking
printf "\n${YELLOW}Cherry-picking:${RESET}\n"
git log --oneline --color=always -1 "$hash"
echo ""

git cherry-pick "$hash"
printf "\n${GREEN} Cherry-pick applied.${RESET}\n"
