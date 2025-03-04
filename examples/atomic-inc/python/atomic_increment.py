import os
import requests
import json
from typing import Dict, Any, Optional
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Get HPKV configuration from environment variables
HPKV_BASE_URL = os.getenv("HPKV_BASE_URL")
HPKV_API_KEY = os.getenv("HPKV_API_KEY")

def create_key(key: str, initial_value: int = 0) -> bool:
    """
    Create a new key with an initial value.
    
    Args:
        key (str): The key to create
        initial_value (int): Initial value for the key (defaults to 0)
    
    Returns:
        bool: True if creation was successful, False otherwise
    """
    try:
        # Store the value as a string representation of the number
        payload = {
            "key": key,
            "value": str(initial_value)  # Convert to string to ensure proper number format
        }
        print(f"Creating key with payload: {json.dumps(payload)}")
        
        response = requests.post(
            f"{HPKV_BASE_URL}/record",
            headers={
                "Content-Type": "application/json",
                "x-api-key": HPKV_API_KEY
            },
            json=payload,
            timeout=10
        )
        
        print(f"Create response status: {response.status_code}")
        print(f"Create response body: {response.text}")
        
        if response.status_code not in [200, 201]:
            error_msg = f"Create failed with status {response.status_code}"
            if response.text:
                try:
                    error_data = response.json()
                    error_msg += f" - {error_data.get('message', response.text)}"
                except:
                    error_msg += f" - {response.text}"
            print(error_msg)
            return False
            
        return True
            
    except Exception as e:
        print(f"Error creating record: {str(e)}")
        return False

def atomic_increment(key: str, increment: int) -> Dict[str, Any]:
    """
    Perform an atomic increment operation on a key in HPKV.
    If the key doesn't exist, it will be created with an initial value of 0.
    
    Args:
        key (str): The key to increment
        increment (int): The value to add (positive) or subtract (negative)
    
    Returns:
        Dict[str, Any]: Response from the HPKV API
    
    Raises:
        ValueError: If the API request fails or returns an error
    """
    if not HPKV_BASE_URL or not HPKV_API_KEY:
        raise ValueError("HPKV_BASE_URL and HPKV_API_KEY must be set in environment variables")

    try:
        # First, try to increment the key
        payload = {
            "key": key,
            "increment": increment
        }
        print(f"Attempting atomic increment with payload: {json.dumps(payload)}")
        
        response = requests.post(
            f"{HPKV_BASE_URL}/record/atomic",
            headers={
                "Content-Type": "application/json",
                "x-api-key": HPKV_API_KEY
            },
            json=payload,
        )
        
        print(f"Atomic increment response status: {response.status_code}")
        print(f"Atomic increment response body: {response.text}")
        
        # If the key doesn't exist (404), create it first and try again
        if response.status_code == 404:
            print(f"Key '{key}' doesn't exist. Creating it with initial value 0...")
            if not create_key(key, 0):
                raise ValueError("Failed to create key with initial value")
            
            # Retry the increment operation
            print("Retrying atomic increment after key creation...")
            response = requests.post(
                f"{HPKV_BASE_URL}/record/atomic",
                headers={
                    "Content-Type": "application/json",
                    "x-api-key": HPKV_API_KEY
                },
                json=payload,
            )
            
            print(f"Retry response status: {response.status_code}")
            print(f"Retry response body: {response.text}")
        
        # Raise an exception for bad status codes
        response.raise_for_status()
        
        result = response.json()
        
        # Check if the API returned an error message
        if not result.get("success", False):
            error_msg = result.get("message", "Unknown error")
            raise ValueError(f"HPKV API error: {error_msg}")
            
        return result
        
    except requests.exceptions.RequestException as e:
        raise ValueError(f"Failed to connect to HPKV API: {str(e)}")
    except ValueError as e:
        raise ValueError(f"Invalid response from HPKV API: {str(e)}")

def main():
    # Example usage
    key = "counter:example"
    
    try:
        # First, let's delete the key if it exists to start fresh
        print("Cleaning up any existing key...")
        response = requests.delete(
            f"{HPKV_BASE_URL}/record/{key}",
            headers={
                "Content-Type": "application/json",
                "x-api-key": HPKV_API_KEY
            }
        )
        create_key(key, 0)
        r1 = requests.get(f"{HPKV_BASE_URL}/record/{key}", headers={
            "Content-Type": "application/json",
            "x-api-key": HPKV_API_KEY
        })
        print(f"Get response status: {r1.status_code}")
        print(f"Delete response status: {response.status_code}")
        print(f"Delete response body: {response.text}")
        print(f"Get response status: {r1.status_code}")
        
        
        # Increment by 1 (will create the key if it doesn't exist)
        print("\nIncrementing counter by 1...")
        result = atomic_increment(key, 1)
        print(f"Result: {result}")
        
        # Increment by 5
        print("\nIncrementing counter by 5...")
        result = atomic_increment(key, 5)
        print(f"Result: {result}")
        
        # Decrement by 2
        print("\nDecrementing counter by 2...")
        result = atomic_increment(key, -2)
        print(f"Result: {result}")
        
    except ValueError as e:
        print(f"Error: {str(e)}")
        exit(1)
    except Exception as e:
        print(f"Unexpected error: {str(e)}")
        exit(1)

if __name__ == "__main__":
    main() 