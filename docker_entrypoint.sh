#!/bin/sh
set -e

# check if the user has provided the required environment variables
if [ -z "$USERNAME" ]; then
  echo "Error:: USERNAME is not set. Exiting."
  exit 1
fi

if [ -z "$PASSWORD" ]; then
  echo "Error: PASSWORD is not set. Exiting."
  exit 1
fi

node setup_docker.js

# Execute the CMD from the Dockerfile (i.e., "npm run start")
exec "$@"
