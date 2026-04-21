#!/bin/sh
set -e

CONFIG_FILE=/app/config/.env

if [ ! -f "$CONFIG_FILE" ]; then
    if [ -f /app/config/.env.example ]; then
        EXAMPLE_FILE=/app/config/.env.example
    else
        EXAMPLE_FILE=/app/.env.example
    fi
    /app/scripts/generate-env.sh "$EXAMPLE_FILE" "$CONFIG_FILE"
fi

exec /app/socks5-proxy "$@"
