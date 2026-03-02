#!/bin/sh
set -e

BOLD='\033[1m'
GREEN='\033[32m'
YELLOW='\033[33m'
RED='\033[31m'
CYAN='\033[36m'
DIM='\033[2m'
RESET='\033[0m'

branch=$(git branch --show-current 2>/dev/null || echo "detached")
printf "${BOLD}${CYAN} Branch:${RESET} %s\n" "$branch"

# Ahead/behind
upstream=$(git rev-parse --abbrev-ref '@{upstream}' 2>/dev/null || true)
if [ -n "$upstream" ]; then
    ahead=$(git rev-list --count '@{upstream}..HEAD' 2>/dev/null || echo 0)
    behind=$(git rev-list --count 'HEAD..@{upstream}' 2>/dev/null || echo 0)
    if [ "$ahead" -gt 0 ] && [ "$behind" -gt 0 ]; then
        printf "${YELLOW} %s ahead, %s behind${RESET} %s\n" "$ahead" "$behind" "$upstream"
    elif [ "$ahead" -gt 0 ]; then
        printf "${GREEN} %s ahead${RESET} of %s\n" "$ahead" "$upstream"
    elif [ "$behind" -gt 0 ]; then
        printf "${RED} %s behind${RESET} %s\n" "$behind" "$upstream"
    else
        printf "${DIM} Up to date with %s${RESET}\n" "$upstream"
    fi
else
    printf "${DIM} No upstream tracking branch${RESET}\n"
fi

echo ""

# Staged
staged=$(git diff --cached --name-status)
if [ -n "$staged" ]; then
    printf "${GREEN}${BOLD} Staged:${RESET}\n"
    echo "$staged" | while IFS=$(printf '\t') read -r status file; do
        printf "  ${GREEN}%-8s${RESET} %s\n" "$status" "$file"
    done
    echo ""
fi

# Unstaged
unstaged=$(git diff --name-status)
if [ -n "$unstaged" ]; then
    printf "${YELLOW}${BOLD} Unstaged:${RESET}\n"
    echo "$unstaged" | while IFS=$(printf '\t') read -r status file; do
        printf "  ${YELLOW}%-8s${RESET} %s\n" "$status" "$file"
    done
    echo ""
fi

# Untracked
untracked=$(git ls-files --others --exclude-standard)
if [ -n "$untracked" ]; then
    printf "${RED}${BOLD} Untracked:${RESET}\n"
    echo "$untracked" | while read -r file; do
        printf "  ${RED}?${RESET}        %s\n" "$file"
    done
    echo ""
fi

if [ -z "$staged" ] && [ -z "$unstaged" ] && [ -z "$untracked" ]; then
    printf "${DIM} Working tree clean${RESET}\n"
fi
