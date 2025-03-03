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

# Determine if the --enable-otp flag should be added.
OTP_FLAG=""
if [ "$ENABLE_OTP" = "true" ]; then
  OTP_FLAG="--enable-otp"
fi

# only set the allowed origins if it is provided
if [ "$ALLOWED_IPS" ]; then
  ALLOWED_IPS_FLAG="--allowed-origins $ALLOWED_IPS"
fi



node setup.js \
  --username $USERNAME \
  --password $PASSWORD \
  $OTP_FLAG \
  $ALLOWED_IPS_FLAG


# Execute the CMD from the Dockerfile (i.e., "npm run start")
exec "$@"
