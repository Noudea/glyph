#!/bin/sh

BOLD='\033[1m'
CYAN='\033[36m'
GREEN='\033[32m'
DIM='\033[2m'
RESET='\033[0m'

printf "${BOLD}${CYAN} Docker Networks:${RESET}\n\n"

docker network ls --format '{{.Name}}\t{{.Driver}}\t{{.Scope}}' | while IFS=$(printf '\t') read -r name driver scope; do
    printf "${GREEN}%-30s${RESET} ${DIM}%-10s %s${RESET}\n" "$name" "$driver" "$scope"

    # Show connected containers
    connected=$(docker network inspect --format '{{range .Containers}}  {{.Name}}{{"\n"}}{{end}}' "$name" 2>/dev/null)
    if [ -n "$connected" ]; then
        echo "$connected" | while read -r container; do
            [ -n "$container" ] && printf "  ${DIM}└─ %s${RESET}\n" "$container"
        done
    fi
done

printf "\n${DIM}Press Enter to return...${RESET}"
read -r _
