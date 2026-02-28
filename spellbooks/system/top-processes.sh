#!/bin/sh

echo "=== Top 10 by CPU ==="
ps aux --sort=-%cpu 2>/dev/null | head -11 || ps aux -r | head -11
echo ""

echo "=== Top 10 by Memory ==="
ps aux --sort=-%mem 2>/dev/null | head -11 || ps aux -m | head -11
