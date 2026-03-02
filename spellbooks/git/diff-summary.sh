#!/bin/sh

BOLD='\033[1m'
CYAN='\033[36m'
DIM='\033[2m'
RESET='\033[0m'

has_staged=$(git diff --cached --stat)
has_unstaged=$(git diff --stat)

if [ -z "$has_staged" ] && [ -z "$has_unstaged" ]; then
    printf "${DIM}No changes to show.${RESET}\n"
    exit 0
fi

if [ -n "$has_staged" ]; then
    printf "${BOLD}${CYAN} Staged changes:${RESET}\n"
    git diff --cached --stat --color=always
    echo ""
    git diff --cached --color=always | head -80
    echo ""
fi

if [ -n "$has_unstaged" ]; then
    printf "${BOLD}${CYAN} Unstaged changes:${RESET}\n"
    git diff --stat --color=always
    echo ""
    git diff --color=always | head -80
fi
