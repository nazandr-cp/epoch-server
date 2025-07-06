#!/usr/bin/env python3

import subprocess
import sys
import os

# Change to the script directory
os.chdir(os.path.dirname(os.path.abspath(__file__)))

# Run the MCP server with uv
cmd = [
    "uv", "run",
    "--with", "mcp==1.0.0",
    "--with", "httpx==0.28.1", 
    "--with", "python-dotenv==1.0.1",
    "--with", "pydantic==2.10.4",
    "python", "main.py"
]

try:
    subprocess.run(cmd, check=True)
except subprocess.CalledProcessError as e:
    print(f"Error running MCP server: {e}")
    sys.exit(1)
except KeyboardInterrupt:
    print("MCP server stopped")
    sys.exit(0)