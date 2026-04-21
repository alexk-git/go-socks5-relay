#!/bin/sh
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Paths
ENV_EXAMPLE="config/.env.example"
ENV_FILE="config/.env"

# Allow overriding paths via arguments
if [ $# -ge 1 ]; then
    ENV_EXAMPLE="$1"
fi
if [ $# -ge 2 ]; then
    ENV_FILE="$2"
fi

# Check if .env.example exists
if [ ! -f "$ENV_EXAMPLE" ]; then
    echo -e "${RED}Error: $ENV_EXAMPLE not found.${NC}"
    exit 1
fi

# Check if .env already exists
if [ -f "$ENV_FILE" ]; then
    echo -e "${YELLOW}Info: $ENV_FILE already exists. Skipping generation.${NC}"
    exit 0
fi

echo -e "${GREEN}Generating $ENV_FILE from $ENV_EXAMPLE with random credentials...${NC}"

# Copy .env.example to .env
cp "$ENV_EXAMPLE" "$ENV_FILE"

# Generate random username (user_ + 8 random alphanumeric characters)
RANDOM_USERNAME="user_$(openssl rand -hex 4 2>/dev/null || cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 8 | head -n 1)"

# Generate random password (16 random alphanumeric characters)
RANDOM_PASSWORD="$(openssl rand -hex 12 2>/dev/null || cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 16 | head -n 1)"

# Replace username and password in .env file
# Using sed compatible with both GNU sed and BSD sed (macOS)
if sed --version 2>&1 | grep -q GNU; then
    # GNU sed
    sed -i "s/^username=.*/username=$RANDOM_USERNAME/" "$ENV_FILE"
    sed -i "s/^password=.*/password=$RANDOM_PASSWORD/" "$ENV_FILE"
else
    # BSD sed
    sed -i "" "s/^username=.*/username=$RANDOM_USERNAME/" "$ENV_FILE"
    sed -i "" "s/^password=.*/password=$RANDOM_PASSWORD/" "$ENV_FILE"
fi

# Also ensure port is set to 1080 (as per .env.example)
# (optional) Uncomment if you want to ensure default port
# sed -i "s/^port=.*/port=1080/" "$ENV_FILE"

echo -e "${GREEN}Generated $ENV_FILE with:${NC}"
echo -e "  Username: ${YELLOW}$RANDOM_USERNAME${NC}"
echo -e "  Password: ${YELLOW}$RANDOM_PASSWORD${NC}"
echo -e "${GREEN}Please save these credentials for connecting to the proxy.${NC}"
