#!/bin/sh

exec git log --graph --oneline --decorate --color=always \
    --format="%C(bold blue)%h%C(reset) %C(white)%s%C(reset) %C(dim)(%cr)%C(reset)%C(bold yellow)%d%C(reset)" \
    -30
