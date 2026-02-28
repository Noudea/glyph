#!/bin/sh

exec git log --oneline --graph --all --decorate --color=always | head -40
