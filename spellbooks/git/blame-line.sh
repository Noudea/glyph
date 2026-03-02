#!/bin/sh
set -e

BOLD='\033[1m'
CYAN='\033[36m'
YELLOW='\033[33m'
DIM='\033[2m'
RESET='\033[0m'

printf "${CYAN}File path: ${RESET}"
read -r filepath

if [ -z "$filepath" ] || [ ! -f "$filepath" ]; then
    printf "${DIM}File not found: %s${RESET}\n" "$filepath"
    exit 1
fi

printf "${CYAN}Line range (e.g. 10,20): ${RESET}"
read -r range

if [ -z "$range" ]; then
    printf "${DIM}Aborted — no range given.${RESET}\n"
    exit 1
fi

echo ""
printf "${BOLD}${YELLOW} Blame for %s (lines %s):${RESET}\n\n" "$filepath" "$range"
git blame -L "$range" --color-by-age --color-lines "$filepath"
