#!/bin/sh
set -e

BOLD='\033[1m'
GREEN='\033[32m'
CYAN='\033[36m'
YELLOW='\033[33m'
DIM='\033[2m'
RESET='\033[0m'

# Show last commit
printf "${BOLD}${YELLOW} Last commit:${RESET}\n"
git log -1 --oneline --color=always
echo ""

# Show current changes
staged=$(git diff --cached --name-status)
unstaged=$(git diff --name-status)

if [ -n "$staged" ]; then
    printf "${GREEN}Staged changes will be included:${RESET}\n"
    echo "$staged"
    echo ""
fi

if [ -n "$unstaged" ]; then
    printf "${DIM}Unstaged changes (won't be included unless you stage them):${RESET}\n"
    echo "$unstaged"
    echo ""
fi

printf "[${GREEN}e${RESET}]dit message  [${CYAN}k${RESET}]eep message  [${YELLOW}q${RESET}]uit: "
read -r action

case "$action" in
    e|edit)
        exec git commit --amend
        ;;
    k|keep)
        git commit --amend --no-edit
        printf "\n${GREEN} Amended (message kept).${RESET}\n"
        ;;
    q|quit)
        printf "${DIM}Aborted.${RESET}\n"
        ;;
    *)
        printf "${DIM}Unknown action. Aborted.${RESET}\n"
        exit 1
        ;;
esac
