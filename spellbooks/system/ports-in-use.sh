#!/bin/sh

echo "=== Listening Ports ==="
lsof -i -P -n | grep LISTEN
