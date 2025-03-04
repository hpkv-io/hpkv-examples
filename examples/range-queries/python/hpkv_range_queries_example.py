import os
import requests
from dotenv import load_dotenv
import json
from typing import Dict, Any, Optional

# Load environment variables
load_dotenv()

class HPKVRangeQueriesExample:
    def __init__(self):
        self.api_key = os.getenv("HPKV_API_KEY")
        self.base_url = os.getenv("HPKV_API_BASE_URL")
        if not self.api_key or not self.base_url:
            raise ValueError("Please set HPKV_API_KEY and HPKV_API_BASE_URL in your .env file")

    def _get_headers(self) -> Dict[str, str]:
        """Get headers for API requests."""
        return {
            "x-api-key": self.api_key
        }

    def create_sample_records(self) -> None:
        """Create sample records for demonstration."""
        # Create some sample user records with sequential IDs
        for i in range(1, 11):
            user_data = {
                "name": f"User {i}",
                "email": f"user{i}@example.com",
                "age": 20 + i,
                "city": "New York" if i % 2 == 0 else "San Francisco"
            }
            
            payload = {
                "key": f"user:{i}",
                "value": json.dumps(user_data)
            }
            
            response = requests.post(
                f"{self.base_url}/record",
                headers=self._get_headers(),
                json=payload
            )
            response.raise_for_status()
            print(f"Created record for user:{i}")

    def perform_range_query(self, start_key: str, end_key: str, limit: Optional[int] = None) -> Dict[str, Any]:
        """
        Perform a range query to retrieve records within a specified key range.
        
        Args:
            start_key: Starting key for the range (inclusive)
            end_key: Ending key for the range (inclusive)
            limit: Maximum number of records to return (optional)
            
        Returns:
            Dictionary containing the response data
        """
        params = {
            "startKey": start_key,
            "endKey": end_key
        }
        
        if limit is not None:
            params["limit"] = limit

        response = requests.get(
            f"{self.base_url}/records",
            headers=self._get_headers(),
            params=params
        )
        response.raise_for_status()
        
        return response.json()

    def cleanup_records(self) -> None:
        """Clean up the sample records."""
        for i in range(1, 11):
            response = requests.delete(
                f"{self.base_url}/record/user:{i}",
                headers=self._get_headers()
            )
            response.raise_for_status()
            print(f"Deleted record for user:{i}")

def main():
    # Initialize the example
    example = HPKVRangeQueriesExample()
    
    try:
        # Create sample records
        print("\nCreating sample records...")
        example.create_sample_records()
        
        # Example 1: Basic range query
        print("\nExample 1: Basic range query (users 1-5)")
        result = example.perform_range_query("user:1", "user:5")
        print(json.dumps(result, indent=2))
        
        # Example 2: Range query with limit
        print("\nExample 2: Range query with limit (users 1-10, limit 3)")
        result = example.perform_range_query("user:1", "user:9", limit=3)
        print(json.dumps(result, indent=2))
        
        # Example 3: Range query for specific city
        print("\nExample 3: Range query for users in New York (even IDs)")
        result = example.perform_range_query("user:2", "user:9")
        new_york_users = [record for record in result["records"] 
                         if json.loads(record["value"])["city"] == "New York"]
        print(json.dumps(new_york_users, indent=2))
       
        
    finally:
        # Clean up the sample records
        print("\nCleaning up sample records...")
        #example.cleanup_records()

if __name__ == "__main__":
    main() 