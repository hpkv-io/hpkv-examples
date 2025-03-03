#!/usr/bin/env python3

import os
import sys
import requests
import json
import time
from typing import Any, Dict, Optional, Union
from dotenv import load_dotenv

class HPKVClient:
    def __init__(self, base_url: str = None, api_key: str = None):
        """Initialize HPKV client with API key.
        
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
        
        self.headers = {
            'Content-Type': 'application/json',
            'x-api-key': self.api_key
        }
    
    def _serialize_value(self, value: Any) -> str:
        """Serialize value to string format.
        
        Args:
            value: Value to serialize
            
        Returns:
            str: Serialized value
        """
        if isinstance(value, str):
            return value
        return json.dumps(value)
    
    def _handle_response(self, response: requests.Response, operation: str) -> Union[Dict, str, None]:
        """Handle API response and extract data.
        
        Args:
            response: Response from API
            operation: Operation name for error messages
            
        Returns:
            Response data if successful, None otherwise
        """
        try:
            if response.status_code in [200, 201]:
                return response.json() if response.text else None
            else:
                error_msg = f"Error in {operation}: Status {response.status_code}"
                if response.text:
                    try:
                        error_data = response.json()
                        error_msg += f" - {error_data.get('message', response.text)}"
                    except:
                        error_msg += f" - {response.text}"
                print(error_msg, file=sys.stderr)
                return None
        except Exception as e:
            print(f"Error parsing response in {operation}: {str(e)}", file=sys.stderr)
            return None
        
    def create(self, key: str, value: Any) -> bool:
        """Create a new key-value pair.
        
        Args:
            key: Key to store
            value: Value to store
            
        Returns:
            bool: True if creation was successful, False otherwise
        """
        try:
            payload = {
                "key": key,
                "value": self._serialize_value(value)
            }
            
            print(f"Sending request to {self.base_url}/record")
            print(f"Payload: {json.dumps(payload, indent=2)}")
            
            response = requests.post(
                f"{self.base_url}/record",
                headers=self.headers,
                json=payload
            )
            
            if response.status_code not in [200, 201]:
                error_msg = f"Create failed with status {response.status_code}"
                if response.text:
                    try:
                        error_data = response.json()
                        error_msg += f" - {error_data.get('message', response.text)}"
                    except:
                        error_msg += f" - {response.text}"
                print(error_msg, file=sys.stderr)
                return False
                
            print(f"Create succeeded with status {response.status_code}")
            return True
                
        except Exception as e:
            print(f"Error creating record: {str(e)}", file=sys.stderr)
            return False

    def read(self, key: str) -> Optional[Any]:
        """Read a value by key."""
        try:
            response = requests.get(
                f"{self.base_url}/record/{key}",
                headers=self.headers
            )
            
            data = self._handle_response(response, "read")
            if data and 'value' in data:
                try:
                    return json.loads(data['value'])
                except (json.JSONDecodeError, TypeError):
                    return data['value']
            return None
            
        except Exception as e:
            print(f"Error reading record: {str(e)}", file=sys.stderr)
            return None

    def update(self, key: str, value: Any, partial_update: bool = False) -> bool:
        """Update an existing key-value pair."""
        try:
            payload = {
                "key": key,
                "value": self._serialize_value(value),
                "partialUpdate": partial_update
            }
            
            response = requests.post(
                f"{self.base_url}/record",
                headers=self.headers,
                json=payload
            )
            
            return response.status_code == 200
            
        except Exception as e:
            print(f"Error updating record: {str(e)}", file=sys.stderr)
            return False

    def delete(self, key: str) -> bool:
        """Delete a key-value pair."""
        try:
            response = requests.delete(
                f"{self.base_url}/record/{key}",
                headers=self.headers
            )
            
            return response.status_code == 200
            
        except Exception as e:
            print(f"Error deleting record: {str(e)}", file=sys.stderr)
            return False

def main():
    try:
        # Load environment variables from .env file
        load_dotenv()
        
        # Initialize HPKV client using environment variables
        client = HPKVClient()
        
        print("HPKV CRUD Operations Example")
        print("===========================")
        print(f"\nUsing HPKV server: {client.base_url}")

        # Create operation
        user_data = {
            "name": "John Doe",
            "email": "john@example.com",
            "age": 30
        }
        print("\n1. Creating a new user record...")
        success = client.create(
            key="user:1",
            value=user_data
        )
        if not success:
            print("Failed to create record. Exiting...")
            sys.exit(1)
        print("Create operation succeeded")

        # Read operation
        print("\n2. Reading the user record...")
        retrieved_data = client.read("user:1")
        if retrieved_data:
            print(f"Retrieved data: {json.dumps(retrieved_data, indent=2)}")
        else:
            print("Failed to retrieve data")

        # Update operation
        print("\n3. Updating the user's age...")
        user_data["age"] = 31
        success = client.update("user:1", user_data)
        print(f"Update operation {'succeeded' if success else 'failed'}")

        # Read after update
        print("\n4. Reading the updated user record...")
        retrieved_data = client.read("user:1")
        if retrieved_data:
            print(f"Retrieved data: {json.dumps(retrieved_data, indent=2)}")
        else:
            print("Failed to retrieve data")

        # Delete operation
        print("\n5. Deleting the user record...")
        success = client.delete("user:1")
        print(f"Delete operation {'succeeded' if success else 'failed'}")

        # Verify deletion
        print("\n6. Attempting to read deleted record...")
        retrieved_data = client.read("user:1")
        if retrieved_data is None:
            print("Record was successfully deleted")
        else:
            print("Record still exists")
            
    except Exception as e:
        print(f"Error running example: {str(e)}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main() 
    