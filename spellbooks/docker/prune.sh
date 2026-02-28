#!/bin/sh
set -e

echo "This will remove all unused containers, networks, images, and optionally volumes."
printf "Continue? [y/N] "
read -r confirm

if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
    docker system prune -a
else
    echo "Aborted."
fi
