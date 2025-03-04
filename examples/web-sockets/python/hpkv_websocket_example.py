#!/usr/bin/env python3

import os
import sys
import json
import asyncio
import websockets
import ssl
from typing import Any, Dict, Optional, Union
from enum import Enum
from dotenv import load_dotenv

class OperationCode(Enum):
    """Enumeration of HPKV WebSocket operation codes."""
    GET = 1
    INSERT = 2
    UPDATE = 3
    DELETE = 4

class HPKVWebSocketClient:
    def __init__(self, base_url: str = None, api_key: str = None):
        """Initialize HPKV WebSocket client with API key.
        
        Args:
            base_url: HPKV server URL (optional, defaults to HPKV_BASE_URL env var)
            api_key: HPKV API key (optional, defaults to HPKV_API_KEY env var)
        """
        # Load environment variables if not already loaded
        if not os.getenv('HPKV_BASE_URL') and not os.getenv('HPKV_API_KEY'):
            load_dotenv()
            
        self.base_url = (base_url or os.getenv('HPKV_BASE_URL')).rstrip('/') if (base_url or os.getenv('HPKV_BASE_URL')) else None
        self.api_key = api_key or os.getenv('HPKV_API_KEY')
        
        if not self.base_url:
            raise ValueError("HPKV base URL not provided. Set HPKV_BASE_URL environment variable or pass base_url parameter.")
        if not self.api_key:
            raise ValueError("HPKV API key not provided. Set HPKV_API_KEY environment variable or pass api_key parameter.")
        
        # Convert HTTP URL to WebSocket URL and add API key as query parameter
        self.ws_url = f"{self.base_url.replace('http://', 'ws://').replace('https://', 'wss://')}/ws?apiKey={self.api_key}"
        self.message_id = 0
        self.response_futures = {}
        self.websocket = None
        self.message_handler_task = None
        
        # Create SSL context
        self.ssl_context = ssl.create_default_context()
        self.ssl_context.check_hostname = False
        self.ssl_context.verify_mode = ssl.CERT_NONE
        
    def _get_next_message_id(self) -> int:
        """Get next available message ID."""
        self.message_id += 1
        return self.message_id
    
    def _serialize_value(self, value: Any) -> str:
        """Serialize value to string format."""
        if isinstance(value, str):
            return value
        return json.dumps(value)
    
    async def _handle_message(self, message: str):
        """Handle incoming WebSocket message."""
        try:
            response = json.loads(message)
            message_id = response.get('messageId')
            if message_id in self.response_futures:
                future = self.response_futures.pop(message_id)
                if 'error' in response:
                    future.set_exception(Exception(response['error']))
                else:
                    future.set_result(response)
        except Exception as e:
            print(f"Error handling message: {str(e)}", file=sys.stderr)
    
    async def _message_handler(self):
        """Handle incoming WebSocket messages."""
        try:
            async for message in self.websocket:
                await self._handle_message(message)
        except websockets.exceptions.ConnectionClosed:
            print("WebSocket connection closed", file=sys.stderr)
        except Exception as e:
            print(f"Error in message handler: {str(e)}", file=sys.stderr)
    
    async def connect(self):
        """Establish WebSocket connection."""
        if not self.websocket:
            self.websocket = await websockets.connect(
                self.ws_url,
                ssl=self.ssl_context
            )
            self.message_handler_task = asyncio.create_task(self._message_handler())
    
    async def disconnect(self):
        """Close WebSocket connection."""
        if self.message_handler_task:
            self.message_handler_task.cancel()
            try:
                await self.message_handler_task
            except asyncio.CancelledError:
                pass
        if self.websocket:
            await self.websocket.close()
            self.websocket = None
        self.message_handler_task = None
    
    async def _send_message(self, message: Dict) -> Any:
        """Send message and wait for response."""
        if not self.websocket:
            await self.connect()
            
        message_id = self._get_next_message_id()
        message['messageId'] = message_id
        
        future = asyncio.Future()
        self.response_futures[message_id] = future
        
        try:
            await self.websocket.send(json.dumps(message))
            return await future
        except Exception as e:
            print(f"Error sending message: {str(e)}", file=sys.stderr)
            raise
    
    async def create(self, key: str, value: Any) -> bool:
        """Create a new key-value pair."""
        try:
            message = {
                "op": OperationCode.INSERT.value,  # Insert operation
                "key": key,
                "value": self._serialize_value(value)
            }
            
            response = await self._send_message(message)
            return 'error' not in response
                
        except Exception as e:
            print(f"Error creating record: {str(e)}", file=sys.stderr)
            return False

    async def read(self, key: str) -> Optional[Any]:
        """Read a value by key."""
        try:
            message = {
                "op": OperationCode.GET.value,  # Get operation
                "key": key
            }
            
            response = await self._send_message(message)
            if 'error' in response:
                return None
                
            try:
                return json.loads(response['value'])
            except (json.JSONDecodeError, TypeError):
                return response['value']
                
        except Exception as e:
            print(f"Error reading record: {str(e)}", file=sys.stderr)
            return None

    async def update(self, key: str, value: Any, partial_update: bool = False) -> bool:
        """Update an existing key-value pair."""
        try:
            message = {
                "op": OperationCode.UPDATE.value if partial_update else OperationCode.INSERT.value,  # Update or Insert operation
                "key": key,
                "value": self._serialize_value(value)
            }
            
            response = await self._send_message(message)
            return 'error' not in response
                
        except Exception as e:
            print(f"Error updating record: {str(e)}", file=sys.stderr)
            return False

    async def delete(self, key: str) -> bool:
        """Delete a key-value pair."""
        try:
            message = {
                "op": OperationCode.DELETE.value,  # Delete operation
                "key": key
            }
            
            response = await self._send_message(message)
            return 'error' not in response
                
        except Exception as e:
            print(f"Error deleting record: {str(e)}", file=sys.stderr)
            return False

async def main():
    try:
        # Load environment variables from .env file
        load_dotenv()
        
        # Initialize HPKV WebSocket client using environment variables
        client = HPKVWebSocketClient()
        
        print("HPKV WebSocket CRUD Operations Example")
        print("=====================================")
        print(f"\nUsing HPKV WebSocket server: {client.ws_url}")

        # Create operation
        user_data = {
            "name": "John Doe",
            "email": "john@example.com",
            "age": 30
        }
        print("\n1. Creating a new user record...")
        success = await client.create(
            key="user:1",
            value=user_data
        )
        if not success:
            print("Failed to create record. Exiting...")
            sys.exit(1)
        print("Create operation succeeded")

        # Read operation
        print("\n2. Reading the user record...")
        retrieved_data = await client.read("user:1")
        if retrieved_data:
            print(f"Retrieved data: {json.dumps(retrieved_data, indent=2)}")
        else:
            print("Failed to retrieve data")

        # Update operation
        print("\n3. Updating the user's age...")
        user_data["age"] = 31
        success = await client.update("user:1", user_data)
        print(f"Update operation {'succeeded' if success else 'failed'}")

        # Read after update
        print("\n4. Reading the updated user record...")
        retrieved_data = await client.read("user:1")
        if retrieved_data:
            print(f"Retrieved data: {json.dumps(retrieved_data, indent=2)}")
        else:
            print("Failed to retrieve data")

        # Delete operation
        print("\n5. Deleting the user record...")
        success = await client.delete("user:1")
        print(f"Delete operation {'succeeded' if success else 'failed'}")

        # Verify deletion
        print("\n6. Attempting to read deleted record...")
        retrieved_data = await client.read("user:1")
        if retrieved_data is None:
            print("Record was successfully deleted")
        else:
            print("Record still exists")
            
        # Clean up
        await client.disconnect()
            
    except Exception as e:
        print(f"Error running example: {str(e)}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    asyncio.run(main()) 