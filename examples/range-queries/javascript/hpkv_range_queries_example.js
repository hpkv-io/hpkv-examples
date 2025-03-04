require('dotenv').config();
const axios = require('axios');

class HPKVRangeQueriesExample {
    constructor() {
        this.apiKey = process.env.HPKV_API_KEY;
        this.baseUrl = process.env.HPKV_API_BASE_URL;
        
        if (!this.apiKey || !this.baseUrl) {
            throw new Error("Please set HPKV_API_KEY and HPKV_API_BASE_URL in your .env file");
        }
    }

    _getHeaders() {
        return {
            "x-api-key": this.apiKey
        };
    }

    async createSampleRecords() {
        // Create some sample user records with sequential IDs
        for (let i = 1; i <= 10; i++) {
            const userData = {
                name: `User ${i}`,
                email: `user${i}@example.com`,
                age: 20 + i,
                city: i % 2 === 0 ? "New York" : "San Francisco"
            };
            
            const payload = {
                key: `user:${i}`,
                value: JSON.stringify(userData)
            };
            
            try {
                await axios.post(
                    `${this.baseUrl}/record`,
                    payload,
                    { headers: this._getHeaders() }
                );
                console.log(`Created record for user:${i}`);
            } catch (error) {
                console.error(`Error creating record for user:${i}:`, error.message);
                throw error;
            }
        }
    }

    async performRangeQuery(startKey, endKey, limit = null) {
        const params = {
            startKey,
            endKey
        };
        
        if (limit !== null) {
            params.limit = limit;
        }

        try {
            const response = await axios.get(
                `${this.baseUrl}/records`,
                {
                    headers: this._getHeaders(),
                    params
                }
            );
            return response.data;
        } catch (error) {
            console.error('Error performing range query:', error.message);
            throw error;
        }
    }

    async cleanupRecords() {
        for (let i = 1; i <= 10; i++) {
            try {
                await axios.delete(
                    `${this.baseUrl}/record/user:${i}`,
                    { headers: this._getHeaders() }
                );
                console.log(`Deleted record for user:${i}`);
            } catch (error) {
                console.error(`Error deleting record for user:${i}:`, error.message);
                throw error;
            }
        }
    }
}

async function main() {
    const example = new HPKVRangeQueriesExample();
    
    try {
        // Create sample records
        console.log("\nCreating sample records...");
        await example.createSampleRecords();
        
        // Example 1: Basic range query
        console.log("\nExample 1: Basic range query (users 1-5)");
        const result1 = await example.performRangeQuery("user:1", "user:5");
        console.log(JSON.stringify(result1, null, 2));
        
        // Example 2: Range query with limit
        console.log("\nExample 2: Range query with limit (users 1-10, limit 3)");
        const result2 = await example.performRangeQuery("user:1", "user:9", 3);
        console.log(JSON.stringify(result2, null, 2));
        
        // Example 3: Range query for specific city
        console.log("\nExample 3: Range query for users in New York (even IDs)");
        const result3 = await example.performRangeQuery("user:2", "user:9");
        const newYorkUsers = result3.records.filter(record => 
            JSON.parse(record.value).city === "New York"
        );
        console.log(JSON.stringify(newYorkUsers, null, 2));
        
    } finally {
        // Clean up the sample records
        console.log("\nCleaning up sample records...");
        // await example.cleanupRecords();
    }
}

main().catch(console.error); 