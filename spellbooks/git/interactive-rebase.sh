#!/bin/sh
set -e

BOLD='\033[1m'
CYAN='\033[36m'
YELLOW='\033[33m'
DIM='\033[2m'
RESET='\033[0m'

printf "${DIM}Recent commits:${RESET}\n\n"
git log --oneline --color=always \
    --format="%C(bold blue)%h%C(reset) %C(white)%s%C(reset) %C(dim)(%cr)%C(reset)" \
    -10
echo ""

printf "${CYAN}How many commits to rebase? ${RESET}"
read -r count

if [ -z "$count" ]; then
    printf "${DIM}Aborted.${RESET}\n"
    exit 1
fi

printf "\n${YELLOW}${BOLD} Opening interactive rebase for last %s commit(s)...${RESET}\n\n"  "$count"
exec git rebase -i "HEAD~${count}"
