#!/bin/sh

BOLD='\033[1m'
CYAN='\033[36m'
DIM='\033[2m'
RESET='\033[0m'

count=$(docker images -q | wc -l | tr -d ' ')
printf "${BOLD}${CYAN} Images${RESET} ${DIM}(%s total)${RESET}\n\n" "$count"

if [ "$count" -eq 0 ]; then
    printf "${DIM}No images.${RESET}\n"
    exit 0
fi

docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedSince}}\t{{.ID}}"
