#!/bin/sh

BOLD='\033[1m'
CYAN='\033[36m'
DIM='\033[2m'
RESET='\033[0m'

total=$(git rev-list --count HEAD)
printf "${BOLD}${CYAN} Contributors${RESET} ${DIM}(%s total commits)${RESET}\n\n" "$total"

git shortlog -sne HEAD | while IFS= read -r line; do
    count=$(echo "$line" | sed 's/^[[:space:]]*//' | cut -f1)
    name=$(echo "$line" | sed 's/^[[:space:]]*//' | cut -f2-)
    pct=$((count * 100 / total))
    printf "  %4s  %-40s %s%%\n" "$count" "$name" "$pct"
done
