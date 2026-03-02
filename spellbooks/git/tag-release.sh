#!/bin/sh
set -e

BOLD='\033[1m'
GREEN='\033[32m'
CYAN='\033[36m'
YELLOW='\033[33m'
DIM='\033[2m'
RESET='\033[0m'

latest=$(git describe --tags --abbrev=0 2>/dev/null || true)
if [ -n "$latest" ]; then
    printf "${DIM}Latest tag: ${YELLOW}%s${RESET}\n\n" "$latest"
else
    printf "${DIM}No existing tags.${RESET}\n\n"
fi

printf "${CYAN}New tag (e.g. v1.0.0): ${RESET}"
read -r tag

if [ -z "$tag" ]; then
    printf "${DIM}Aborted — empty tag.${RESET}\n"
    exit 1
fi

git tag "$tag"
printf "${GREEN} Created tag '${tag}'.${RESET}\n\n"

printf "Push tag to remote? [y/N] "
read -r confirm

if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
    git push origin "$tag"
    printf "\n${GREEN} Pushed '${tag}' to origin.${RESET}\n"
else
    printf "${DIM}Tag created locally only.${RESET}\n"
fi
