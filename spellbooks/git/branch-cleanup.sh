#!/bin/sh
set -e

current=$(git branch --show-current)
merged=$(git branch --merged | grep -v "^\*" | grep -v "main" | grep -v "master" | sed 's/^[ ]*//')

if [ -z "$merged" ]; then
    echo "No merged branches to clean up."
    exit 0
fi

echo "Merged branches (current: $current):"
echo "$merged"
echo ""
printf "Delete these branches? [y/N] "
read -r confirm

if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
    echo "$merged" | while read -r branch; do
        git branch -d "$branch"
    done
    echo "Done."
else
    echo "Aborted."
fi
