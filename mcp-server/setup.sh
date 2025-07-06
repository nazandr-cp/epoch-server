#!/bin/bash

# Install Python MCP server dependencies
echo "Installing Python MCP server dependencies..."

# Check if python3 is available
if ! command -v python3 &> /dev/null; then
    echo "Error: python3 is not installed"
    exit 1
fi

# Install pip if not available
if ! command -v pip3 &> /dev/null; then
    echo "Error: pip3 is not installed"
    exit 1
fi

# Install dependencies
pip3 install -r requirements.txt

echo "MCP server dependencies installed successfully!"
echo "To use the MCP server, make sure your epoch-server is running on http://localhost:8088"