#!/bin/bash

# This script runs the bot with a dummy token for testing purposes
echo "Starting bot with dummy token for testing..."
TELEGRAM_BOT_TOKEN="dummy_token_for_testing" docker-compose up --remove-orphans bot