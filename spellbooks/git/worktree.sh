#!/bin/sh
set -e

BOLD='\033[1m'
GREEN='\033[32m'
CYAN='\033[36m'
YELLOW='\033[33m'
RED='\033[31m'
DIM='\033[2m'
RESET='\033[0m'

# Show existing worktrees
printf "${BOLD}${CYAN} Worktrees:${RESET}\n\n"
git worktree list | while IFS= read -r line; do
    printf "  %s\n" "$line"
done
echo ""

printf "[${GREEN}a${RESET}]dd  [${RED}r${RESET}]emove  [${YELLOW}q${RESET}]uit: "
read -r action

case "$action" in
    a|add)
        printf "${CYAN}Branch name: ${RESET}"
        read -r branch
        if [ -z "$branch" ]; then
            printf "${DIM}Aborted.${RESET}\n"
            exit 1
        fi
        path="../worktree-${branch}"
        printf "${CYAN}Path [%s]: ${RESET}" "$path"
        read -r custom_path
        if [ -n "$custom_path" ]; then
            path="$custom_path"
        fi
        git worktree add "$path" -b "$branch" 2>&1 || git worktree add "$path" "$branch" 2>&1
        printf "\n${GREEN} Worktree created at %s${RESET}\n" "$path"
        ;;
    r|remove)
        printf "${CYAN}Worktree path to remove: ${RESET}"
        read -r path
        if [ -z "$path" ]; then
            printf "${DIM}Aborted.${RESET}\n"
            exit 1
        fi
        git worktree remove "$path"
        printf "\n${RED} Removed worktree at %s${RESET}\n" "$path"
        ;;
    q|quit)
        printf "${DIM}Done.${RESET}\n"
        ;;
    *)
        printf "${DIM}Unknown action.${RESET}\n"
        exit 1
        ;;
esac
