#!/bin/sh
set -e

stash_count=$(git stash list | wc -l | tr -d ' ')

if [ "$stash_count" -eq 0 ]; then
    echo "No stashes found."
    exit 0
fi

echo "=== Stash List ==="
git stash list
echo ""

printf "Enter stash number (0-%d): " "$((stash_count - 1))"
read -r num

printf "Action? [a]pply / [p]op / [d]rop: "
read -r action

case "$action" in
    a|apply)  git stash apply "stash@{${num}}" ;;
    p|pop)    git stash pop "stash@{${num}}" ;;
    d|drop)   git stash drop "stash@{${num}}" ;;
    *)        echo "Unknown action: $action"; exit 1 ;;
esac
