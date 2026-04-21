#!/bin/sh
set -e

# Generate .env file in /app/config if config directory exists
if [ -f /app/scripts/generate-env.sh ]; then
    # Use /app/config as the configuration directory if it exists
    if [ -d /app/config ]; then
        # Ensure .env.example exists in config directory
        if [ ! -f /app/config/.env.example ]; then
            cp /app/.env.example /app/config/.env.example 2>/dev/null || true
        fi
        cd /app/config
        /app/scripts/generate-env.sh
    else
        cd /app
        ./scripts/generate-env.sh
    fi
else
    echo "Warning: generate-env.sh not found. Using existing .env file."
fi

# Run SOCKS5 proxy
exec /app/socks5-proxy "$@"