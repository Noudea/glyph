#!/bin/sh
set -e

BOLD='\033[1m'
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
DIM='\033[2m'
RESET='\033[0m'

containers=$(docker ps -a --format '{{.Names}}')

if [ -z "$containers" ]; then
    printf "${DIM}No containers.${RESET}\n"
    exit 0
fi

printf "${BOLD}${CYAN} All Containers:${RESET}\n\n"
i=1
for name in $containers; do
    status=$(docker inspect --format '{{.State.Status}}' "$name")
    if [ "$status" = "running" ]; then
        printf "  ${GREEN}%d)${RESET} %s ${GREEN}(%s)${RESET}\n" "$i" "$name" "$status"
    else
        printf "  ${DIM}%d)${RESET} %s ${DIM}(%s)${RESET}\n" "$i" "$name" "$status"
    fi
    i=$((i + 1))
done

echo ""
printf "${CYAN}Container number: ${RESET}"
read -r num

target=$(echo "$containers" | sed -n "${num}p")

if [ -z "$target" ]; then
    printf "${DIM}Invalid selection.${RESET}\n"
    exit 1
fi

echo ""
printf "${BOLD}${YELLOW} %s${RESET}\n\n" "$target"

# Image
image=$(docker inspect --format '{{.Config.Image}}' "$target")
printf "${CYAN}Image:${RESET}   %s\n" "$image"

# Status
state=$(docker inspect --format '{{.State.Status}}' "$target")
printf "${CYAN}State:${RESET}   %s\n" "$state"

# IP Address
ip=$(docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$target")
printf "${CYAN}IP:${RESET}      %s\n" "${ip:-none}"

# Ports
ports=$(docker inspect --format '{{range $p, $conf := .NetworkSettings.Ports}}{{$p}} -> {{range $conf}}{{.HostIp}}:{{.HostPort}}{{end}} {{end}}' "$target")
printf "${CYAN}Ports:${RESET}   %s\n" "${ports:-none}"

# Mounts
printf "${CYAN}Mounts:${RESET}\n"
docker inspect --format '{{range .Mounts}}  {{.Type}}: {{.Source}} -> {{.Destination}}{{"\n"}}{{end}}' "$target"

# Env
printf "${CYAN}Env:${RESET}\n"
docker inspect --format '{{range .Config.Env}}  {{.}}{{"\n"}}{{end}}' "$target"

printf "\n${DIM}Press Enter to return...${RESET}"
read -r _
