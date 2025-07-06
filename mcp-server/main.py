#!/usr/bin/env python3

import asyncio
import os
import sys
from typing import Any, Dict, List, Optional

import httpx
from dotenv import load_dotenv
from mcp.server import Server
from mcp.server.stdio import stdio_server
from mcp.types import Tool, TextContent, CallToolRequest, CallToolResult
from pydantic import BaseModel, Field

# Load environment variables
load_dotenv()

# Configuration
EPOCH_SERVER_URL = os.getenv("EPOCH_SERVER_URL", "http://localhost:8080")

class EpochMCPServer:
    def __init__(self):
        self.server = Server("epoch-server-mcp")
        self.http_client = httpx.AsyncClient(timeout=30.0)
        self.setup_tools()

    def setup_tools(self):
        """Set up MCP tools for epoch server endpoints"""
        
        # Health Check Tool
        health_check_tool = Tool(
            name="health_check",
            description="Check the health status of the epoch server",
            inputSchema={
                "type": "object",
                "properties": {},
                "required": [],
            },
        )
        
        # Start Epoch Tool
        start_epoch_tool = Tool(
            name="start_epoch",
            description="Start a new epoch for the lending platform",
            inputSchema={
                "type": "object",
                "properties": {
                    "epoch_id": {
                        "type": "string",
                        "description": "The ID of the epoch to start",
                    }
                },
                "required": ["epoch_id"],
            },
        )
        
        # Distribute Subsidies Tool
        distribute_subsidies_tool = Tool(
            name="distribute_subsidies",
            description="Distribute subsidies for an epoch",
            inputSchema={
                "type": "object",
                "properties": {
                    "epoch_id": {
                        "type": "string",
                        "description": "The ID of the epoch to distribute subsidies for",
                    }
                },
                "required": ["epoch_id"],
            },
        )

        # Register tools
        @self.server.call_tool()
        async def handle_health_check(name: str, arguments: Dict[str, Any]) -> List[TextContent]:
            """Handle health check requests"""
            try:
                response = await self.http_client.get(f"{EPOCH_SERVER_URL}/health")
                response.raise_for_status()
                
                result = response.json()
                return [
                    TextContent(
                        type="text",
                        text=f"Health check successful: {result}",
                    )
                ]
            except Exception as e:
                return [
                    TextContent(
                        type="text",
                        text=f"Health check failed: {str(e)}",
                    )
                ]

        @self.server.call_tool()
        async def handle_start_epoch(name: str, arguments: Dict[str, Any]) -> List[TextContent]:
            """Handle start epoch requests"""
            try:
                epoch_id = arguments.get("epoch_id")
                if not epoch_id:
                    return [
                        TextContent(
                            type="text",
                            text="Error: epoch_id is required",
                        )
                    ]
                
                response = await self.http_client.post(
                    f"{EPOCH_SERVER_URL}/epochs/{epoch_id}/start"
                )
                response.raise_for_status()
                
                result = response.json()
                return [
                    TextContent(
                        type="text",
                        text=f"Epoch {epoch_id} started successfully: {result}",
                    )
                ]
            except httpx.HTTPStatusError as e:
                error_detail = ""
                try:
                    error_detail = e.response.json()
                except:
                    error_detail = e.response.text
                
                return [
                    TextContent(
                        type="text",
                        text=f"Failed to start epoch {epoch_id}: HTTP {e.response.status_code} - {error_detail}",
                    )
                ]
            except Exception as e:
                return [
                    TextContent(
                        type="text",
                        text=f"Failed to start epoch {epoch_id}: {str(e)}",
                    )
                ]

        @self.server.call_tool()
        async def handle_distribute_subsidies(name: str, arguments: Dict[str, Any]) -> List[TextContent]:
            """Handle distribute subsidies requests"""
            try:
                epoch_id = arguments.get("epoch_id")
                if not epoch_id:
                    return [
                        TextContent(
                            type="text",
                            text="Error: epoch_id is required",
                        )
                    ]
                
                response = await self.http_client.post(
                    f"{EPOCH_SERVER_URL}/epochs/{epoch_id}/distribute"
                )
                response.raise_for_status()
                
                result = response.json()
                return [
                    TextContent(
                        type="text",
                        text=f"Subsidies distributed for epoch {epoch_id}: {result}",
                    )
                ]
            except httpx.HTTPStatusError as e:
                error_detail = ""
                try:
                    error_detail = e.response.json()
                except:
                    error_detail = e.response.text
                
                return [
                    TextContent(
                        type="text",
                        text=f"Failed to distribute subsidies for epoch {epoch_id}: HTTP {e.response.status_code} - {error_detail}",
                    )
                ]
            except Exception as e:
                return [
                    TextContent(
                        type="text",
                        text=f"Failed to distribute subsidies for epoch {epoch_id}: {str(e)}",
                    )
                ]

        # Register tools with the server
        self.server.list_tools = lambda: [
            health_check_tool,
            start_epoch_tool,
            distribute_subsidies_tool,
        ]

    async def run(self):
        """Run the MCP server"""
        async with stdio_server() as (read_stream, write_stream):
            await self.server.run(
                read_stream,
                write_stream,
                self.server.create_initialization_options(),
            )

async def main():
    """Main entry point"""
    server = EpochMCPServer()
    await server.run()

if __name__ == "__main__":
    asyncio.run(main())