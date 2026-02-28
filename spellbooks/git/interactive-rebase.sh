#!/bin/sh
set -e

printf "How many commits back to rebase? [default: 5] "
read -r n
n="${n:-5}"

exec git rebase -i "HEAD~${n}"
