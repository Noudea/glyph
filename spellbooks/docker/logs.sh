#!/bin/sh
set -e

containers=$(docker ps --format '{{.Names}}')

if [ -z "$containers" ]; then
    echo "No running containers."
    exit 0
fi

echo "=== Running Containers ==="
i=1
echo "$containers" | while read -r name; do
    echo "  $i) $name"
    i=$((i + 1))
done
echo ""

printf "Container number to tail logs: "
read -r num

target=$(echo "$containers" | sed -n "${num}p")

if [ -z "$target" ]; then
    echo "Invalid selection."
    exit 1
fi

exec docker logs -f "$target"
