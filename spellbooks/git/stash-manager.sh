#!/bin/sh
set -e

BOLD='\033[1m'
CYAN='\033[36m'
YELLOW='\033[33m'
GREEN='\033[32m'
RED='\033[31m'
DIM='\033[2m'
RESET='\033[0m'

stash_list=$(git stash list 2>/dev/null || true)

if [ -z "$stash_list" ]; then
    printf "${DIM}No stashes found.${RESET}\n"
    exit 0
fi

stash_count=$(echo "$stash_list" | wc -l | tr -d ' ')

printf "${BOLD}${CYAN} Stashes (${stash_count}):${RESET}\n\n"
echo "$stash_list" | while IFS= read -r line; do
    num=$(echo "$line" | sed 's/stash@{\([0-9]*\)}.*/\1/')
    msg=$(echo "$line" | sed 's/stash@{[0-9]*}: //')
    printf "  ${YELLOW}[%s]${RESET} %s\n" "$num" "$msg"
done

echo ""
printf "Stash number (0-%d): " "$((stash_count - 1))"
read -r num

if [ -z "$num" ]; then
    printf "${DIM}Aborted.${RESET}\n"
    exit 0
fi

printf "[${GREEN}a${RESET}]pply  [${YELLOW}p${RESET}]op  [${RED}d${RESET}]rop: "
read -r action

case "$action" in
    a|apply)
        git stash apply "stash@{${num}}"
        printf "\n${GREEN} Applied stash@{${num}}.${RESET}\n"
        ;;
    p|pop)
        git stash pop "stash@{${num}}"
        printf "\n${GREEN} Popped stash@{${num}}.${RESET}\n"
        ;;
    d|drop)
        git stash drop "stash@{${num}}"
        printf "\n${RED} Dropped stash@{${num}}.${RESET}\n"
        ;;
    *)
        printf "${DIM}Unknown action. Aborted.${RESET}\n"
        exit 1
        ;;
esac
