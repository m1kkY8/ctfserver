#!/bin/sh

# Create directories if they don't exist
mkdir -p /opt/tools /opt/loot

# Ensure proper permissions
chmod 755 /opt/tools /opt/loot

# Start the CTF server
exec /usr/local/bin/ctfserver "$@"
