#!/bin/sh

BOLD='\033[1m'
CYAN='\033[36m'
DIM='\033[2m'
RESET='\033[0m'

running=$(docker ps -q | wc -l | tr -d ' ')

if [ "$running" -eq 0 ]; then
    printf "${DIM}No running containers.${RESET}\n"
    exit 0
fi

printf "${BOLD}${CYAN} Resource Usage${RESET} ${DIM}(%s containers)${RESET}\n\n" "$running"
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}\t{{.NetIO}}\t{{.BlockIO}}"
